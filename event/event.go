// Package event 提供 OneBot 上报事件的解析和基础类型。
package event

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"

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

// ParseFailure 表示保留了原始数据的事件解析失败。
type ParseFailure interface {
	Event
	ParseError() error
}

// ID 是兼容 JSON 整数和十进制整数字符串的事件 ID。
type ID int64

// UnmarshalJSON 解码 int64 范围内的整数 ID。
func (id *ID) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if len(data) == 0 {
		return fmt.Errorf("无效事件 ID: 空 JSON")
	}
	if data[0] != '"' {
		if data[0] != '-' && (data[0] < '0' || data[0] > '9') {
			return fmt.Errorf("无效事件 ID %q: 必须是整数或十进制整数字符串", data)
		}
		var value int64
		if err := json.Unmarshal(data, &value); err != nil {
			return fmt.Errorf("无效事件 ID %q: %w", data, err)
		}
		*id = ID(value)
		return nil
	}

	var text string
	if err := json.Unmarshal(data, &text); err != nil {
		return err
	}
	value, err := strconv.ParseInt(text, 10, 64)
	if err != nil {
		return fmt.Errorf("无效事件 ID %q: %w", text, err)
	}
	*id = ID(value)
	return nil
}

// Int64 返回 ID 的 int64 值。
func (id ID) Int64() int64 { return int64(id) }

// Base 保存所有事件共有字段。
type Base struct {
	EventTime     int64  `json:"time"`
	EventPostType string `json:"post_type"`
	EventSelfID   ID     `json:"self_id"`
	rawData       []byte
}

// PostType 返回事件类型。
func (b Base) PostType() string { return b.EventPostType }

// SelfID 返回接收事件的机器人 QQ。
func (b Base) SelfID() int64 { return b.EventSelfID.Int64() }

// Time 返回事件发生时间戳。
func (b Base) Time() int64 { return b.EventTime }

// Raw 返回原始事件 JSON。
func (b Base) Raw() []byte { return bytes.Clone(b.rawData) }

// PrivateSender 表示私聊发送者信息。
type PrivateSender struct {
	UserID   ID     `json:"user_id"`
	Nickname string `json:"nickname"`
	Sex      string `json:"sex,omitempty"`
	Age      int    `json:"age,omitempty"`
}

// GroupSender 表示群聊发送者信息。
type GroupSender struct {
	UserID   ID     `json:"user_id"`
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
	MessageID   ID            `json:"message_id"`
	UserID      ID            `json:"user_id"`
	Message     message.Chain `json:"message"`
	RawMessage  string        `json:"raw_message"`
	Sender      PrivateSender `json:"sender"`
}

// GroupMessage 表示群消息事件。
type GroupMessage struct {
	Base
	MessageType string        `json:"message_type"`
	SubType     string        `json:"sub_type"`
	MessageID   ID            `json:"message_id"`
	GroupID     ID            `json:"group_id"`
	UserID      ID            `json:"user_id"`
	Message     message.Chain `json:"message"`
	RawMessage  string        `json:"raw_message"`
	Sender      GroupSender   `json:"sender"`
}

// UnknownEvent 表示当前 SDK 未细分的事件。
type UnknownEvent struct {
	Base
}

// RawEvent 表示无法强类型解析时保留的原始事件。
type RawEvent struct {
	Base
	parseErr error
}

// ParseError 返回强类型解析失败的原始错误。
func (e *RawEvent) ParseError() error { return e.parseErr }

type envelope struct {
	Time        int64  `json:"time"`
	PostType    string `json:"post_type"`
	SelfID      ID     `json:"self_id"`
	MessageType string `json:"message_type"`
}

type rawEnvelope struct {
	Time        json.RawMessage `json:"time"`
	PostType    json.RawMessage `json:"post_type"`
	SelfID      json.RawMessage `json:"self_id"`
	MessageType json.RawMessage `json:"message_type"`
}

// Parse 将 OneBot 事件 JSON 解析为具体事件类型，失败时返回 RawEvent。
func Parse(data []byte) Event {
	env, err := parseEnvelope(data)
	if err != nil {
		return &RawEvent{
			Base: Base{
				EventTime:     env.Time,
				EventPostType: env.PostType,
				EventSelfID:   env.SelfID,
				rawData:       bytes.Clone(data),
			},
			parseErr: fmt.Errorf("解析事件 envelope 失败: %w", err),
		}
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
			return &RawEvent{Base: base, parseErr: fmt.Errorf("解析私聊消息失败: %w", err)}
		}
		ev.Base = base
		return &ev
	case env.PostType == "message" && env.MessageType == "group":
		var ev GroupMessage
		if err := json.Unmarshal(data, &ev); err != nil {
			return &RawEvent{Base: base, parseErr: fmt.Errorf("解析群消息失败: %w", err)}
		}
		ev.Base = base
		return &ev
	default:
		return &UnknownEvent{Base: base}
	}
}

func parseEnvelope(data []byte) (envelope, error) {
	var raw rawEnvelope
	if err := json.Unmarshal(data, &raw); err != nil {
		return envelope{}, err
	}

	var env envelope
	var errs []error
	if len(raw.Time) > 0 {
		if err := json.Unmarshal(raw.Time, &env.Time); err != nil {
			errs = append(errs, fmt.Errorf("time: %w", err))
		}
	}
	if len(raw.PostType) > 0 {
		if err := json.Unmarshal(raw.PostType, &env.PostType); err != nil {
			errs = append(errs, fmt.Errorf("post_type: %w", err))
		}
	}
	if len(raw.SelfID) > 0 {
		if err := json.Unmarshal(raw.SelfID, &env.SelfID); err != nil {
			errs = append(errs, fmt.Errorf("self_id: %w", err))
		}
	}
	if len(raw.MessageType) > 0 {
		if err := json.Unmarshal(raw.MessageType, &env.MessageType); err != nil {
			errs = append(errs, fmt.Errorf("message_type: %w", err))
		}
	}
	return env, errors.Join(errs...)
}
