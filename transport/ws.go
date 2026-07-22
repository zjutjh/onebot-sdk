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
	// EventBuffer 已满时丢弃新事件，避免阻塞同一连接上的 API 响应。
	EventBuffer int
	DialOptions *websocket.DialOptions
	Logger      *slog.Logger
}

// WebSocketCaller 使用 OneBot WebSocket action 模型调用 NapCat。
type WebSocketCaller struct {
	conn           *websocket.Conn
	requestTimeout time.Duration
	events         chan event.Event
	pending        map[string]chan []byte
	pendingMu      sync.Mutex
	closed         chan struct{}
	closeOnce      sync.Once
	seq            atomic.Int64
	logger         *slog.Logger
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
		conn:           conn,
		requestTimeout: timeout,
		events:         make(chan event.Event, eventBuffer),
		pending:        make(map[string]chan []byte),
		closed:         make(chan struct{}),
		logger:         opts.Logger,
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
		if errors.Is(waitCtx.Err(), context.DeadlineExceeded) && ctx.Err() == nil {
			return errorsx.ErrTimeout
		}
		return &errorsx.TransportError{Op: action, Err: err}
	}

	select {
	case raw := <-respCh:
		return decodeEnvelope(action, raw, result)
	case <-waitCtx.Done():
		if errors.Is(waitCtx.Err(), context.DeadlineExceeded) && ctx.Err() == nil {
			return errorsx.ErrTimeout
		}
		return waitCtx.Err()
	case <-c.closed:
		return errorsx.ErrClosed
	}
}

// Close 关闭连接并唤醒所有 pending 调用。
func (c *WebSocketCaller) Close() error {
	var err error
	c.closeOnce.Do(func() {
		c.pendingMu.Lock()
		close(c.closed)
		c.pendingMu.Unlock()
		err = c.conn.Close(websocket.StatusNormalClosure, "")
	})
	return err
}

func (c *WebSocketCaller) addPending(echo string, ch chan []byte) error {
	c.pendingMu.Lock()
	defer c.pendingMu.Unlock()
	select {
	case <-c.closed:
		return errorsx.ErrClosed
	default:
	}
	c.pending[echo] = ch
	return nil
}

func (c *WebSocketCaller) removePending(echo string) {
	c.pendingMu.Lock()
	delete(c.pending, echo)
	c.pendingMu.Unlock()
}

func (c *WebSocketCaller) readLoop() {
	defer close(c.events)
	defer c.Close()
	for {
		_, data, err := c.conn.Read(context.Background())
		if err != nil {
			return
		}
		var probe echoProbe
		if err := json.Unmarshal(data, &probe); err == nil && probe.Echo != "" {
			c.pendingMu.Lock()
			ch := c.pending[probe.Echo]
			c.pendingMu.Unlock()
			if ch != nil {
				select {
				case ch <- data:
				default:
				}
			}
			continue
		}
		ev, err := event.Parse(data)
		if err != nil {
			if c.logger != nil {
				c.logger.Warn("无法解析 NapCat WebSocket 消息", "error", err)
			}
			continue
		}
		select {
		case c.events <- ev:
		case <-c.closed:
			return
		default:
			if c.logger != nil {
				c.logger.Warn("NapCat 事件缓冲区已满，丢弃事件", "post_type", ev.PostType())
			}
		}
	}
}
