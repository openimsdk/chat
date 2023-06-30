package sms

import (
	"context"
	"fmt"
	"github.com/OpenIMSDK/chat/pkg/common/config"
	"strings"
)

func New() (SMS, error) {
	switch strings.ToLower(config.Config.VerifyCode.Use) {
	case "":
		return empty{}, nil
	case "ali":
		return newAli()
	default:
		return nil, fmt.Errorf("not support sms: `%s`", config.Config.VerifyCode.Use)
	}
}

type SMS interface {
	Name() string
	SendCode(ctx context.Context, areaCode string, phoneNumber string, verifyCode string) error
}

type empty struct{}

func (e empty) Name() string {
	return "empty-sms"
}

func (e empty) SendCode(ctx context.Context, areaCode string, phoneNumber string, verifyCode string) error {
	return nil
}
