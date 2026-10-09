package onebot

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/zjutjh/onebot-sdk/api"
	"github.com/zjutjh/onebot-sdk/dialect"
	"github.com/zjutjh/onebot-sdk/event"
	"github.com/zjutjh/onebot-sdk/transport"
)

// Option 配置 SDK client。
type Option func(*options)

// Backend 指定后端类型,用于跳过方言检测或强制保守模式。
type Backend int

const (
	// BackendAuto 按连接自报的 app_name 自动检测方言。
	BackendAuto Backend = iota
	// BackendNapCat 显式指定 NapCat。
	BackendNapCat
	// BackendSnowLuma 显式指定 SnowLuma。
	BackendSnowLuma
	// BackendGeneric 显式指定未知后端,只用 OneBot 11 核心能力。
	BackendGeneric
)

type options struct {
	token                string
	httpTimeout          time.Duration
	httpClient           *http.Client
	wsDialOptions        *websocket.DialOptions
	requestTimeout       time.Duration
	eventBuffer          int
	eventDeliveryTimeout time.Duration
	logger               *slog.Logger
	backend              Backend
}

// dialectDetectTimeout 限制 get_version_info 方言检测的耗时。
const dialectDetectTimeout = 5 * time.Second

// defaultEventBuffer 是未配置 WithEventBuffer 时事件队列的默认容量;
// transport 包对直接构造的调用方另持同值兜底,正常路径经根包统一传入。
const defaultEventBuffer = 16

// Client 是 SDK 的统一入口。
// 事件流默认经过归一化管线(强类型解析 + 方言补拉);RawEvents 提供原始旁路。
type Client struct {
	api    *api.Client
	ws     *transport.WebSocketCaller
	events chan event.Event
	opts   options

	dialectMu       sync.Mutex
	dialectVal      dialect.Dialect
	dialectResolved bool
}

// WithToken 设置后端访问 token。
func WithToken(token string) Option {
	return func(o *options) { o.token = token }
}

// WithHTTPTimeout 设置 HTTP 超时。
func WithHTTPTimeout(timeout time.Duration) Option {
	return func(o *options) { o.httpTimeout = timeout }
}

// WithHTTPClient 注入自定义 HTTP client。
func WithHTTPClient(client *http.Client) Option {
	return func(o *options) { o.httpClient = client }
}

// WithWebSocketDialOptions 注入 WebSocket 连接参数。
func WithWebSocketDialOptions(opts *websocket.DialOptions) Option {
	return func(o *options) { o.wsDialOptions = opts }
}

// WithRequestTimeout 设置 WebSocket 请求默认等待时间。
func WithRequestTimeout(timeout time.Duration) Option {
	return func(o *options) { o.requestTimeout = timeout }
}

// WithEventBuffer 设置待交付事件队列大小。默认模式下队列满时丢弃新事件。
func WithEventBuffer(size int) Option {
	return func(o *options) { o.eventBuffer = size }
}

// WithEventDeliveryTimeout 设置事件交付超时。
// 用户在超时时长内持续不消费事件队列时,以 ErrEventBackpressure 终止连接。
// 该超时只作用于管线到用户这一跳,归一化补拉的 API 往返不计入,不会因后端响应慢误杀连接。
func WithEventDeliveryTimeout(timeout time.Duration) Option {
	return func(o *options) { o.eventDeliveryTimeout = timeout }
}

// WithLogger 设置可选日志器。
func WithLogger(logger *slog.Logger) Option {
	return func(o *options) { o.logger = logger }
}

// WithBackend 显式指定后端,默认按 get_version_info 自动检测。
func WithBackend(backend Backend) Option {
	return func(o *options) { o.backend = backend }
}

// NewHTTPClient 创建基于 HTTP API 的 client。
// HTTP 没有事件流;方言在首次需要时懒检测。
func NewHTTPClient(baseURL string, opts ...Option) *Client {
	cfg := collectOptions(opts...)
	caller := transport.NewHTTPCaller(baseURL, transport.HTTPOptions{
		Token:   cfg.token,
		Timeout: cfg.httpTimeout,
		Client:  cfg.httpClient,
	})
	// 事件流返回已关闭通道:消费方 range 立即结束,nil channel 会永久阻塞
	events := make(chan event.Event)
	close(events)
	return &Client{api: api.NewClient(caller), opts: cfg, events: events}
}

