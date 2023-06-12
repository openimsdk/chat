package table

import (
	"context"
	"time"
)

// Admin 后台管理员
type Admin struct {
	Account    string    `gorm:"column:account;primary_key;type:char(64)"`
	Password   string    `gorm:"column:password;type:char(64)"`
	FaceURL    string    `gorm:"column:face_url;type:char(64)"`
	Nickname   string    `gorm:"column:nickname;type:char(64)"`
	UserID     string    `gorm:"column:user_id;type:char(64)"` //openIM userID
	Level      int32     `gorm:"column:level;default:1"  `
	CreateTime time.Time `gorm:"column:create_time"`
}

func (Admin) TableName() string {
	return "admins"
}

type AdminInterface interface {
	Take(ctx context.Context, account string) (*Admin, error)
	TakeUserID(ctx context.Context, userID string) (*Admin, error)
	Update(ctx context.Context, account string, update map[string]any) error
}
