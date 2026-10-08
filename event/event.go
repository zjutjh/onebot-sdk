// Package event 提供 OneBot 上报事件的解析和强类型。
//
// 事件字段取已实现后端(NapCat、SnowLuma)上报的并集,单个后端缺席的字段为零值;
// 未识别事件通过 UnknownEvent 保留原始 JSON,强类型解析失败通过 RawEvent 交付。
// 私聊 target_id、forward 段内容等方言缺口由根包归一化管线补齐,不在本包处理。
package event

import (
	"bytes"
	"errors"
	"fmt"

	json "github.com/zjutjh/onebot-sdk/internal/jsonx"
	"github.com/zjutjh/onebot-sdk/message"
)

// ID 是兼容 JSON 整数和十进制整数字符串的 OneBot ID,实现复用 message.ID。
type ID = message.ID

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
// 临时私聊携带来源群,GroupID 仅在该场景非零。
type PrivateSender struct {
	UserID   ID     `json:"user_id"`
	Nickname string `json:"nickname"`
	GroupID  ID     `json:"group_id,omitempty"`
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

// PrivateMessage 表示私聊消息事件,post_type 为 message 或 message_sent(自发)。
// target_id 恒为会话对端:NapCat 两个方向都上报(收到的私聊里对端就是发送者,与
// user_id 相等),SnowLuma 只在自发时上报,缺失时归一化管线按 user_id 补齐。
type PrivateMessage struct {
	Base
	MessageType string        `json:"message_type"`
	SubType     string        `json:"sub_type"`
	MessageID   ID            `json:"message_id"`
	MessageSeq  int64         `json:"message_seq"`
	UserID      ID            `json:"user_id"`
	TargetID    ID            `json:"target_id,omitempty"`
	Message     message.Chain `json:"message"`
	RawMessage  string        `json:"raw_message"`
	Font        int           `json:"font"`
	Sender      PrivateSender `json:"sender"`
}

// GroupMessage 表示群消息事件,post_type 为 message 或 message_sent(自发)。
type GroupMessage struct {
	Base
	MessageType string        `json:"message_type"`
	SubType     string        `json:"sub_type"`
	MessageID   ID            `json:"message_id"`
	MessageSeq  int64         `json:"message_seq"`
	GroupID     ID            `json:"group_id"`
	GroupName   string        `json:"group_name,omitempty"`
	UserID      ID            `json:"user_id"`
	Message     message.Chain `json:"message"`
	RawMessage  string        `json:"raw_message"`
	Font        int           `json:"font"`
	Sender      GroupSender   `json:"sender"`
}

// FriendRecall 表示好友消息撤回。
type FriendRecall struct {
	Base
	UserID    ID `json:"user_id"`
	MessageID ID `json:"message_id"`
}

// GroupRecall 表示群消息撤回。
type GroupRecall struct {
	Base
	GroupID    ID `json:"group_id"`
	OperatorID ID `json:"operator_id"`
	UserID     ID `json:"user_id"`
	MessageID  ID `json:"message_id"`
}

// GroupIncrease 表示群成员入群,SubType 为 approve 或 invite。
type GroupIncrease struct {
	Base
	SubType    string `json:"sub_type"`
	GroupID    ID     `json:"group_id"`
	OperatorID ID     `json:"operator_id"`
	UserID     ID     `json:"user_id"`
}

// GroupDecrease 表示群成员退群,SubType 为 leave、kick、kick_me 或 disband。
type GroupDecrease struct {
	Base
	SubType    string `json:"sub_type"`
	GroupID    ID     `json:"group_id"`
	OperatorID ID     `json:"operator_id"`
	UserID     ID     `json:"user_id"`
}

// FriendAdd 表示好友添加完成。
type FriendAdd struct {
	Base
	UserID ID `json:"user_id"`
}

// GroupBan 表示群禁言,SubType 为 ban 或 lift_ban。
type GroupBan struct {
	Base
	SubType    string `json:"sub_type"`
	GroupID    ID     `json:"group_id"`
	OperatorID ID     `json:"operator_id"`
	UserID     ID     `json:"user_id"`
	Duration   int64  `json:"duration"`
}

// GroupAdmin 表示群管理员设置,SubType 为 set 或 unset。
type GroupAdmin struct {
	Base
	SubType string `json:"sub_type"`
	GroupID ID     `json:"group_id"`
	UserID  ID     `json:"user_id"`
}

// GroupFile 表示群文件上传的文件信息。
type GroupFile struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Size  int64  `json:"size"`
	BusID int64  `json:"busid"`
}

