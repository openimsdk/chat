package admin

import (
	"context"
	"time"
)

// 邀请码被注册使用
type InvitationRegister struct {
	InvitationCode string    `gorm:"column:invitation_code;primary_key;type:char(32)"`
	UsedByUserID   string    `gorm:"column:user_id;index:userID;type:char(64)"`
	CreateTime     time.Time `gorm:"column:create_time"`
}

func (InvitationRegister) TableName() string {
	return "invitation_registers"
}

type InvitationRegisterInterface interface {
	NewTx(tx any) InvitationRegisterInterface
	Find(ctx context.Context, codes []string) ([]*InvitationRegister, error)
	Del(ctx context.Context, codes []string) error
	Create(ctx context.Context, v ...*InvitationRegister) error
	Take(ctx context.Context, code string) (*InvitationRegister, error)
	Update(ctx context.Context, code string, data map[string]any) error
	Search(ctx context.Context, keyword string, state int32, userIDs []string, codes []string, page int32, size int32) (uint32, []*InvitationRegister, error)
}
