package core

import (
	"context"
	"devper/app/core/notify"
	"errors"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

func InitContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	return ctx, cancel
}

func NotifyMassage(massage string) (*notify.Response, error) {
	token := os.Getenv("LINE_TOKEN")
	if token == "" {
		err := errors.New("line token empty")
		return nil, err
	}
	c := notify.NewClient()
	res, err := c.NotifyMessage(context.Background(), token, massage)
	if err != nil {
		logrus.Error(err)
		return res, err
	}
	return res, nil
}
