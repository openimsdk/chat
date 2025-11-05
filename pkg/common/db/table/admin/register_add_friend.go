package admin

import (
	"context"
	"github.com/openimsdk/tools/db/pagination"
	"time"
)

type RegisterAddFriend struct {
	UserID     string    `bson:"user_id"`
	CreateTime time.Time `bson:"create_time"`
}

func (RegisterAddFriend) TableName() string {
	return "register_add_friends"
}

type RegisterAddFriendInterface interface {
	Add(ctx context.Context, registerAddFriends []*RegisterAddFriend) error
	Del(ctx context.Context, userIDs []string) error
	FindUserID(ctx context.Context, userIDs []string) ([]string, error)
	Search(ctx context.Context, keyword string, pagination pagination.Pagination) (int64, []*RegisterAddFriend, error)
}