// GroupUpload 表示群文件上传。
type GroupUpload struct {
	Base
	GroupID ID        `json:"group_id"`
	UserID  ID        `json:"user_id"`
	File    GroupFile `json:"file"`
}

// Essence 表示群精华消息设置,SubType 为 add 或 delete。
type Essence struct {
	Base
	SubType    string `json:"sub_type"`
	GroupID    ID     `json:"group_id"`
	UserID     ID     `json:"user_id"`
	SenderID   ID     `json:"sender_id,omitempty"`
	OperatorID ID     `json:"operator_id"`
	MessageID  ID     `json:"message_id"`
	MessageSeq int64  `json:"message_seq,omitempty"`
}

// Poke 表示戳一戳(notify/poke),字段取群聊与好友私聊并集:
// 群聊携带 GroupID,好友私聊额外携带 SenderID 与 ActionImgURL,缺席为零值。
type Poke struct {
	Base
	GroupID      ID     `json:"group_id,omitempty"`
	UserID       ID     `json:"user_id"`
	SenderID     ID     `json:"sender_id,omitempty"`
	TargetID     ID     `json:"target_id,omitempty"`
	Action       string `json:"action,omitempty"`
	Suffix       string `json:"suffix,omitempty"`
	ActionImgURL string `json:"action_img_url,omitempty"`
}

// FriendRequest 表示加好友请求。两后端均只上报 user_id/comment/flag。
type FriendRequest struct {
	Base
	UserID  ID     `json:"user_id"`
	Comment string `json:"comment"`
	Flag    string `json:"flag"`
}

// GroupRequest 表示入群请求或邀请,SubType 为 add 或 invite。
type GroupRequest struct {
	Base
	SubType   string `json:"sub_type"`
	GroupID   ID     `json:"group_id"`
	UserID    ID     `json:"user_id"`
	InvitedID ID     `json:"invited_id,omitempty"`
	Comment   string `json:"comment"`
	Flag      string `json:"flag"`
}

// Heartbeat 表示心跳事件。
type Heartbeat struct {
	Base
	SubType  string         `json:"sub_type"`
	Interval int64          `json:"interval"`
	Status   map[string]any `json:"status,omitempty"`
}

// Lifecycle 表示连接生命周期事件,SubType 为 connect 或 enable。
type Lifecycle struct {
	Base
	SubType string `json:"sub_type"`
}

// UnknownEvent 表示当前 SDK 未细分的事件,原始 JSON 可通过 Raw() 读取。
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
	Time          int64  `json:"time"`
	PostType      string `json:"post_type"`
	SelfID        ID     `json:"self_id"`
	MessageType   string `json:"message_type"`
	NoticeType    string `json:"notice_type"`
	RequestType   string `json:"request_type"`
	MetaEventType string `json:"meta_event_type"`
	SubType       string `json:"sub_type"`
}

type rawEnvelope struct {
	Time          json.RawMessage `json:"time"`
	PostType      json.RawMessage `json:"post_type"`
	SelfID        json.RawMessage `json:"self_id"`
	MessageType   json.RawMessage `json:"message_type"`
	NoticeType    json.RawMessage `json:"notice_type"`
	RequestType   json.RawMessage `json:"request_type"`
	MetaEventType json.RawMessage `json:"meta_event_type"`
	SubType       json.RawMessage `json:"sub_type"`
}

// Parse 将 OneBot 事件 JSON 解析为具体事件类型,失败时返回 RawEvent。
func Parse(data []byte) Event {
	env, err := parseEnvelope(data)
	base := Base{
		EventTime:     env.Time,
		EventPostType: env.PostType,
		EventSelfID:   env.SelfID,
		rawData:       bytes.Clone(data),
	}
	if err != nil {
		return &RawEvent{Base: base, parseErr: fmt.Errorf("解析事件 envelope 失败: %w", err)}
	}

	switch {
	case (env.PostType == "message" || env.PostType == "message_sent") && env.MessageType == "private":
		return parseTyped(data, base, &PrivateMessage{}, "私聊消息")
	case (env.PostType == "message" || env.PostType == "message_sent") && env.MessageType == "group":
		return parseTyped(data, base, &GroupMessage{}, "群消息")
	case env.PostType == "notice":
		return parseNotice(data, base, env)
	case env.PostType == "request":
		return parseRequest(data, base, env)
	case env.PostType == "meta_event":
		return parseMetaEvent(data, base, env)
	default:
		return &UnknownEvent{Base: base}
	}
}

