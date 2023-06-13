package chat

import (
	"context"
	"time"
)

// 用户登录信息表
type UserLoginRecord struct {
	UserID    string    `gorm:"column:user_id;size:64"`
	LoginTime time.Time `gorm:"column:login_time"`
	IP        string    `gorm:"column:ip;type:varchar(32)"`
	DeviceID  string    `gorm:"column:device_id;type:varchar(255)"`
	Platform  string    `gorm:"column:platform;type:varchar(32)"`
}

func (UserLoginRecord) TableName() string {
	return "user_login_records"
}

type UserLoginRecordInterface interface {
	NewTx(tx any) UserLoginRecordInterface
	Create(ctx context.Context, records ...*UserLoginRecord) error
}
