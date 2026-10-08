// Package message 提供强类型 OneBot 消息段、消息链的解析与构造。
package message

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	json "github.com/zjutjh/onebot-sdk/internal/jsonx"
)

// SegmentData 是消息段 data 的强类型载体。
// 实现方都是值类型;未识别的段类型由 UnknownSegmentData 保留原始字段。
type SegmentData interface {
	isSegmentData()
}

// TextData 是文本段。
type TextData struct {
	Text string `json:"text"`
}

// AtData 是 @ 段,QQ 为 "all" 或十进制 QQ 号字符串。
type AtData struct {
	QQ   StrNum `json:"qq"`
	Name string `json:"name,omitempty"`
}

// FaceData 是 QQ 表情段。
type FaceData struct {
	ID StrNum `json:"id"`
}

// ReplyData 是回复段。
type ReplyData struct {
	ID ID `json:"id"`
}

// ImageData 是图片段,字段取两后端并集,缺席为零值。
type ImageData struct {
	File    string `json:"file"`
	SubType StrNum `json:"sub_type,omitempty"`
	URL     string `json:"url,omitempty"`
	Summary string `json:"summary,omitempty"`
}

// RecordData 是语音段。
type RecordData struct {
	File string `json:"file"`
	URL  string `json:"url,omitempty"`
}

// VideoData 是视频段。
type VideoData struct {
	File string `json:"file"`
	URL  string `json:"url,omitempty"`
}

// FileData 是文件段(发送接口使用)。
type FileData struct {
	FileID string `json:"file_id"`
}

// JSONData 是 JSON 消息段。
type JSONData struct {
	Data string `json:"data"`
}

// MFaceData 是商城表情段,字段形态与 NapCat 上报一致;
// SnowLuma 在上报前已把商城表情转换为图片段,不会出现该类型。
type MFaceData struct {
	EmojiPackageID int64  `json:"emoji_package_id"`
	EmojiID        string `json:"emoji_id"`
	Key            string `json:"key"`
	Summary        string `json:"summary,omitempty"`
}

// CardData 是卡片消息段,内层 data 形态随卡片类型变化(对象或字符串),
// 按原始 JSON 保留,不做二次解释。
type CardData struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data,omitempty"`
}

// ForwardData 是合并转发段。
// NapCat 事件内已内嵌 Content;SnowLuma 仅上报 ID,
// 由归一化管线补拉填充 Content,补拉失败时 Content 为空且 FetchErr 可见。
type ForwardData struct {
	ID       string        `json:"id"`
	Content  []ForwardNode `json:"-"`
	FetchErr error         `json:"-"`
}

