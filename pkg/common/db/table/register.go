package table

import (
	"context"
	"time"
)

// Register 注册信息表
type Register struct {
	UserID      string    `gorm:"column:user_id;primary_key;type:char(64)"`
	DeviceID    string    `gorm:"column:device_id;type:varchar(255)"`
	IP          string    `gorm:"column:ip;type:varchar(64)"`
	Platform    string    `gorm:"column:platform;type:varchar(32)"`
	AccountType string    `gorm:"column:account_type;type:varchar(32)"` //email phone account
	Mode        string    `gorm:"column:mode;type:varchar(32)"`         //user admin
	CreateTime  time.Time `gorm:"column:create_time"`
}

func (Register) TableName() string {
	return "registers"
}

type RegisterInterface interface {
	NewTx(tx any) RegisterInterface
	Create(ctx context.Context, registers ...*Register) error
}
