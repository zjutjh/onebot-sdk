package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/zjutjh/onebot-sdk"
	"github.com/zjutjh/onebot-sdk/api"
	"github.com/zjutjh/onebot-sdk/message"
)

func main() {
	ctx := context.Background()
	token := os.Getenv("ONEBOT_TOKEN")
	raw := os.Getenv("ONEBOT_TARGET_USER_ID")
	if raw == "" {
		panic("请设置 ONEBOT_TARGET_USER_ID")
	}
	userID, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		panic(err)
	}

	client := onebot.NewHTTPClient("http://127.0.0.1:3000", onebot.WithToken(token))

	login, err := client.API().GetLoginInfo(ctx, api.GetLoginInfoRequest{})
	if err != nil {
		panic(err)
	}
	fmt.Println("当前账号:", login.UserID)

	msg, err := api.NewOB11Message(message.Text("来自 onebot-sdk 的 HTTP 消息"))
	if err != nil {
		panic(err)
	}
	_, err = client.API().SendPrivateMsg(ctx, api.SendPrivateMsgRequest{
		UserID:  message.ID(userID),
		Message: msg,
	})
	if err != nil {
		panic(err)
	}
}
