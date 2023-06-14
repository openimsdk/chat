package admin

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/ormutil"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/chat/pkg/common/constant"
	"github.com/OpenIMSDK/chat/pkg/common/db/table/admin"
	"gorm.io/gorm"
)

func NewInvitationRegister(db *gorm.DB) admin.InvitationRegisterInterface {
	return &InvitationRegister{db: db}
}

type InvitationRegister struct {
	db *gorm.DB
}

func (o *InvitationRegister) NewTx(tx any) admin.InvitationRegisterInterface {
	return &InvitationRegister{db: tx.(*gorm.DB)}
}

func (o *InvitationRegister) Find(ctx context.Context, codes []string) ([]*admin.InvitationRegister, error) {
	var ms []*admin.InvitationRegister
	return ms, errs.Wrap(o.db.WithContext(ctx).Where("invitation_code in ?", codes).Find(&ms).Error)
}

func (o *InvitationRegister) Del(ctx context.Context, codes []string) error {
	return errs.Wrap(o.db.WithContext(ctx).Where("invitation_code in ?", codes).Delete(&admin.InvitationRegister{}).Error)
}

func (o *InvitationRegister) Create(ctx context.Context, v ...*admin.InvitationRegister) error {
	return errs.Wrap(o.db.WithContext(ctx).Create(v).Error)
}

func (o *InvitationRegister) Take(ctx context.Context, code string) (*admin.InvitationRegister, error) {
	var c admin.InvitationRegister
	return &c, errs.Wrap(o.db.WithContext(ctx).Where("code = ?", code).Take(&c).Error)
}

func (o *InvitationRegister) Update(ctx context.Context, code string, data map[string]any) error {
	return errs.Wrap(o.db.WithContext(ctx).Model(&admin.InvitationRegister{}).Where("code = ?", code).Updates(data).Error)
}

func (o *InvitationRegister) Search(ctx context.Context, keyword string, state int32, userIDs []string, codes []string, page int32, size int32) (uint32, []*admin.InvitationRegister, error) {
	db := o.db.WithContext(ctx)
	switch state {
	case constant.InvitationCodeUsed:
		db = db.Where("user_id <> ?", "")
	case constant.InvitationCodeUnused:
		db = db.Where("user_id = ?", "")
	}
	ormutil.GormIn(&db, "user_id", userIDs)
	ormutil.GormIn(&db, "invitation_code", codes)
	return ormutil.GormSearch[admin.InvitationRegister](db, []string{"invitation_code", "user_id"}, keyword, page, size)
}
