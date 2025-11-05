package chat

import (
	"context"
	"time"
)

type Register struct {
	UserID      string    `bson:"user_id"`
	DeviceID    string    `bson:"device_id"`
	IP          string    `bson:"ip"`
	Platform    string    `bson:"platform"`
	AccountType string    `bson:"account_type"`
	Mode        string    `bson:"mode"`
	CreateTime  time.Time `bson:"create_time"`
}

func (Register) TableName() string {
	return "registers"
}

type RegisterInterface interface {
	// NewTx(tx any) RegisterInterface
	Create(ctx context.Context, registers ...*Register) error
	CountTotal(ctx context.Context, before *time.Time) (int64, error)
	Delete(ctx context.Context, userIDs []string) error
}
