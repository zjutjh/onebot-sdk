package onebot

import "github.com/zjutjh/onebot-sdk/internal/errorsx"

// ErrClosed 表示 client 或底层连接已经关闭。
var ErrClosed = errorsx.ErrClosed

// ErrTimeout 表示请求等待响应超时。
var ErrTimeout = errorsx.ErrTimeout

// ErrEventBackpressure 表示事件消费者持续无法接收事件。
var ErrEventBackpressure = errorsx.ErrEventBackpressure

// TransportError 表示 HTTP 或 WebSocket 传输层失败。
type TransportError = errorsx.TransportError

// APIError 表示后端返回了合法 envelope，但业务状态失败。
type APIError = errorsx.APIError

// ProtocolError 表示 SDK 无法解析 envelope、响应或事件。
type ProtocolError = errorsx.ProtocolError
