package admin

import (
	"context"
	"time"
)

// 限制userID只能在某些ip登录
type LimitUserLoginIP struct {
	UserID     string    `gorm:"column:user_id;primary_key;type:char(64)"`
	IP         string    `gorm:"column:ip;primary_key;type:char(32)"`
	CreateTime time.Time `gorm:"column:create_time" `
}

func (LimitUserLoginIP) TableName() string {
	return "limit_user_login_ips"
}

type LimitUserLoginIPInterface interface {
	Create(ctx context.Context, ms []*LimitUserLoginIP) error
	Delete(ctx context.Context, ms []*LimitUserLoginIP) error
	Count(ctx context.Context, userID string) (uint32, error)
	Take(ctx context.Context, userID string, ip string) (*LimitUserLoginIP, error)
	Search(ctx context.Context, keyword string, page int32, size int32) (uint32, []*LimitUserLoginIP, error)
}
