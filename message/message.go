// Package message 提供 OneBot 消息段和消息链类型。
package message

import (
	"fmt"
	"strconv"
	"strings"
)

// Segment 表示一个 OneBot 消息段。
type Segment struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}

// Chain 表示按顺序排列的消息段。
type Chain []Segment

// ChainOf 从消息段构造消息链。
func ChainOf(segments ...Segment) Chain {
	out := make(Chain, len(segments))
	copy(out, segments)
	return out
}

// Text 构造文本消息段。
func Text(text string) Segment {
	return Segment{Type: "text", Data: map[string]any{"text": text}}
}

// At 构造 @ 指定 QQ 的消息段。
func At(qq int64) Segment {
	return Segment{Type: "at", Data: map[string]any{"qq": strconv.FormatInt(qq, 10)}}
}

// AtAll 构造 @全体成员 的消息段。
func AtAll() Segment {
	return Segment{Type: "at", Data: map[string]any{"qq": "all"}}
}

// Image 构造图片消息段。
func Image(file string) Segment {
	return Segment{Type: "image", Data: map[string]any{"file": file}}
}

// Reply 构造回复消息段。
func Reply(id int64) Segment {
	return Segment{Type: "reply", Data: map[string]any{"id": strconv.FormatInt(id, 10)}}
}

// Record 构造语音消息段。
func Record(file string) Segment {
	return Segment{Type: "record", Data: map[string]any{"file": file}}
}

// Video 构造视频消息段。
func Video(file string) Segment {
	return Segment{Type: "video", Data: map[string]any{"file": file}}
}

// Face 构造 QQ 表情消息段。
func Face(id string) Segment {
	return Segment{Type: "face", Data: map[string]any{"id": id}}
}

// NodeCustom 构造自定义合并转发节点。
func NodeCustom(name string, uin int64, content ...Segment) Segment {
	return Segment{Type: "node", Data: map[string]any{
		"name":    name,
		"uin":     strconv.FormatInt(uin, 10),
		"content": ChainOf(content...),
	}}
}

// String 读取消息段 data 中指定字段的字符串形式。
func (s Segment) String(key string) string {
	switch data := s.Data.(type) {
	case map[string]any:
		return stringify(data[key])
	case map[string]string:
		return data[key]
	default:
		return ""
	}
}

// Text 提取消息链中所有 text 段的文本并拼接。
func (c Chain) Text() string {
	var b strings.Builder
	for _, seg := range c {
		if seg.Type == "text" {
			b.WriteString(seg.String("text"))
		}
	}
	return b.String()
}

// OfType 返回指定类型的消息段。
func (c Chain) OfType(segmentType string) Chain {
	var out Chain
	for _, seg := range c {
		if seg.Type == segmentType {
			out = append(out, seg)
		}
	}
	return out
}

func stringify(v any) string {
	if v == nil {
		return ""
	}
	return fmt.Sprint(v)
}
