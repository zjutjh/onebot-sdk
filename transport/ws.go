package transport

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/coder/websocket"
	"github.com/zjutjh/napcat-sdk/event"
	"github.com/zjutjh/napcat-sdk/internal/errorsx"
	json "github.com/zjutjh/napcat-sdk/internal/jsonx"
)

// WebSocketOptions 配置 WebSocket 调用器。
type WebSocketOptions struct {
	Token          string
	RequestTimeout time.Duration
	// EventBuffer 是待交付事件队列大小。
	EventBuffer int
	// EventDeliveryTimeout 为正时，事件队列持续满到超时会终止连接。
	EventDeliveryTimeout time.Duration
	DialOptions          *websocket.DialOptions
	Logger               *slog.Logger
}

// WebSocketCaller 使用 OneBot WebSocket action 模型调用 NapCat。
type WebSocketCaller struct {
	conn                 *websocket.Conn
	requestTimeout       time.Duration
	eventDeliveryTimeout time.Duration
	events               chan event.Event
	pending              map[string]chan []byte
	stateMu              sync.Mutex
	terminalErr          error
	closed               chan struct{}
	terminateOnce        sync.Once
	seq                  atomic.Int64
	logger               *slog.Logger
}

const webSocketReadLimit = 64 << 20

type wsRequest struct {
	Action string `json:"action"`
	Params any    `json:"params"`
	Echo   string `json:"echo"`
}

type echoProbe struct {
	Echo string `json:"echo"`
}

// DialWebSocket 连接正向 WebSocket，并启动后台读循环。
func DialWebSocket(ctx context.Context, url string, opts WebSocketOptions) (*WebSocketCaller, error) {
	dialOpts := &websocket.DialOptions{}
	if opts.DialOptions != nil {
		*dialOpts = *opts.DialOptions
	}
	if opts.Token != "" {
		dialOpts.HTTPHeader = dialOpts.HTTPHeader.Clone()
		if dialOpts.HTTPHeader == nil {
			dialOpts.HTTPHeader = make(http.Header)
		}
		dialOpts.HTTPHeader.Set("Authorization", "Bearer "+opts.Token)
	}
	conn, _, err := websocket.Dial(ctx, url, dialOpts)
	if err != nil {
		return nil, &errorsx.TransportError{Op: "websocket dial", Err: err}
	}
	return NewWebSocketCaller(conn, opts), nil
}

// NewWebSocketCaller 从已有连接创建调用器，供反向 WebSocket server 复用。
func NewWebSocketCaller(conn *websocket.Conn, opts WebSocketOptions) *WebSocketCaller {
	eventBuffer := opts.EventBuffer
	if eventBuffer <= 0 {
		eventBuffer = 16
	}
	timeout := opts.RequestTimeout
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	c := &WebSocketCaller{
		conn:                 conn,
		requestTimeout:       timeout,
		eventDeliveryTimeout: opts.EventDeliveryTimeout,
		events:               make(chan event.Event, eventBuffer),
		pending:              make(map[string]chan []byte),
		closed:               make(chan struct{}),
		logger:               opts.Logger,
	}
	conn.SetReadLimit(webSocketReadLimit)
	go c.readLoop()
	return c
}

// Events 返回事件流。
func (c *WebSocketCaller) Events() <-chan event.Event {
	return c.events
}

// Call 通过 WebSocket action 调用 NapCat。
func (c *WebSocketCaller) Call(ctx context.Context, action string, params any, result any) error {
	echo := strconv.FormatInt(c.seq.Add(1), 10)
	respCh := make(chan []byte, 1)
	if err := c.addPending(echo, respCh); err != nil {
		return err
	}
	defer c.removePending(echo)

	waitCtx, cancel := context.WithTimeout(ctx, c.requestTimeout)
	defer cancel()

	payload, err := json.Marshal(wsRequest{Action: action, Params: params, Echo: echo})
	if err != nil {
		return fmt.Errorf("编码 WebSocket 请求失败: %w", err)
	}
	if err := c.conn.Write(waitCtx, websocket.MessageText, payload); err != nil {
		if terminalErr := c.Err(); terminalErr != nil {
			return terminalErr
		}
		if errors.Is(waitCtx.Err(), context.DeadlineExceeded) && ctx.Err() == nil {
			return errorsx.ErrTimeout
		}
		return &errorsx.TransportError{Op: action, Err: err}
	}

	select {
	case raw := <-respCh:
		return decodeEnvelope(action, raw, result)
	case <-waitCtx.Done():
		select {
		case raw := <-respCh:
			return decodeEnvelope(action, raw, result)
		default:
		}
		if errors.Is(waitCtx.Err(), context.DeadlineExceeded) && ctx.Err() == nil {
			return errorsx.ErrTimeout
		}
		return waitCtx.Err()
	case <-c.closed:
		select {
		case raw := <-respCh:
			return decodeEnvelope(action, raw, result)
		default:
		}
		return c.Err()
	}
}

