package main

import (
	"context"
	"os"
	"os/signal"
	"strconv"

	napcat "github.com/zjutjh/napcat-sdk"
	"github.com/zjutjh/napcat-sdk/api"
	"github.com/zjutjh/napcat-sdk/event"
	"github.com/zjutjh/napcat-sdk/message"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	pong, err := api.NewOB11Message(message.Text("pong"))
	if err != nil {
		panic(err)
	}

	err = napcat.ServeReverseWebSocket(ctx, ":8080", func(client *napcat.Client) {
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
	}, napcat.WithToken(os.Getenv("NAPCAT_TOKEN")))
	if err != nil {
		panic(err)
	}
}
