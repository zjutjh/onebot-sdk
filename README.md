# onebot-sdk (Go)

这是一个面向 OneBot 11 的纯 Go 统一接入层 SDK，同时支持 **NapCat** 和 **SnowLuma** 两种后端。SDK 以 NapCat 官方 OpenAPI 4.18.33 为来源生成强类型 API 方法，方言差异（字段缺口、扩展 action、消息形态）在 SDK 内部归一：上游只见统一 API，无需关心后端是谁。提供 HTTP、正向 WebSocket、反向 WebSocket server、事件解析和消息段构造能力。

## 特性

- 基于方言表的双后端归一化：连接建立后自动检测后端（`get_version_info` 的 `app_name`），未知后端降级为通用 OneBot 11 模式（零补拉、零扩展）。
- 归一化事件流 `Events()`：强类型事件 + 类型化消息段 + 方言字段补齐 + 合并转发补拉。
- 原始事件旁路 `RawEvents()`：透传未解析的原始 JSON 帧，与 `Events()` 并存。
- 基于 NapCat 4.18.33 OpenAPI 生成 178 个 action 绑定；对 SnowLuma 调用不存在的 action 会得到携带 retcode 1404 的 `*onebot.APIError`（后端返回 `UNKNOWN_ACTION`）。
- HTTP 调用基于标准库 `net/http`；WebSocket 调用基于 coder/websocket，支持 echo 匹配、有界事件流和终止原因查询。
- JSON 编解码使用 Sonic 的标准库兼容配置。
- 不内置插件、命令或路由框架。

## 安装

```sh
go get github.com/zjutjh/onebot-sdk
```

（v2 起 module 路径由 `github.com/zjutjh/napcat-sdk` 改名为 `github.com/zjutjh/onebot-sdk`，属破坏性变更。）

## 后端与方言

| 行为 | NapCat | SnowLuma | SDK 归一化 |
|---|---|---|---|
| 方言检测 | 自报 `app_name` = `NapCat.Onebot` | 自报 `app_name` = `SnowLuma` | 自动检测；未知 → Generic |
| forward 段内容 | 事件内嵌 `data.content` | 仅上报 res_id | SnowLuma 由管线补拉 `get_forward_msg` 注入内容 |
| 表情回应读取 | `get_msg` 的 `emoji_likes_list`（`likes_cnt` 为字符串） | 独立 action `get_msg_emoji_likes` | `GetMessageReactions` 双路径归一为 `{EmojiID, Count}` |
| 私聊 `target_id` | 有，值为会话对端（收到的私聊等于发送者） | 仅自发消息上报 | 缺失时按对端补齐，统一为「会话对端」 |
| `set_friend_add_request` 的 remark | 支持 | 不支持 | SnowLuma 不传 remark |
| `message_id` 数值语义 | 31 位正数（MD5 清符号位） | 有符号 int32（SHA-1 截断） | 无法归一，见下方契约 |

检测时机：WebSocket（正向/反向）在连接建立时同步检测（超时 5 秒，失败降级 Generic 并记日志）；HTTP client 在首次需要方言时懒检测，检测失败不缓存、下次调用重试，只有识别出的后端和明确判定为未知后端这两种稳定结果才会缓存。`WithBackend(...)` 可显式指定后端跳过检测，取值在 `onebot.BackendNapCat`、`onebot.BackendSnowLuma`、`onebot.BackendGeneric` 三选一（默认 `BackendAuto` 自动检测）。

## message_id 契约

`message_id` 在两个后端的数值语义不同（NapCat 槽位含义与 SnowLuma 完全不可比），SDK 只保证：

- 它是不透明的外部定位符，**只能做等值比较**，禁止排序、大小比较或当水位用。
- 正常路径下保证非零；消息事件中缺失时 SDK 记日志警告，该消息无法用于回复、撤回或去重。
- NapCat 可能上报负数 ID 的边缘情况 SDK 照常接受，不假设递增。

## HTTP 调用

```go
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/zjutjh/onebot-sdk"
	"github.com/zjutjh/onebot-sdk/api"
	"github.com/zjutjh/onebot-sdk/message"
)

func main() {
	ctx := context.Background()
	userID := "123456"

	client := onebot.NewHTTPClient("http://127.0.0.1:3000", onebot.WithToken(os.Getenv("ONEBOT_TOKEN")))

	login, err := client.API().GetLoginInfo(ctx, api.GetLoginInfoRequest{})
	if err != nil {
		panic(err)
	}
	fmt.Println("当前账号:", login.UserID, "后端:", client.Dialect(ctx).Name)

	msg, err := api.NewOB11Message(message.Text("pong"))
	if err != nil {
		panic(err)
	}
	_, err = client.API().SendPrivateMsg(ctx, api.SendPrivateMsgRequest{
		UserID:  &userID,
		Message: msg,
	})
	if err != nil {
		panic(err)
	}
}
```

