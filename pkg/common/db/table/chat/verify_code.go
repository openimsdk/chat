package chat

import (
	"context"
	"time"
)

type VerifyCode struct {
	ID         uint      `gorm:"column:id;primary_key;autoIncrement"`
	Account    string    `gorm:"column:account;type:char(64)"`
	Platform   string    `gorm:"column:platform;type:varchar(32)"`
	Code       string    `gorm:"column:verify_code;type:varchar(16)"`
	Duration   uint      `gorm:"column:duration;type:int(11)"`
	Count      int       `gorm:"column:count;type:int(11)"`
	Used       bool      `gorm:"column:used"`
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime"`
}

func (VerifyCode) TableName() string {
	return "verify_codes"
}

type VerifyCodeInterface interface {
	NewTx(tx any) VerifyCodeInterface
	Add(ctx context.Context, ms []*VerifyCode) error
	RangeNum(ctx context.Context, account string, start time.Time, end time.Time) (uint32, error)
	TakeLast(ctx context.Context, account string) (*VerifyCode, error)
	Incr(ctx context.Context, id uint) error
	Delete(ctx context.Context, id uint) error
}
