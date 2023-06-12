package gormimpl

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/ormutil"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/chat/pkg/common/db/table"
	"gorm.io/gorm"
)

func NewRegisterAddGroup(db *gorm.DB) table.RegisterAddGroupInterface {
	return &RegisterAddGroup{db: db}
}

type RegisterAddGroup struct {
	db *gorm.DB
}

func (o *RegisterAddGroup) Add(ctx context.Context, registerAddGroups []*table.RegisterAddGroup) error {
	return errs.Wrap(o.db.WithContext(ctx).Create(registerAddGroups).Error)
}

func (o *RegisterAddGroup) Del(ctx context.Context, userIDs []string) error {
	return errs.Wrap(o.db.WithContext(ctx).Where("group_id in ?", userIDs).Delete(&table.RegisterAddGroup{}).Error)
}

func (o *RegisterAddGroup) FindGroupID(ctx context.Context, userIDs []string) ([]string, error) {
	db := o.db.WithContext(ctx).Model(&table.RegisterAddGroup{})
	if len(userIDs) > 0 {
		db = db.Where("group_id in ?", userIDs)
	}
	var ms []string
	if err := db.Pluck("group_id", &ms).Error; err != nil {
		return nil, errs.Wrap(err)
	}
	return ms, nil
}

func (o *RegisterAddGroup) Search(ctx context.Context, keyword string, page int32, size int32) (uint32, []*table.RegisterAddGroup, error) {
	return ormutil.GormSearch[table.RegisterAddGroup](o.db.WithContext(ctx), []string{"group_id"}, keyword, page, size)
}
