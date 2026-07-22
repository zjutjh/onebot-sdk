package transport

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/zjutjh/napcat-sdk/internal/errorsx"
	json "github.com/zjutjh/napcat-sdk/internal/jsonx"
)

// HTTPOptions 配置 HTTP 调用器。
type HTTPOptions struct {
	Token       string
	Timeout     time.Duration
	RetryCount  int
	RetryWait   time.Duration
	RestyClient *resty.Client
}

// HTTPCaller 使用 NapCat HTTP API 调用 action。
type HTTPCaller struct {
	baseURL string
	token   string
	client  *resty.Client
}

// NewHTTPCaller 创建 HTTP 调用器。
func NewHTTPCaller(baseURL string, opts HTTPOptions) *HTTPCaller {
	client := opts.RestyClient
	if client == nil {
		client = resty.New()
	}
	if opts.Timeout > 0 {
		client.SetTimeout(opts.Timeout)
	}
	if opts.RetryCount > 0 {
		client.SetRetryCount(opts.RetryCount)
	}
	if opts.RetryWait > 0 {
		client.SetRetryWaitTime(opts.RetryWait)
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
	Echo    string          `json:"echo,omitempty"`
}

// Call 通过 HTTP POST 调用一个 NapCat action。
func (c *HTTPCaller) Call(ctx context.Context, action string, params any, result any) error {
	body, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("编码 HTTP 请求失败: %w", err)
	}

	req := c.client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(body)
	if c.token != "" {
		req.SetHeader("Authorization", "Bearer "+c.token)
	}

	resp, err := req.Post(c.baseURL + "/" + strings.TrimLeft(action, "/"))
	if err != nil {
		return &errorsx.TransportError{Op: action, Err: err}
	}
	raw := append([]byte(nil), resp.Body()...)
	if resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
		return &errorsx.TransportError{Op: action, Status: resp.StatusCode(), Body: raw}
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
	if result == nil || len(env.Data) == 0 || string(env.Data) == "null" {
		return nil
	}
	if err := json.Unmarshal(env.Data, result); err != nil {
		return &errorsx.ProtocolError{Message: "解析响应 data 失败", Raw: raw}
	}
	return nil
}
