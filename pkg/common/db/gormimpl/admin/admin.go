package admin

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/chat/pkg/common/db/table/admin"
	"gorm.io/gorm"
)

func NewAdmin(db *gorm.DB) *Admin {
	return &Admin{
		db: db,
	}
}

type Admin struct {
	db *gorm.DB
}

func (o *Admin) Take(ctx context.Context, account string) (*admin.Admin, error) {
	var a admin.Admin
	return &a, errs.Wrap(o.db.WithContext(ctx).Where("account = ?", account).Take(&a).Error)
}

func (o *Admin) TakeUserID(ctx context.Context, userID string) (*admin.Admin, error) {
	var a admin.Admin
	return &a, errs.Wrap(o.db.WithContext(ctx).Where("user_id = ?", userID).Take(&a).Error)
}

func (o *Admin) Update(ctx context.Context, account string, update map[string]any) error {
	return errs.Wrap(o.db.WithContext(ctx).Model(&admin.Admin{}).Where("user_id = ?", account).Updates(update).Error)
}
