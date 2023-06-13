package chat

import (
	"context"
	"time"
)

// Account 账号密码表
type Account struct {
	UserID         string    `gorm:"column:user_id;primary_key;type:char(64)"`
	Password       string    `gorm:"column:password;type:varchar(32)"`
	CreateTime     time.Time `gorm:"column:create_time;autoCreateTime"`
	ChangeTime     time.Time `gorm:"column:change_time;autoUpdateTime"`
	OperatorUserID string    `gorm:"column:operator_user_id;type:varchar(64)"`
}

func (Account) TableName() string {
	return "accounts"
}

type AccountInterface interface {
	NewTx(tx any) AccountInterface
	Create(ctx context.Context, accounts ...*Account) error
	Take(ctx context.Context, userId string) (*Account, error)
	Update(ctx context.Context, userID string, data map[string]any) error
	UpdatePassword(ctx context.Context, userId string, password string) error
}
