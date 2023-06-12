package gormimpl

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/ormutil"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/chat/pkg/common/db/table"
	"gorm.io/gorm"
)

func NewForbiddenAccount(db *gorm.DB) table.ForbiddenAccountInterface {
	return &ForbiddenAccount{db: db}
}

type ForbiddenAccount struct {
	db *gorm.DB
}

func (o *ForbiddenAccount) Create(ctx context.Context, ms []*table.ForbiddenAccount) error {
	return errs.Wrap(o.db.WithContext(ctx).Create(&ms).Error)
}

func (o *ForbiddenAccount) Take(ctx context.Context, userID string) (*table.ForbiddenAccount, error) {
	var f table.ForbiddenAccount
	return &f, errs.Wrap(o.db.WithContext(ctx).Where("user_id = ?", userID).Take(&f).Error)
}

func (o *ForbiddenAccount) Delete(ctx context.Context, userIDs []string) error {
	return errs.Wrap(o.db.WithContext(ctx).Where("user_id in ?", userIDs).Delete(&table.ForbiddenAccount{}).Error)
}

func (o *ForbiddenAccount) Find(ctx context.Context, userIDs []string) ([]*table.ForbiddenAccount, error) {
	var ms []*table.ForbiddenAccount
	return ms, errs.Wrap(o.db.WithContext(ctx).Where("user_id in ?", userIDs).Find(&ms).Error)
}

func (o *ForbiddenAccount) Search(ctx context.Context, keyword string, page int32, size int32) (uint32, []*table.ForbiddenAccount, error) {
	return ormutil.GormSearch[table.ForbiddenAccount](o.db.WithContext(ctx), []string{"user_id", "reason", "operator_user_id"}, keyword, page, size)
}
