package admin

import (
	"context"
	"time"
)

// RegisterAddGroup 注册时默认群组
type RegisterAddGroup struct {
	GroupID    string    `gorm:"column:group_id;primary_key;type:char(64)"`
	CreateTime time.Time `gorm:"column:create_time"`
}

func (RegisterAddGroup) TableName() string {
	return "register_add_groups"
}

type RegisterAddGroupInterface interface {
	Add(ctx context.Context, registerAddGroups []*RegisterAddGroup) error
	Del(ctx context.Context, userIDs []string) error
	FindGroupID(ctx context.Context, userIDs []string) ([]string, error)
	Search(ctx context.Context, keyword string, page int32, size int32) (uint32, []*RegisterAddGroup, error)
}