## 正向 WebSocket 事件循环

```go
package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/zjutjh/onebot-sdk"
	"github.com/zjutjh/onebot-sdk/api"
	"github.com/zjutjh/onebot-sdk/event"
	"github.com/zjutjh/onebot-sdk/message"
)

func main() {
	ctx := context.Background()
	client, err := onebot.DialWebSocket(
		ctx,
		"ws://127.0.0.1:3001",
		onebot.WithToken(os.Getenv("ONEBOT_TOKEN")),
		onebot.WithEventBuffer(1024),
		onebot.WithEventDeliveryTimeout(time.Second),
	)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	fmt.Println("后端方言:", client.Dialect(ctx).Name)

	pong, err := api.NewOB11Message(message.Text("pong"))
	if err != nil {
		panic(err)
	}
	for ev := range client.Events() {
		switch e := ev.(type) {
		case *event.PrivateMessage:
			if e.Message.Text() != "/ping" {
				continue
			}
			userID := strconv.FormatInt(e.UserID.Int64(), 10)
			_, _ = client.API().SendPrivateMsg(ctx, api.SendPrivateMsgRequest{
				UserID:  &userID,
				Message: pong,
			})
		}
	}
	if err := client.Err(); err != nil && !errors.Is(err, onebot.ErrClosed) {
		panic(err)
	}
}
```

### 事件管线与保序

事件流默认经过归一化管线：单 goroutine 串行处理（强类型解析 → 方言归一 → 投递），**投递顺序与到达顺序一致**。普通事件（不含 forward 段）零 API 调用快速通过；SnowLuma 的 forward 段触发补拉时有独立超时（10 秒），补拉失败照常投递（`ForwardData.Content` 为空、`FetchErr` 暴露错误），不伪造内容也不吞事件。

背压语义只有一跳（管线到用户）：

- 未配置 `WithEventDeliveryTimeout`（默认）：队列满时丢弃新事件并记日志（配置 logger 时）。
- 配置正数超时后：用户持续不消费超过该时长会以 `ErrEventBackpressure` 终止连接。**归一化补拉的 API 往返不计入交付超时**，不会因后端响应慢误杀连接。

传输层读循环永不阻塞：它是 action(echo) 响应唯一的路由方，阻塞会让在途调用（包括管线自己的补拉请求）永远等不到响应，所以传输到管线这一跳固定为非阻塞投递，队列满一律丢弃并记日志。事件队列有界。

事件 channel 关闭后可通过 `Client.Err()` 获取固定的首次终止原因：主动关闭及对端正常关闭返回 `ErrClosed`，网络读取失败返回 `TransportError`。HTTP client 的 `Err()` 始终返回 `nil`。

### 原始事件旁路

`RawEvents()` 返回原始事件 JSON 的旁路流，供绕过归一化管线自行解析（新事件类型、特殊字段、自管解析逻辑）：

```go
raw := client.RawEvents()
for frame := range raw {
	// frame 是未解析的原始 JSON 字节
}
```

与 `Events()` **并存而非二选一**：每个事件帧同时进入两者。注意：旁路在首次调用 `RawEvents()` 时才启用，**须在事件开始到达前调用**，之后才调用的帧不会被回放。旁路队列与事件队列同容量，满时丢弃旁路帧并记日志，不影响主事件流。HTTP client 返回 `nil`。

## 方言门面方法

少量方言差异通过手写门面方法归一（生成方法之外）：

```go
// 表情回应：NapCat 走 get_msg.emoji_likes_list，SnowLuma 走 get_msg_emoji_likes
reactions, err := client.GetMessageReactions(ctx, messageID)

// 合并转发内容补拉
nodes, err := client.GetForwardContent(ctx, resID)

// 好友申请处理：remark 仅 NapCat 支持，SnowLuma 下静默忽略
err = client.SetFriendRequest(ctx, flag, true, "备注名")
```

## 原始 action 兜底

当后端新增接口但 SDK 还未生成强类型方法时，可以直接调用原始 action：

```go
var result map[string]any
err := client.Call(ctx, "some_new_action", map[string]any{"value": 1}, &result)
```

## 按方言的 action 适用范围

178 个生成方法来自 NapCat 官方 OpenAPI：

- 标准 OneBot 11 核心接口（send_msg、get_group_member_info 等）两端都可用。
- NapCat 扩展接口对 SnowLuma 调用会得到携带 retcode 1404 的 `*onebot.APIError`（`UNKNOWN_ACTION`）——这是后端行为，不是 SDK 伪造的。
- SnowLuma 扩展 action 通过 `Call` 逃生口调用。

## 重新生成 API

```sh
go generate ./...
```

生成器读取 `internal/openapi/4.18.33/openapi.json`，把 action、client 和所有模型输出到 `api/`。生成文件头会标明“请勿手动修改”。
