package message

import (
	"bytes"
	"fmt"
	"strconv"

	json "github.com/zjutjh/onebot-sdk/internal/jsonx"
)

// ID 是兼容 JSON 整数和十进制整数字符串的 OneBot ID。
// message_id 等字段在不同实现间数值语义不同(NapCat 为 31 位正数,
// SnowLuma 为有符号 int32),ID 只承诺非零时可用于等值匹配,不承诺大小关系。
type ID int64

// UnmarshalJSON 解码 int64 范围内的整数 ID。
func (id *ID) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if len(data) == 0 {
		return fmt.Errorf("无效 ID: 空 JSON")
	}
	if data[0] != '"' {
		if data[0] != '-' && (data[0] < '0' || data[0] > '9') {
			return fmt.Errorf("无效 ID %q: 必须是整数或十进制整数字符串", data)
		}
		var value int64
		if err := json.Unmarshal(data, &value); err != nil {
			return fmt.Errorf("无效 ID %q: %w", data, err)
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
		return fmt.Errorf("无效 ID %q: %w", text, err)
	}
	*id = ID(value)
	return nil
}

// Int64 返回 ID 的 int64 值。
func (id ID) Int64() int64 { return int64(id) }

// StrNum 是兼容 JSON 字符串和整数的字符串字段,如 at.qq("all" 或 QQ 号)和 face.id。
type StrNum string

// UnmarshalJSON 接受字符串或整数字面量,统一解码为十进制字符串。
func (s *StrNum) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if len(data) == 0 {
		return fmt.Errorf("无效字符串字段: 空 JSON")
	}
	if data[0] == '"' {
		var text string
		if err := json.Unmarshal(data, &text); err != nil {
			return err
		}
		*s = StrNum(text)
		return nil
	}
	if value, err := strconv.ParseInt(string(data), 10, 64); err == nil {
		*s = StrNum(strconv.FormatInt(value, 10))
		return nil
	}
	return fmt.Errorf("无效字符串字段 %q: 必须是字符串或整数", data)
}
