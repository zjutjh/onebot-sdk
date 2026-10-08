package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/zjutjh/onebot-sdk"
	"github.com/zjutjh/onebot-sdk/api"
	"github.com/zjutjh/onebot-sdk/event"
	"github.com/zjutjh/onebot-sdk/message"
)

func main() {
	ctx := context.Background()
	wsURL := os.Getenv("ONEBOT_WS_URL")
	if wsURL == "" {
		wsURL = "ws://127.0.0.1:3001"
	}

	client, err := onebot.DialWebSocket(ctx, wsURL, onebot.WithToken(os.Getenv("ONEBOT_TOKEN")))
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// 后端自动检测:打印当前连接的方言(NapCat/SnowLuma/Generic)
	fmt.Println("后端方言:", client.Dialect(ctx).Name)

	pong, err := api.NewOB11Message(message.Text("pong"))
	if err != nil {
		panic(err)
	}
	for ev := range client.Events() {
		switch e := ev.(type) {
		case *event.PrivateMessage:
			if e.Message.Text() == "/ping" {
				userID := strconv.FormatInt(e.UserID.Int64(), 10)
				_, _ = client.API().SendPrivateMsg(ctx, api.SendPrivateMsgRequest{
					UserID:  &userID,
					Message: pong,
				})
			}
		}
	}
	if err := client.Err(); err != nil && !errors.Is(err, onebot.ErrClosed) {
		panic(err)
	}
}
