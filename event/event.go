// Package event 提供 OneBot 上报事件的解析和基础类型。
package event

import (
	"bytes"
	"fmt"

	json "github.com/zjutjh/napcat-sdk/internal/jsonx"
	"github.com/zjutjh/napcat-sdk/message"
)

// Event 是所有上报事件的公共接口。
type Event interface {
	PostType() string
	SelfID() int64
	Time() int64
	Raw() []byte
}

// Base 保存所有事件共有字段。
type Base struct {
	EventTime     int64  `json:"time"`
	EventPostType string `json:"post_type"`
	EventSelfID   int64  `json:"self_id"`
	rawData       []byte
}

// PostType 返回事件类型。
func (b Base) PostType() string { return b.EventPostType }

// SelfID 返回接收事件的机器人 QQ。
func (b Base) SelfID() int64 { return b.EventSelfID }

// Time 返回事件发生时间戳。
func (b Base) Time() int64 { return b.EventTime }

// Raw 返回原始事件 JSON。
func (b Base) Raw() []byte { return bytes.Clone(b.rawData) }

// PrivateSender 表示私聊发送者信息。
type PrivateSender struct {
	UserID   int64  `json:"user_id"`
	Nickname string `json:"nickname"`
	Sex      string `json:"sex,omitempty"`
	Age      int    `json:"age,omitempty"`
}

// GroupSender 表示群聊发送者信息。
type GroupSender struct {
	UserID   int64  `json:"user_id"`
	Nickname string `json:"nickname"`
	Card     string `json:"card,omitempty"`
	Role     string `json:"role,omitempty"`
	Sex      string `json:"sex,omitempty"`
	Age      int    `json:"age,omitempty"`
}

// PrivateMessage 表示私聊消息事件。
type PrivateMessage struct {
	Base
	MessageType string        `json:"message_type"`
	SubType     string        `json:"sub_type"`
	MessageID   int64         `json:"message_id"`
	UserID      int64         `json:"user_id"`
	Message     message.Chain `json:"message"`
	RawMessage  string        `json:"raw_message"`
	Sender      PrivateSender `json:"sender"`
}

// GroupMessage 表示群消息事件。
type GroupMessage struct {
	Base
	MessageType string        `json:"message_type"`
	SubType     string        `json:"sub_type"`
	MessageID   int64         `json:"message_id"`
	GroupID     int64         `json:"group_id"`
	UserID      int64         `json:"user_id"`
	Message     message.Chain `json:"message"`
	RawMessage  string        `json:"raw_message"`
	Sender      GroupSender   `json:"sender"`
}

// UnknownEvent 表示当前 SDK 未细分的事件。
type UnknownEvent struct {
	Base
}

type envelope struct {
	Time        int64  `json:"time"`
	PostType    string `json:"post_type"`
	SelfID      int64  `json:"self_id"`
	MessageType string `json:"message_type"`
}

// Parse 将 OneBot 事件 JSON 解析为具体事件类型。
func Parse(data []byte) (Event, error) {
	var env envelope
	if err := json.Unmarshal(data, &env); err != nil {
		return nil, fmt.Errorf("解析事件 envelope 失败: %w", err)
	}

	base := Base{
		EventTime:     env.Time,
		EventPostType: env.PostType,
		EventSelfID:   env.SelfID,
		rawData:       bytes.Clone(data),
	}

	switch {
	case env.PostType == "message" && env.MessageType == "private":
		var ev PrivateMessage
		if err := json.Unmarshal(data, &ev); err != nil {
			return nil, fmt.Errorf("解析私聊消息失败: %w", err)
		}
		ev.Base = base
		return &ev, nil
	case env.PostType == "message" && env.MessageType == "group":
		var ev GroupMessage
		if err := json.Unmarshal(data, &ev); err != nil {
			return nil, fmt.Errorf("解析群消息失败: %w", err)
		}
		ev.Base = base
		return &ev, nil
	default:
		return &UnknownEvent{Base: base}, nil
	}
}
