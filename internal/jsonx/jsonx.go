// Package jsonx 提供与 encoding/json 默认行为一致的 Sonic 编解码入口。
package jsonx

import (
	stdjson "encoding/json"

	"github.com/bytedance/sonic"
)

// RawMessage 延迟解析一段 JSON。
type RawMessage = stdjson.RawMessage

// Marshal 使用 Sonic 的标准库兼容配置编码 JSON。
func Marshal(value any) ([]byte, error) {
	return sonic.ConfigStd.Marshal(value)
}

// Unmarshal 使用 Sonic 的标准库兼容配置解码 JSON。
func Unmarshal(data []byte, value any) error {
	return sonic.ConfigStd.Unmarshal(data, value)
}
