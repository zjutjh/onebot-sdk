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
	"github.com/zjutjh/onebot-sdk/internal/errorsx"
	json "github.com/zjutjh/onebot-sdk/internal/jsonx"
)

// WebSocketOptions 配置 WebSocket 调用器。
type WebSocketOptions struct {
	Token          string
	RequestTimeout time.Duration
	// EventBuffer 是待归一化事件的原始队列大小,用户事件通道取同尺寸。
	EventBuffer int
	DialOptions *websocket.DialOptions
	Logger      *slog.Logger
}

// WebSocketCaller 使用 OneBot WebSocket action 模型调用后端。
type WebSocketCaller struct {
	conn           *websocket.Conn
	requestTimeout time.Duration
	rawEvents      chan []byte
	// bypass 是 RawEvents() 首次调用后启用的原始事件旁路,读循环把每帧复制一份进来
	bypass        chan []byte
	bypassOnce    sync.Once
	pending       map[string]chan []byte
	stateMu       sync.Mutex
	terminalErr   error
	closed        chan struct{}
	terminateOnce sync.Once
	seq           atomic.Int64
	logger        *slog.Logger
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
		rawEvents:      make(chan []byte, eventBuffer),
		pending:        make(map[string]chan []byte),
		closed:         make(chan struct{}),
		logger:         opts.Logger,
	}
	conn.SetReadLimit(webSocketReadLimit)
	go c.readLoop()
	return c
}

// TakeRawEvents 返回原始事件队列,由归一化管线唯一消费。
// 读循环不阻塞投递:队列满即丢弃事件并记日志,保证 action 响应永不被事件阻塞
// (读循环是 echo 响应唯一的路由方,一旦阻塞,在途调用——包括管线自身的补拉——都会饿死)。
// 交付超时只作用于管线到用户这一跳,见 pipeline.go。
func (c *WebSocketCaller) TakeRawEvents() <-chan []byte {
	return c.rawEvents
}

// closedRawEvents 是连接已终止时返回的已关闭通道,保证消费者立即收到关闭信号。
var closedRawEvents = func() <-chan []byte {
	ch := make(chan []byte)
	close(ch)
	return ch
}()

// TeeRawEvents 返回原始事件 JSON 的旁路通道。
// 首次调用后读循环把每个事件帧复制一份进来,与归一化输出并存,不是二选一。
// 旁路队列与事件队列同容量,满时丢弃该旁路帧并记日志,不影响归一化投递。
// 须在开始接收事件前调用;调用前已投递的帧不会回放。
// 连接已终止时返回已关闭的通道;通道关闭只由读循环负责。
func (c *WebSocketCaller) TeeRawEvents() <-chan []byte {
	c.bypassOnce.Do(func() {
		ch := make(chan []byte, cap(c.rawEvents))
		c.stateMu.Lock()
		select {
		case <-c.closed:
			// 连接已终止,读循环不会再复制帧,不登记通道
		default:
			c.bypass = ch
		}
		c.stateMu.Unlock()
	})
	c.stateMu.Lock()
	bypass := c.bypass
	c.stateMu.Unlock()
	if bypass == nil {
		return closedRawEvents
	}
	return bypass
}

// Call 通过 WebSocket 调用一个 action。
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
		if isRequestTimeout(waitCtx, ctx) {
			return errorsx.ErrTimeout
		}
		return &errorsx.TransportError{Op: action, Err: err}
	}

	select {
	case raw := <-respCh:
		return decodeEnvelope(action, raw, result)
	case <-waitCtx.Done():
		if raw, ok := tryResponse(respCh); ok {
			return decodeEnvelope(action, raw, result)
		}
		if isRequestTimeout(waitCtx, ctx) {
			return errorsx.ErrTimeout
		}
		return waitCtx.Err()
	case <-c.closed:
		if raw, ok := tryResponse(respCh); ok {
			return decodeEnvelope(action, raw, result)
		}
		return c.Err()
	}
}

// tryResponse 非阻塞读取已到达的响应。
// 等待与响应同时就绪时优先返回真实响应,避免把已成功的结果判成超时或关闭。
func tryResponse(respCh <-chan []byte) ([]byte, bool) {
	select {
	case raw := <-respCh:
		return raw, true
	default:
		return nil, false
	}
}

// isRequestTimeout 判断等待是否因本次请求超时结束,父 context 取消不算。
func isRequestTimeout(waitCtx, parent context.Context) bool {
	return errors.Is(waitCtx.Err(), context.DeadlineExceeded) && parent.Err() == nil
}

// Close 关闭连接并唤醒所有 pending 调用。
func (c *WebSocketCaller) Close() error {
	return c.terminate(errorsx.ErrClosed, false)
}

// Terminate 以指定原因立即终止连接,供归一化管线在事件交付超时时调用。
func (c *WebSocketCaller) Terminate(cause error) error {
	return c.terminate(cause, true)
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

// Done 返回连接终止信号,连接以任意原因终止后关闭。
// 供归一化管线在交付阻塞期间观察终止,不必等满交付超时。
func (c *WebSocketCaller) Done() <-chan struct{} {
	return c.closed
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
	defer close(c.rawEvents)
	defer c.closeBypass()
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
		c.teeBypass(data)
		if !c.enqueueRaw(data) {
			return
		}
	}
}

// teeBypass 把事件帧复制进已启用的旁路队列。
func (c *WebSocketCaller) teeBypass(data []byte) {
	c.stateMu.Lock()
	bypass := c.bypass
	c.stateMu.Unlock()
	if bypass == nil {
		return
	}
	select {
	case bypass <- data:
	default:
		if c.logger != nil {
			c.logger.Warn("事件旁路队列已满，丢弃旁路事件")
		}
	}
}

// closeBypass 在读循环退出时关闭已启用的旁路通道。
// 读循环是旁路通道唯一的发送方与关闭方,因此不存在向已关闭通道发送的竞态。
func (c *WebSocketCaller) closeBypass() {
	c.stateMu.Lock()
	defer c.stateMu.Unlock()
	if c.bypass != nil {
		close(c.bypass)
	}
}

// enqueueRaw 把事件原始 JSON 非阻塞投入队列,满队列直接丢弃并记日志。
// 读循环是 echo 响应唯一的路由方,阻塞会让在途调用(含管线的补拉 API 往返)永远等不到响应。
func (c *WebSocketCaller) enqueueRaw(data []byte) bool {
	select {
	case c.rawEvents <- data:
		return true
	case <-c.closed:
		return false
	default:
	}
	if c.logger != nil {
		c.logger.Warn("事件原始队列已满，丢弃事件")
	}
	return true
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
