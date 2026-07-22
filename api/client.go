//go:generate go run ../cmd/napcatgen -spec ../internal/openapi/4.18.13/openapi.json -out .

// Package api 提供生成 API 方法共享的运行时类型。
package api

import (
	"context"

	json "github.com/zjutjh/napcat-sdk/internal/jsonx"
	"github.com/zjutjh/napcat-sdk/transport"
)

// Client 调用 NapCat action。强类型方法由生成代码补充。
type Client struct {
	caller transport.Caller
}

// NewOB11Message 把字符串、消息段或消息链编码为发送接口使用的消息联合类型。
func NewOB11Message(value any) (OB11MessageMixTypeUnion, error) {
	raw, err := json.Marshal(value)
	return OB11MessageMixTypeUnion{Raw: raw}, err
}

// NewClient 创建 API client。
func NewClient(caller transport.Caller) *Client {
	return &Client{caller: caller}
}

// Call 调用原始 action，供尚未强类型化的新接口兜底。
func (c *Client) Call(ctx context.Context, action string, params any, result any) error {
	return c.caller.Call(ctx, action, params, result)
}
