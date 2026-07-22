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
	userID := os.Getenv("NAPCAT_TARGET_USER_ID")
	if userID == "" {
		panic("请设置 NAPCAT_TARGET_USER_ID")
	}

	client := napcat.NewHTTPClient("http://127.0.0.1:3000", napcat.WithToken(token))

	login, err := client.API().GetLoginInfo(ctx, api.GetLoginInfoRequest{})
	if err != nil {
		panic(err)
	}
	fmt.Println("当前账号:", login.UserID)

	msg, err := api.NewOB11Message(message.Text("来自 NapCat Go SDK 的 HTTP 消息"))
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
