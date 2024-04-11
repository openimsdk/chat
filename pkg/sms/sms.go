package sms

import "context"

type SMS interface {
	Name() string
	SendCode(ctx context.Context, areaCode string, phoneNumber string, verifyCode string) error
}
