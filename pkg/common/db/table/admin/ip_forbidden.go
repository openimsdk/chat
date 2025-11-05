package admin

import (
	"context"
	"github.com/openimsdk/tools/db/pagination"
	"time"
)

type IPForbidden struct {
	IP            string    `bson:"ip"`
	LimitRegister bool      `bson:"limit_register"`
	LimitLogin    bool      `bson:"limit_login"`
	CreateTime    time.Time `bson:"create_time"`
}

func (IPForbidden) IPForbidden() string {
	return "ip_forbiddens"
}

type IPForbiddenInterface interface {
	Take(ctx context.Context, ip string) (*IPForbidden, error)
	Find(ctx context.Context, ips []string) ([]*IPForbidden, error)
	Search(ctx context.Context, keyword string, state int32, pagination pagination.Pagination) (int64, []*IPForbidden, error)
	Create(ctx context.Context, ms []*IPForbidden) error
	Delete(ctx context.Context, ips []string) error
}
