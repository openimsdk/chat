package admin

import (
	"context"
	"github.com/openimsdk/tools/db/pagination"
	"time"
)

// ForbiddenAccount table
type ForbiddenAccount struct {
	UserID         string    `bson:"user_id"`
	Reason         string    `bson:"reason"`
	OperatorUserID string    `bson:"operator_user_id"`
	CreateTime     time.Time `bson:"create_time"`
}

func (ForbiddenAccount) TableName() string {
	return "forbidden_accounts"
}

type ForbiddenAccountInterface interface {
	Create(ctx context.Context, ms []*ForbiddenAccount) error
	Take(ctx context.Context, userID string) (*ForbiddenAccount, error)
	Delete(ctx context.Context, userIDs []string) error
	Find(ctx context.Context, userIDs []string) ([]*ForbiddenAccount, error)
	Search(ctx context.Context, keyword string, pagination pagination.Pagination) (int64, []*ForbiddenAccount, error)
	FindAllIDs(ctx context.Context) ([]string, error)
}