// DialWebSocket 创建基于正向 WebSocket 的 client。
// 连接建立后同步执行方言检测(失败降级为 Generic),随后启动归一化管线。
func DialWebSocket(ctx context.Context, url string, opts ...Option) (*Client, error) {
	cfg := collectOptions(opts...)
	caller, err := transport.DialWebSocket(ctx, url, transport.WebSocketOptions{
		Token:          cfg.token,
		RequestTimeout: cfg.requestTimeout,
		EventBuffer:    cfg.eventBuffer,
		DialOptions:    cfg.wsDialOptions,
		Logger:         cfg.logger,
	})
	if err != nil {
		return nil, err
	}
	d, _ := resolveDialect(ctx, caller, cfg.backend, cfg.logger)
	return newWSClient(caller, cfg, d), nil
}

// ServeReverseWebSocket 监听反向 WebSocket，并把每个连接包装为 Client。
// 每个连接独立执行方言检测和归一化管线。
func ServeReverseWebSocket(ctx context.Context, addr string, handler func(*Client), opts ...Option) error {
	cfg := collectOptions(opts...)
	var connMu sync.Mutex
	conns := make(map[*transport.WebSocketCaller]struct{})
	var handlers sync.WaitGroup
	var shuttingDown bool

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if cfg.token != "" && r.Header.Get("Authorization") != "Bearer "+cfg.token {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		conn, err := websocket.Accept(w, r, nil)
		if err != nil {
			return
		}
		caller := transport.NewWebSocketCaller(conn, transport.WebSocketOptions{
			RequestTimeout: cfg.requestTimeout,
			EventBuffer:    cfg.eventBuffer,
			Logger:         cfg.logger,
		})
		// 关闭判定、在途连接登记和 WaitGroup 登记在同一临界区:
		// 关闭后接受的连接直接拒绝;关闭前接受的连接必定在 shutdown 快照里,
		// 且 Add 不会与 Wait 并发。
		connMu.Lock()
		if shuttingDown {
			connMu.Unlock()
			_ = caller.CloseNow()
			return
		}
		conns[caller] = struct{}{}
		handlers.Add(1)
		connMu.Unlock()
		go func() {
			defer handlers.Done()
			defer func() {
				connMu.Lock()
				delete(conns, caller)
				connMu.Unlock()
				_ = caller.Close()
			}()
			// 方言检测最长 5 秒,期间连接可能已被 shutdown 关闭,
			// 登记与调用 handler 前须重新确认,避免关闭后仍进入 handler 阻塞关闭流程。
			d, _ := resolveDialect(context.Background(), caller, cfg.backend, cfg.logger)
			connMu.Lock()
			closing := shuttingDown
			connMu.Unlock()
			if closing {
				return
			}
			handler(newWSClient(caller, cfg, d))
		}()
	})
	server := &http.Server{Addr: addr, Handler: mux}
	var shutdownOnce sync.Once
	shutdown := func() {
		shutdownOnce.Do(func() {
			connMu.Lock()
			shuttingDown = true
			active := make([]*transport.WebSocketCaller, 0, len(conns))
			for caller := range conns {
				active = append(active, caller)
			}
			connMu.Unlock()
			_ = server.Close()
			for _, caller := range active {
				_ = caller.CloseNow()
			}
		})
	}
	done := make(chan struct{})
	go func() {
		select {
		case <-ctx.Done():
			shutdown()
		case <-done:
		}
	}()
	err := server.ListenAndServe()
	close(done)
	shutdown()
	handlers.Wait()
	if err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

// API 返回强类型 API client。
func (c *Client) API() *api.Client {
	return c.api
}

// Events 返回归一化事件流。HTTP client 没有事件流,返回已关闭通道,range 立即结束。
func (c *Client) Events() <-chan event.Event {
	return c.events
}

