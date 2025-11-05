package admin

import (
	"context"
	"github.com/openimsdk/tools/db/pagination"
	"time"
)

type LimitUserLoginIP struct {
	UserID     string    `bson:"user_id"`
	IP         string    `bson:"ip"`
	CreateTime time.Time `bson:"create_time"`
}

func (LimitUserLoginIP) TableName() string {
	return "limit_user_login_ips"
}

type LimitUserLoginIPInterface interface {
	Create(ctx context.Context, ms []*LimitUserLoginIP) error
	Delete(ctx context.Context, ms []*LimitUserLoginIP) error
	Count(ctx context.Context, userID string) (uint32, error)
	Take(ctx context.Context, userID string, ip string) (*LimitUserLoginIP, error)
	Search(ctx context.Context, keyword string, pagination pagination.Pagination) (int64, []*LimitUserLoginIP, error)
}