// UnmarshalJSON 接受两种上报形态:NapCat 在 id 旁内嵌 content 节点数组,
// SnowLuma 仅有 id。FetchErr 不参与反序列化。
func (d *ForwardData) UnmarshalJSON(data []byte) error {
	var raw struct {
		ID      string          `json:"id"`
		Content json.RawMessage `json:"content"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("解析 forward 段 data 失败: %w", err)
	}
	d.ID = raw.ID
	if json.IsNull(raw.Content) {
		d.Content = nil
		return nil
	}
	if err := json.Unmarshal(raw.Content, &d.Content); err != nil {
		return fmt.Errorf("解析 forward 段 content 失败: %w", err)
	}
	return nil
}

// MarshalJSON 编码回 wire 形态,仅在存在节点时携带 content。
func (d ForwardData) MarshalJSON() ([]byte, error) {
	if len(d.Content) == 0 {
		return json.Marshal(struct {
			ID string `json:"id"`
		}{ID: d.ID})
	}
	return json.Marshal(struct {
		ID      string        `json:"id"`
		Content []ForwardNode `json:"content"`
	}{ID: d.ID, Content: d.Content})
}

// UnknownSegmentData 保留未识别段类型的原始字段。
type UnknownSegmentData struct {
	Fields map[string]any
}

func (TextData) isSegmentData()           {}
func (AtData) isSegmentData()             {}
func (FaceData) isSegmentData()           {}
func (ReplyData) isSegmentData()          {}
func (ImageData) isSegmentData()          {}
func (RecordData) isSegmentData()         {}
func (VideoData) isSegmentData()          {}
func (FileData) isSegmentData()           {}
func (JSONData) isSegmentData()           {}
func (MFaceData) isSegmentData()          {}
func (CardData) isSegmentData()           {}
func (ForwardData) isSegmentData()        {}
func (UnknownSegmentData) isSegmentData() {}

// MarshalJSON 把保留字段编码为段 data 对象。
func (d UnknownSegmentData) MarshalJSON() ([]byte, error) {
	if d.Fields == nil {
		return []byte("{}"), nil
	}
	return json.Marshal(d.Fields)
}

// Segment 表示一个 OneBot 消息段。
type Segment struct {
	Type string      `json:"type"`
	Data SegmentData `json:"data"`
}

// UnmarshalJSON 按 type 分派到强类型 data,未知类型保留原始字段。
// data 缺失或为 null 时使用该类型的零值。
func (s *Segment) UnmarshalJSON(data []byte) error {
	var raw struct {
		Type string          `json:"type"`
		Data json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("解析消息段失败: %w", err)
	}
	if raw.Type == "" {
		return fmt.Errorf("消息段缺少 type: %s", data)
	}
	s.Type = raw.Type

	empty := json.IsNull(raw.Data)
	switch raw.Type {
	case "text":
		return decodeInto(s, TextData{}, raw.Data, empty)
	case "at":
		return decodeInto(s, AtData{}, raw.Data, empty)
	case "face":
		return decodeInto(s, FaceData{}, raw.Data, empty)
	case "reply":
		return decodeInto(s, ReplyData{}, raw.Data, empty)
	case "image":
		return decodeInto(s, ImageData{}, raw.Data, empty)
	case "record":
		return decodeInto(s, RecordData{}, raw.Data, empty)
	case "video":
		return decodeInto(s, VideoData{}, raw.Data, empty)
	case "file":
		return decodeInto(s, FileData{}, raw.Data, empty)
	case "json":
		return decodeInto(s, JSONData{}, raw.Data, empty)
	case "mface":
		return decodeInto(s, MFaceData{}, raw.Data, empty)
	case "card":
		return decodeInto(s, CardData{}, raw.Data, empty)
	case "forward":
		return decodeInto(s, ForwardData{}, raw.Data, empty)
	default:
		var fields map[string]any
		if !empty {
			if err := json.Unmarshal(raw.Data, &fields); err != nil {
				return fmt.Errorf("解析 %s 段 data 失败: %w", raw.Type, err)
			}
		}
		s.Data = UnknownSegmentData{Fields: fields}
		return nil
	}
}

func decodeInto[T SegmentData](s *Segment, target T, raw json.RawMessage, empty bool) error {
	if empty {
		s.Data = target
		return nil
	}
	// target 是值类型,取地址反序列化后仍以值形态存入,保证接口断言按值匹配
	if err := json.Unmarshal(raw, any(&target)); err != nil {
		return fmt.Errorf("解析 %s 段 data 失败: %w", s.Type, err)
	}
	s.Data = target
	return nil
}

// Chain 表示按顺序排列的消息段。
type Chain []Segment

// UnmarshalJSON 接受段数组或字符串形态;
// 字符串形态按 OneBot 约定解码为单个 text 段。
func (c *Chain) UnmarshalJSON(data []byte) error {
	trimmed := bytes.TrimSpace(data)
	if len(trimmed) > 0 && trimmed[0] == '"' {
		var text string
		if err := json.Unmarshal(trimmed, &text); err != nil {
			return fmt.Errorf("解析消息链失败: %w", err)
		}
		*c = Chain{{Type: "text", Data: TextData{Text: text}}}
		return nil
	}
	var segs []Segment
	if err := json.Unmarshal(data, &segs); err != nil {
		return fmt.Errorf("解析消息链失败: %w", err)
	}
	*c = segs
	return nil
}

// ChainOf 从消息段构造消息链。
func ChainOf(segments ...Segment) Chain {
	out := make(Chain, len(segments))
	copy(out, segments)
	return out
}

// Text 构造文本消息段。
func Text(text string) Segment {
	return Segment{Type: "text", Data: TextData{Text: text}}
}

// At 构造 @ 指定 QQ 的消息段。
func At(qq int64) Segment {
	return Segment{Type: "at", Data: AtData{QQ: StrNum(strconv.FormatInt(qq, 10))}}
}

// AtAll 构造 @全体成员 的消息段。
func AtAll() Segment {
	return Segment{Type: "at", Data: AtData{QQ: "all"}}
}

// Image 构造图片消息段。
func Image(file string) Segment {
	return Segment{Type: "image", Data: ImageData{File: file}}
}

// Reply 构造回复消息段。
func Reply(id int64) Segment {
	return Segment{Type: "reply", Data: ReplyData{ID: ID(id)}}
}

// Record 构造语音消息段。
func Record(file string) Segment {
	return Segment{Type: "record", Data: RecordData{File: file}}
}

// Video 构造视频消息段。
func Video(file string) Segment {
	return Segment{Type: "video", Data: VideoData{File: file}}
}

// Face 构造 QQ 表情消息段。
func Face(id string) Segment {
	return Segment{Type: "face", Data: FaceData{ID: StrNum(id)}}
}

// JSON 构造 JSON 消息段。
func JSON(data string) Segment {
	return Segment{Type: "json", Data: JSONData{Data: data}}
}

// MFace 构造商城表情消息段。
func MFace(packageID int64, emojiID string) Segment {
	return Segment{Type: "mface", Data: MFaceData{EmojiPackageID: packageID, EmojiID: emojiID}}
}

// Card 构造卡片消息段,cardType 为卡片类型,data 为卡片内容的原始 JSON 文本。
func Card(cardType, data string) Segment {
	return Segment{Type: "card", Data: CardData{Type: cardType, Data: json.RawMessage(data)}}
}

// Forward 构造合并转发消息段,id 为转发资源标识。
func Forward(id string) Segment {
	return Segment{Type: "forward", Data: ForwardData{ID: id}}
}

// File 构造文件消息段。
func File(fileID string) Segment {
	return Segment{Type: "file", Data: FileData{FileID: fileID}}
}

// NodeCustom 构造自定义合并转发节点。
func NodeCustom(name string, uin int64, content ...Segment) Segment {
	return Segment{Type: "node", Data: UnknownSegmentData{Fields: map[string]any{
		"name":    name,
		"uin":     strconv.FormatInt(uin, 10),
		"content": ChainOf(content...),
	}}}
}

// NodeID 构造引用已有消息的合并转发节点。
func NodeID(id string) Segment {
	return Segment{Type: "node", Data: UnknownSegmentData{Fields: map[string]any{"id": id}}}
}

// Text 提取消息链中所有 text 段的文本并拼接。
func (c Chain) Text() string {
	var b strings.Builder
	for _, seg := range c {
		if d, ok := seg.Data.(TextData); ok {
			b.WriteString(d.Text)
		}
	}
	return b.String()
}
