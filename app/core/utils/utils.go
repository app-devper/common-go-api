package utils

import (
	"context"
	"crypto/rand"
	"devper/app/core/notify"
	"errors"
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
		return res, err
	}
	return res, nil
}

func ToFormat(date time.Time) string {
	location, _ := time.LoadLocation("Asia/Bangkok")
	format := "02 Jan 2006 15:04"
	return date.In(location).Format(format)
}

const otpChars = "1234567890"

func GenerateCode(length int) (string, error) {
	buffer := make([]byte, length)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}
	otpCharsLength := len(otpChars)
	for i := 0; i < length; i++ {
		buffer[i] = otpChars[int(buffer[i])%otpCharsLength]
	}
	return string(buffer), nil
}

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func GenerateRefId(length int) (string, error) {
	buffer := make([]byte, length)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}
	otpCharsLength := len(alphabet)
	for i := 0; i < length; i++ {
		buffer[i] = alphabet[int(buffer[i])%otpCharsLength]
	}
	return string(buffer), nil
}
