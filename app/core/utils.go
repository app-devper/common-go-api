package core

import (
	"context"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

func InitContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	return ctx, cancel
}

func PushMassage(massage string) {
	secret := os.Getenv("LINE_SECRET")
	accessToken := os.Getenv("LINE_ACCESS_TOKEN")
	to := os.Getenv("LINE_USER_ID")
	bot, err := linebot.New(secret, accessToken)
	if bot != nil {
		_, err = bot.PushMessage(to, linebot.NewTextMessage(massage)).Do()
		if err != nil {
			logrus.Error(err)
		}
	}
}