func parseNotice(data []byte, base Base, env envelope) Event {
	switch {
	case env.NoticeType == "friend_recall":
		return parseTyped(data, base, &FriendRecall{}, "好友撤回")
	case env.NoticeType == "group_recall":
		return parseTyped(data, base, &GroupRecall{}, "群撤回")
	case env.NoticeType == "group_increase":
		return parseTyped(data, base, &GroupIncrease{}, "入群")
	case env.NoticeType == "group_decrease":
		return parseTyped(data, base, &GroupDecrease{}, "退群")
	case env.NoticeType == "friend_add":
		return parseTyped(data, base, &FriendAdd{}, "好友添加")
	case env.NoticeType == "group_ban":
		return parseTyped(data, base, &GroupBan{}, "群禁言")
	case env.NoticeType == "group_admin":
		return parseTyped(data, base, &GroupAdmin{}, "群管理员")
	case env.NoticeType == "group_upload":
		return parseTyped(data, base, &GroupUpload{}, "群文件上传")
	case env.NoticeType == "essence":
		return parseTyped(data, base, &Essence{}, "精华消息")
	case env.NoticeType == "notify" && env.SubType == "poke":
		return parseTyped(data, base, &Poke{}, "戳一戳")
	default:
		return &UnknownEvent{Base: base}
	}
}

func parseRequest(data []byte, base Base, env envelope) Event {
	switch env.RequestType {
	case "friend":
		return parseTyped(data, base, &FriendRequest{}, "好友请求")
	case "group":
		return parseTyped(data, base, &GroupRequest{}, "群请求")
	default:
		return &UnknownEvent{Base: base}
	}
}

func parseMetaEvent(data []byte, base Base, env envelope) Event {
	switch env.MetaEventType {
	case "heartbeat":
		return parseTyped(data, base, &Heartbeat{}, "心跳")
	case "lifecycle":
		return parseTyped(data, base, &Lifecycle{}, "生命周期")
	default:
		return &UnknownEvent{Base: base}
	}
}

func parseTyped[E Event](data []byte, base Base, ev E, label string) Event {
	if err := json.Unmarshal(data, ev); err != nil {
		return &RawEvent{Base: base, parseErr: fmt.Errorf("解析%s事件失败: %w", label, err)}
	}
	// 重新注入 base:Base 的 rawData 不参与 JSON 反序列化。
	// 所有事件都内嵌 Base,提升的 setBase 方法覆盖全部类型,无需逐一维护。
	if b, ok := any(ev).(interface{ setBase(Base) }); ok {
		b.setBase(base)
	}
	return ev
}

// setBase 把已解析的 envelope 基础字段写回事件(Base.rawData 私有,无法通过 JSON 注入)。
func (b *Base) setBase(base Base) { *b = base }

func parseEnvelope(data []byte) (envelope, error) {
	var raw rawEnvelope
	if err := json.Unmarshal(data, &raw); err != nil {
		return envelope{}, err
	}

	var env envelope
	var errs []error
	assign := func(dst any, src json.RawMessage, label string) {
		if len(src) == 0 {
			return
		}
		if err := json.Unmarshal(src, dst); err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", label, err))
		}
	}
	assign(&env.Time, raw.Time, "time")
	assign(&env.PostType, raw.PostType, "post_type")
	assign(&env.SelfID, raw.SelfID, "self_id")
	assign(&env.MessageType, raw.MessageType, "message_type")
	assign(&env.NoticeType, raw.NoticeType, "notice_type")
	assign(&env.RequestType, raw.RequestType, "request_type")
	assign(&env.MetaEventType, raw.MetaEventType, "meta_event_type")
	assign(&env.SubType, raw.SubType, "sub_type")
	return env, errors.Join(errs...)
}
