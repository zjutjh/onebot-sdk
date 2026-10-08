package main

import (
	"context"
	"os"
	"os/signal"
	"strconv"

	"github.com/zjutjh/onebot-sdk"
	"github.com/zjutjh/onebot-sdk/api"
	"github.com/zjutjh/onebot-sdk/event"
	"github.com/zjutjh/onebot-sdk/message"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	pong, err := api.NewOB11Message(message.Text("pong"))
	if err != nil {
		panic(err)
	}

	err = onebot.ServeReverseWebSocket(ctx, ":8080", func(client *onebot.Client) {
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
	}, onebot.WithToken(os.Getenv("ONEBOT_TOKEN")))
	if err != nil {
		panic(err)
	}
}
