package table

import (
	"context"
	"time"
)

// RegisterAddFriend 注册时默认好友
type RegisterAddFriend struct {
	UserID     string    `gorm:"column:user_id;primary_key;type:char(64)"`
	CreateTime time.Time `gorm:"column:create_time"`
}

func (RegisterAddFriend) TableName() string {
	return "register_add_friends"
}

type RegisterAddFriendInterface interface {
	Add(ctx context.Context, registerAddFriends []*RegisterAddFriend) error
	Del(ctx context.Context, userIDs []string) error
	FindUserID(ctx context.Context, userIDs []string) ([]string, error)
	Search(ctx context.Context, keyword string, page int32, size int32) (uint32, []*RegisterAddFriend, error)
}
