package napcat

import (
	"context"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/zjutjh/napcat-sdk/api"
	"github.com/zjutjh/napcat-sdk/event"
	"github.com/zjutjh/napcat-sdk/transport"
)

// Option 配置 SDK client。
type Option func(*options)

type options struct {
	token                string
	httpTimeout          time.Duration
	httpClient           *http.Client
	wsDialOptions        *websocket.DialOptions
	requestTimeout       time.Duration
	eventBuffer          int
	eventDeliveryTimeout time.Duration
	logger               *slog.Logger
}

// Client 是 SDK 的统一入口。
type Client struct {
	api *api.Client
	ws  *transport.WebSocketCaller
}

// WithToken 设置 NapCat token。
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

// WithEventDeliveryTimeout 设置事件缓冲区满后的最长等待时间。超时会终止 WebSocket。
func WithEventDeliveryTimeout(timeout time.Duration) Option {
	return func(o *options) { o.eventDeliveryTimeout = timeout }
}

// WithLogger 设置可选日志器。
func WithLogger(logger *slog.Logger) Option {
	return func(o *options) { o.logger = logger }
}

// NewHTTPClient 创建基于 HTTP API 的 client。
func NewHTTPClient(baseURL string, opts ...Option) *Client {
	cfg := collectOptions(opts...)
	caller := transport.NewHTTPCaller(baseURL, transport.HTTPOptions{
		Token:   cfg.token,
		Timeout: cfg.httpTimeout,
		Client:  cfg.httpClient,
	})
	return newClient(caller, nil)
}

// DialWebSocket 创建基于正向 WebSocket 的 client。
func DialWebSocket(ctx context.Context, url string, opts ...Option) (*Client, error) {
	cfg := collectOptions(opts...)
	caller, err := transport.DialWebSocket(ctx, url, transport.WebSocketOptions{
		Token:                cfg.token,
		RequestTimeout:       cfg.requestTimeout,
		EventBuffer:          cfg.eventBuffer,
		EventDeliveryTimeout: cfg.eventDeliveryTimeout,
		DialOptions:          cfg.wsDialOptions,
		Logger:               cfg.logger,
	})
	if err != nil {
		return nil, err
	}
	return newClient(caller, caller), nil
}

// ServeReverseWebSocket 监听反向 WebSocket，并把每个连接包装为 Client。
func ServeReverseWebSocket(ctx context.Context, addr string, handler func(*Client), opts ...Option) error {
	cfg := collectOptions(opts...)
	var clientsMu sync.Mutex
	clients := make(map[*Client]struct{})
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
			RequestTimeout:       cfg.requestTimeout,
			EventBuffer:          cfg.eventBuffer,
			EventDeliveryTimeout: cfg.eventDeliveryTimeout,
			Logger:               cfg.logger,
		})
		client := newClient(caller, caller)
		clientsMu.Lock()
		if shuttingDown {
			clientsMu.Unlock()
			_ = caller.CloseNow()
			return
		}
		clients[client] = struct{}{}
		handlers.Add(1)
		clientsMu.Unlock()
		go func() {
			defer handlers.Done()
			defer func() {
				clientsMu.Lock()
				delete(clients, client)
				clientsMu.Unlock()
				_ = client.Close()
			}()
			handler(client)
		}()
	})
	server := &http.Server{Addr: addr, Handler: mux}
	var shutdownOnce sync.Once
	shutdown := func() {
		shutdownOnce.Do(func() {
			clientsMu.Lock()
			shuttingDown = true
			active := make([]*Client, 0, len(clients))
			for client := range clients {
				active = append(active, client)
			}
			clientsMu.Unlock()
			_ = server.Close()
			for _, client := range active {
				_ = client.ws.CloseNow()
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

// Events 返回事件流。HTTP client 没有事件流，返回 nil。
func (c *Client) Events() <-chan event.Event {
	if c.ws == nil {
		return nil
	}
	return c.ws.Events()
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
	return cfg
}

func newClient(caller transport.Caller, ws *transport.WebSocketCaller) *Client {
	return &Client{api: api.NewClient(caller), ws: ws}
}
