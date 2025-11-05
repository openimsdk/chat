package chat

import (
	"context"
	"time"
)

type VerifyCode struct {
	ID         string    `bson:"_id"`
	Account    string    `bson:"account"`
	Platform   string    `bson:"platform"`
	Code       string    `bson:"code"`
	Duration   uint      `bson:"duration"`
	Count      int       `bson:"count"`
	Used       bool      `bson:"used"`
	CreateTime time.Time `bson:"create_time"`
}

func (VerifyCode) TableName() string {
	return "verify_codes"
}

type VerifyCodeInterface interface {
	Add(ctx context.Context, ms []*VerifyCode) error
	RangeNum(ctx context.Context, account string, start time.Time, end time.Time) (int64, error)
	TakeLast(ctx context.Context, account string) (*VerifyCode, error)
	Incr(ctx context.Context, id string) error
	Delete(ctx context.Context, id string) error
}
