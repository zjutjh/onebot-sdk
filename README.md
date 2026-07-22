# NapCat Go SDK

这是一个面向 NapCat/OneBot 11 的纯 Go SDK。SDK 以 NapCat 官方 OpenAPI 4.18.13 为来源生成强类型 API 方法，并提供 HTTP、正向 WebSocket、反向 WebSocket server、事件解析和消息段构造能力。

## 特性

- 基于 NapCat 4.18.13 OpenAPI 生成 170 个 action 绑定。
- HTTP 调用基于 Resty。
- WebSocket 调用基于 coder/websocket，支持 echo 匹配和事件流。
- JSON 编解码使用 Sonic 的标准库兼容配置。
- 手写 OneBot 消息段和事件解析，不内置插件、命令或路由框架。

## 安装

```sh
go get github.com/zjutjh/napcat-sdk
```

## HTTP 调用

```go
package main

import (
	"context"
	"fmt"
	"os"

	napcat "github.com/zjutjh/napcat-sdk"
	"github.com/zjutjh/napcat-sdk/api"
	"github.com/zjutjh/napcat-sdk/message"
)

func main() {
	ctx := context.Background()
	token := os.Getenv("NAPCAT_TOKEN")
	userID := "123456"

	client := napcat.NewHTTPClient("http://127.0.0.1:3000", napcat.WithToken(token))

	login, err := client.API().GetLoginInfo(ctx, api.GetLoginInfoRequest{})
	if err != nil {
		panic(err)
	}
	fmt.Println("当前账号:", login.UserID)

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
	"os"
	"strconv"

	napcat "github.com/zjutjh/napcat-sdk"
	"github.com/zjutjh/napcat-sdk/api"
	"github.com/zjutjh/napcat-sdk/event"
	"github.com/zjutjh/napcat-sdk/message"
)

func main() {
	ctx := context.Background()
	client, err := napcat.DialWebSocket(ctx, "ws://127.0.0.1:3001", napcat.WithToken(os.Getenv("NAPCAT_TOKEN")))
	if err != nil {
		panic(err)
	}
	defer client.Close()

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
			userID := strconv.FormatInt(e.UserID, 10)
			_, _ = client.API().SendPrivateMsg(ctx, api.SendPrivateMsgRequest{
				UserID:  &userID,
				Message: pong,
			})
		}
	}
}
```

## 原始 action 兜底

当 NapCat 新增接口但 SDK 还未生成强类型方法时，可以直接调用原始 action：

```go
var result map[string]any
err := client.Call(ctx, "some_new_action", map[string]any{"value": 1}, &result)
```

## 重新生成 API

```sh
go generate ./...
```

生成器读取 `internal/openapi/4.18.13/openapi.json`，把 action、client 和所有模型输出到 `api/`。生成文件头会标明“请勿手动修改”。
