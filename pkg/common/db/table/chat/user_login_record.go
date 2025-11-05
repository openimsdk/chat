package chat

import (
	"context"
	"time"
)

type UserLoginRecord struct {
	UserID    string    `bson:"user_id"`
	LoginTime time.Time `bson:"login_time"`
	IP        string    `bson:"ip"`
	DeviceID  string    `bson:"device_id"`
	Platform  string    `bson:"platform"`
}

func (UserLoginRecord) TableName() string {
	return "user_login_records"
}

type UserLoginRecordInterface interface {
	Create(ctx context.Context, records ...*UserLoginRecord) error
	CountTotal(ctx context.Context, before *time.Time) (int64, error)
	CountRangeEverydayTotal(ctx context.Context, start *time.Time, end *time.Time) (map[string]int64, int64, error)
}
