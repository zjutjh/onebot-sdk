package transport

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/zjutjh/onebot-sdk/internal/errorsx"
	json "github.com/zjutjh/onebot-sdk/internal/jsonx"
)

// HTTPOptions 配置 HTTP 调用器。
type HTTPOptions struct {
	Token   string
	Timeout time.Duration
	Client  *http.Client
}

// HTTPCaller 使用 OneBot HTTP API 调用 action。
type HTTPCaller struct {
	baseURL string
	token   string
	client  *http.Client
}

// NewHTTPCaller 创建 HTTP 调用器。
func NewHTTPCaller(baseURL string, opts HTTPOptions) *HTTPCaller {
	client := opts.Client
	if client == nil {
		client = &http.Client{Timeout: opts.Timeout}
	}
	if opts.Timeout > 0 && client.Timeout != opts.Timeout {
		clone := *client
		clone.Timeout = opts.Timeout
		client = &clone
	}
	return &HTTPCaller{
		baseURL: strings.TrimRight(baseURL, "/"),
		token:   opts.Token,
		client:  client,
	}
}

type envelope struct {
	Status  *string         `json:"status"`
	RetCode *int            `json:"retcode"`
	Data    json.RawMessage `json:"data"`
	Message string          `json:"message"`
	Wording string          `json:"wording"`
}

// Call 通过 HTTP POST 调用一个 action。
func (c *HTTPCaller) Call(ctx context.Context, action string, params any, result any) error {
	body, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("编码 HTTP 请求失败: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/"+strings.TrimLeft(action, "/"), bytes.NewReader(body))
	if err != nil {
		return &errorsx.TransportError{Op: action, Err: err}
	}
	req.Header.Set("Content-Type", "application/json")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return &errorsx.TransportError{Op: action, Err: err}
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return &errorsx.TransportError{Op: action, Status: resp.StatusCode, Err: err}
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return &errorsx.TransportError{Op: action, Status: resp.StatusCode, Body: raw}
	}

	return decodeEnvelope(action, raw, result)
}

func decodeEnvelope(action string, raw []byte, result any) error {
	var env envelope
	if err := json.Unmarshal(raw, &env); err != nil {
		return &errorsx.ProtocolError{Message: "解析响应 envelope 失败", Raw: raw}
	}
	if env.Status == nil || env.RetCode == nil {
		return &errorsx.ProtocolError{Message: "响应 envelope 缺少 status 或 retcode", Raw: raw}
	}
	if *env.Status != "ok" || *env.RetCode != 0 {
		return &errorsx.APIError{
			Action:  action,
			Status:  *env.Status,
			RetCode: *env.RetCode,
			Message: env.Message,
			Wording: env.Wording,
			Raw:     raw,
		}
	}
	if result == nil || json.IsNull(env.Data) {
		return nil
	}
	if err := json.Unmarshal(env.Data, result); err != nil {
		return &errorsx.ProtocolError{Message: "解析响应 data 失败", Raw: raw}
	}
	return nil
}
