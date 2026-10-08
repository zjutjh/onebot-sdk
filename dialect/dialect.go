// Package dialect 描述 OneBot 11 各实现的协议方言差异。
//
// 已知方言差异以布尔位表达,调用方据此选择补拉路径或参数裁剪;
// 未识别的后端统一落到 Generic(纯 OneBot 11 核心能力,零补拉零扩展)。
package dialect

import (
	"context"
	"errors"
	"fmt"

	"github.com/zjutjh/onebot-sdk/transport"
)

// ErrUnknownBackend 表示 app_name 不在方言表内。
// 这是稳定判定而非临时故障,调用方可以放心缓存由此得出的 Generic。
var ErrUnknownBackend = errors.New("未知后端")

// Dialect 描述单个 OneBot 实现的行为差异。
type Dialect struct {
	// Name 是方言标识,用于日志和错误信息,不为空。
	Name string
	// ForwardFetch 表示事件内 forward 段仅有 id,需要补拉 get_forward_msg 填充内容。
	ForwardFetch bool
	// ReactionsViaGetMsg 表示消息表情回应通过 get_msg 的 emoji_likes_list 读取。
	ReactionsViaGetMsg bool
	// EmojiLikesAction 表示消息表情回应通过独立 action get_msg_emoji_likes 读取。
	EmojiLikesAction bool
	// FriendRemark 表示 set_friend_add_request 支持附加 remark 参数。
	FriendRemark bool
}

// 预置方言表。新增实现时在此登记,不要在调用方散落字符串判断。
var (
	// NapCat 内嵌 forward 内容无需补拉,表情回应挂在 get_msg,加好友支持备注。
	NapCat = Dialect{
		Name:               "napcat",
		ReactionsViaGetMsg: true,
		FriendRemark:       true,
	}
	// SnowLuma 事件内 forward 仅有 resId,需要补拉;表情回应有独立 action。
	SnowLuma = Dialect{
		Name:             "snowluma",
		ForwardFetch:     true,
		EmojiLikesAction: true,
	}
	// Generic 是未识别后端的保守形态:只用 OneBot 11 核心能力。
	Generic = Dialect{Name: "generic"}
)

// Detect 调用 get_version_info,按 app_name 识别后端方言。
// 调用失败时返回 Generic 和错误;app_name 未知时返回 Generic 和 ErrUnknownBackend。
// 调用方据此区分"临时故障(可重试)"与"稳定判定(可缓存)"。
func Detect(ctx context.Context, caller transport.Caller) (Dialect, error) {
	var res struct {
		AppName string `json:"app_name"`
	}
	if err := caller.Call(ctx, "get_version_info", nil, &res); err != nil {
		return Generic, fmt.Errorf("方言检测失败: %w", err)
	}
	switch res.AppName {
	case "NapCat.Onebot":
		return NapCat, nil
	case "SnowLuma":
		return SnowLuma, nil
	default:
		return Generic, fmt.Errorf("%w: app_name %q", ErrUnknownBackend, res.AppName)
	}
}
