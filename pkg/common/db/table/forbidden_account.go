package table

import (
	"context"
	"time"
)

// ForbiddenAccount 封号表
type ForbiddenAccount struct {
	UserID         string    `gorm:"column:user_id;index:userID;primary_key;type:char(64)"`
	Reason         string    `gorm:"column:reason;type:varchar(255)" `
	OperatorUserID string    `gorm:"column:operator_user_id;type:varchar(255)"`
	CreateTime     time.Time `gorm:"column:create_time" `
}

func (ForbiddenAccount) TableName() string {
	return "forbidden_accounts"
}

type ForbiddenAccountInterface interface {
	Create(ctx context.Context, ms []*ForbiddenAccount) error
	Take(ctx context.Context, userID string) (*ForbiddenAccount, error)
	Delete(ctx context.Context, userIDs []string) error
	Find(ctx context.Context, userIDs []string) ([]*ForbiddenAccount, error)
	Search(ctx context.Context, keyword string, page int32, size int32) (uint32, []*ForbiddenAccount, error)
}