// Close 关闭连接并唤醒所有 pending 调用。
func (c *WebSocketCaller) Close() error {
	return c.terminate(errorsx.ErrClosed, false)
}

// CloseNow 不等待关闭握手，立即关闭连接。
func (c *WebSocketCaller) CloseNow() error {
	return c.terminate(errorsx.ErrClosed, true)
}

// Err 返回连接的首次终止原因，连接运行中返回 nil。
func (c *WebSocketCaller) Err() error {
	c.stateMu.Lock()
	defer c.stateMu.Unlock()
	return c.terminalErr
}

func (c *WebSocketCaller) addPending(echo string, ch chan []byte) error {
	c.stateMu.Lock()
	defer c.stateMu.Unlock()
	select {
	case <-c.closed:
		return c.terminalErr
	default:
	}
	c.pending[echo] = ch
	return nil
}

func (c *WebSocketCaller) removePending(echo string) {
	c.stateMu.Lock()
	delete(c.pending, echo)
	c.stateMu.Unlock()
}

func (c *WebSocketCaller) readLoop() {
	defer close(c.events)
	for {
		_, data, err := c.conn.Read(context.Background())
		if err != nil {
			c.terminate(webSocketReadError(err), true)
			return
		}
		var probe echoProbe
		if err := json.Unmarshal(data, &probe); err == nil && probe.Echo != "" {
			c.stateMu.Lock()
			select {
			case <-c.closed:
			default:
				if ch := c.pending[probe.Echo]; ch != nil {
					select {
					case ch <- data:
					default:
					}
				}
			}
			c.stateMu.Unlock()
			continue
		}
		ev := event.Parse(data)
		if failure, ok := ev.(event.ParseFailure); ok {
			if c.logger != nil {
				c.logger.Warn("无法强类型解析 NapCat WebSocket 事件，已保留原始事件", "error", failure.ParseError())
			}
		}
		if !c.enqueueEvent(ev) {
			return
		}
	}
}

func (c *WebSocketCaller) enqueueEvent(ev event.Event) bool {
	select {
	case c.events <- ev:
		return true
	case <-c.closed:
		return false
	default:
	}

	if c.eventDeliveryTimeout <= 0 {
		if c.logger != nil {
			c.logger.Warn("NapCat 事件缓冲区已满，丢弃事件", "post_type", ev.PostType())
		}
		return true
	}

	timer := time.NewTimer(c.eventDeliveryTimeout)
	defer timer.Stop()
	select {
	case c.events <- ev:
		return true
	case <-c.closed:
		return false
	case <-timer.C:
		if c.logger != nil {
			c.logger.Warn("NapCat 事件交付超时，终止连接", "timeout", c.eventDeliveryTimeout)
		}
		c.terminate(errorsx.ErrEventBackpressure, true)
		return false
	}
}

func (c *WebSocketCaller) terminate(cause error, immediate bool) error {
	won := false
	c.terminateOnce.Do(func() {
		won = true
		c.stateMu.Lock()
		c.terminalErr = cause
		close(c.closed)
		c.stateMu.Unlock()
	})
	if !won {
		if immediate {
			return c.conn.CloseNow()
		}
		return nil
	}
	if immediate {
		return c.conn.CloseNow()
	}
	return c.conn.Close(websocket.StatusNormalClosure, "")
}

func webSocketReadError(err error) error {
	status := websocket.CloseStatus(err)
	if status == websocket.StatusNormalClosure || status == websocket.StatusGoingAway {
		return errorsx.ErrClosed
	}
	return &errorsx.TransportError{Op: "websocket read", Err: err}
}
