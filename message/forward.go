package message

import (
	"fmt"

	json "github.com/zjutjh/onebot-sdk/internal/jsonx"
)

// ForwardNode 是合并转发消息的单个节点。
// 两后端的 get_forward_msg 均返回消息形态对象(user_id/sender.nickname/time/message),
// 部分实现把消息体放在 content 字段,此处一并兼容。
type ForwardNode struct {
	UserID   ID
	Nickname string
	Time     int64
	Message  Chain
}

type wireForwardNode struct {
	UserID ID `json:"user_id"`
	Sender struct {
		Nickname string `json:"nickname"`
	} `json:"sender"`
	Time    int64           `json:"time"`
	Message json.RawMessage `json:"message,omitempty"`
	Content json.RawMessage `json:"content,omitempty"`
}

// UnmarshalJSON 从消息形态节点归一化字段,消息体优先取 message,缺失时取 content。
func (n *ForwardNode) UnmarshalJSON(data []byte) error {
	var raw wireForwardNode
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("解析转发节点失败: %w", err)
	}
	n.UserID = raw.UserID
	n.Nickname = raw.Sender.Nickname
	n.Time = raw.Time

	payload := raw.Message
	if json.IsNull(payload) {
		payload = raw.Content
	}
	if json.IsNull(payload) {
		n.Message = nil
		return nil
	}
	var chain Chain
	if err := json.Unmarshal(payload, &chain); err != nil {
		return fmt.Errorf("解析转发节点消息体失败: %w", err)
	}
	n.Message = chain
	return nil
}

// MarshalJSON 编码回消息形态节点。
func (n ForwardNode) MarshalJSON() ([]byte, error) {
	var raw wireForwardNode
	raw.UserID = n.UserID
	raw.Sender.Nickname = n.Nickname
	raw.Time = n.Time
	payload, err := json.Marshal(n.Message)
	if err != nil {
		return nil, fmt.Errorf("编码转发节点消息体失败: %w", err)
	}
	if len(n.Message) > 0 {
		raw.Message = payload
	}
	return json.Marshal(raw)
}
