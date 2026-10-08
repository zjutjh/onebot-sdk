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
	token := os.Getenv("ONEBOT_TOKEN")
	userID := os.Getenv("ONEBOT_TARGET_USER_ID")
	if userID == "" {
		panic("请设置 ONEBOT_TARGET_USER_ID")
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
		UserID:  &userID,
		Message: msg,
	})
	if err != nil {
		panic(err)
	}
}