// closedRawEvents 供 HTTP client 的 RawEvents 返回已关闭通道,避免 nil channel 永久阻塞。
var closedRawEvents = func() chan []byte {
	ch := make(chan []byte)
	close(ch)
	return ch
}()

// RawEvents 返回原始事件 JSON 的旁路流,供绕过归一化管线自行解析。
// 与 Events 并存:每个事件帧同时进入两者,不是二选一。
// 须在事件开始到达前调用,之后才调用的帧不会回放;HTTP client 返回已关闭通道。
func (c *Client) RawEvents() <-chan []byte {
	if c.ws == nil {
		return closedRawEvents
	}
	return c.ws.TeeRawEvents()
}

// Dialect 返回当前连接解析出的方言。
// WebSocket 连接在建立时已解析;HTTP client 首次调用时懒检测。
func (c *Client) Dialect(ctx context.Context) dialect.Dialect {
	return c.currentDialect(ctx)
}

// Call 调用原始 action。
func (c *Client) Call(ctx context.Context, action string, params any, result any) error {
	return c.api.Call(ctx, action, params, result)
}

// Close 关闭底层连接。HTTP client 调用该方法无副作用。
func (c *Client) Close() error {
	if c.ws == nil {
		return nil
	}
	return c.ws.Close()
}

// Err 返回 WebSocket 的终止原因。HTTP client 始终返回 nil。
func (c *Client) Err() error {
	if c.ws == nil {
		return nil
	}
	return c.ws.Err()
}

func collectOptions(opts ...Option) options {
	var cfg options
	for _, opt := range opts {
		opt(&cfg)
	}
	if cfg.eventBuffer <= 0 {
		cfg.eventBuffer = defaultEventBuffer
	}
	return cfg
}

func newWSClient(caller *transport.WebSocketCaller, cfg options, d dialect.Dialect) *Client {
	c := &Client{
		api:    api.NewClient(caller),
		ws:     caller,
		events: make(chan event.Event, cfg.eventBuffer),
		opts:   cfg,
	}
	c.dialectMu.Lock()
	c.dialectVal, c.dialectResolved = d, true
	c.dialectMu.Unlock()
	go c.runPipeline()
	return c
}

// currentDialect 返回已解析方言,HTTP client 首次调用时懒检测。
// 持锁检测让并发调用共用同一次结果;只有稳定判定才落缓存,
// 否则一次网络抖动会把健康后端永久判成 Generic。
func (c *Client) currentDialect(ctx context.Context) dialect.Dialect {
	c.dialectMu.Lock()
	defer c.dialectMu.Unlock()
	if c.dialectResolved {
		return c.dialectVal
	}
	d, cacheable := resolveDialect(ctx, c.api, c.opts.backend, c.opts.logger)
	if cacheable {
		c.dialectVal, c.dialectResolved = d, true
	}
	return d
}

// resolveDialect 按配置解析方言,第二个返回值表示结果是否可缓存。
// 显式指定与已识别的 app_name 可缓存;未知 app_name 是稳定判定,缓存为 Generic;
// 其余调用失败不缓存,下次调用重试。
func resolveDialect(ctx context.Context, caller transport.Caller, backend Backend, logger *slog.Logger) (dialect.Dialect, bool) {
	if backend != BackendAuto {
		return backendDialect(backend), true
	}
	detectCtx, cancel := context.WithTimeout(ctx, dialectDetectTimeout)
	defer cancel()
	d, err := dialect.Detect(detectCtx, caller)
	if err != nil {
		if logger != nil {
			logger.Warn("方言检测失败，降级为通用 OneBot 11 模式", "error", err)
		}
		return dialect.Generic, errors.Is(err, dialect.ErrUnknownBackend)
	}
	return d, true
}

func backendDialect(backend Backend) dialect.Dialect {
	switch backend {
	case BackendNapCat:
		return dialect.NapCat
	case BackendSnowLuma:
		return dialect.SnowLuma
	default:
		return dialect.Generic
	}
}
