package admin

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/ormutil"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/chat/pkg/common/db/table/admin"
	"gorm.io/gorm"
)

func NewRegisterAddGroup(db *gorm.DB) admin.RegisterAddGroupInterface {
	return &RegisterAddGroup{db: db}
}

type RegisterAddGroup struct {
	db *gorm.DB
}

func (o *RegisterAddGroup) Add(ctx context.Context, registerAddGroups []*admin.RegisterAddGroup) error {
	return errs.Wrap(o.db.WithContext(ctx).Create(registerAddGroups).Error)
}

func (o *RegisterAddGroup) Del(ctx context.Context, userIDs []string) error {
	return errs.Wrap(o.db.WithContext(ctx).Where("group_id in ?", userIDs).Delete(&admin.RegisterAddGroup{}).Error)
}

func (o *RegisterAddGroup) FindGroupID(ctx context.Context, userIDs []string) ([]string, error) {
	db := o.db.WithContext(ctx).Model(&admin.RegisterAddGroup{})
	if len(userIDs) > 0 {
		db = db.Where("group_id in ?", userIDs)
	}
	var ms []string
	if err := db.Pluck("group_id", &ms).Error; err != nil {
		return nil, errs.Wrap(err)
	}
	return ms, nil
}

func (o *RegisterAddGroup) Search(ctx context.Context, keyword string, page int32, size int32) (uint32, []*admin.RegisterAddGroup, error) {
	return ormutil.GormSearch[admin.RegisterAddGroup](o.db.WithContext(ctx), []string{"group_id"}, keyword, page, size)
}
