// Package errorsx 保存 SDK 内部共享的错误类型。
package errorsx

import (
	"errors"
	"fmt"
)

// ErrClosed 表示 client 或底层连接已经关闭。
var ErrClosed = errors.New("onebot: client closed")

// ErrTimeout 表示请求等待响应超时。
var ErrTimeout = errors.New("onebot: request timeout")

// ErrEventBackpressure 表示事件消费者持续无法接收事件。
var ErrEventBackpressure = errors.New("onebot: event delivery backpressure")

// TransportError 表示 HTTP 或 WebSocket 传输层失败。
type TransportError struct {
	Op     string
	Status int
	Body   []byte
	Err    error
}

func (e *TransportError) Error() string {
	message := "onebot: transport error"
	if e.Op != "" {
		message += " during " + e.Op
	}
	if e.Status != 0 {
		message += fmt.Sprintf(" (HTTP %d)", e.Status)
	}
	if e.Err != nil {
		return message + ": " + e.Err.Error()
	}
	return message
}

func (e *TransportError) Unwrap() error {
	return e.Err
}

// APIError 表示后端返回了合法 envelope，但业务状态失败。
type APIError struct {
	Action  string
	Status  string
	RetCode int
	Message string
	Wording string
	Raw     []byte
}

func (e *APIError) Error() string {
	message := fmt.Sprintf("onebot: api error during %s (retcode %d)", e.Action, e.RetCode)
	if e.Message != "" {
		return message + ": " + e.Message
	}
	if e.Wording != "" {
		return message + ": " + e.Wording
	}
	return message
}

// ProtocolError 表示 SDK 无法解析 envelope、响应或事件。
type ProtocolError struct {
	Message string
	Raw     []byte
}

func (e *ProtocolError) Error() string {
	if e.Message == "" {
		return "onebot: protocol error"
	}
	return "onebot: protocol error: " + e.Message
}
