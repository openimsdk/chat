package admin

import (
	"context"
	"github.com/openimsdk/tools/db/pagination"
	"time"
)

type InvitationRegister struct {
	InvitationCode string    `bson:"invitation_code"`
	UsedByUserID   string    `bson:"used_by_user_id"`
	CreateTime     time.Time `bson:"create_time"`
}

func (InvitationRegister) TableName() string {
	return "invitation_registers"
}

type InvitationRegisterInterface interface {
	Find(ctx context.Context, codes []string) ([]*InvitationRegister, error)
	Del(ctx context.Context, codes []string) error
	Create(ctx context.Context, v []*InvitationRegister) error
	Take(ctx context.Context, code string) (*InvitationRegister, error)
	Update(ctx context.Context, code string, data map[string]any) error
	Search(ctx context.Context, keyword string, state int32, userIDs []string, codes []string, pagination pagination.Pagination) (int64, []*InvitationRegister, error)
}
