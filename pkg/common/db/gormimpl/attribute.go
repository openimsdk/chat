package gormimpl

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/ormutil"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/chat/pkg/common/db/table"
	"gorm.io/gorm"
)

func NewAttribute(db *gorm.DB) table.AttributeInterface {
	return &Attribute{db: db}
}

type Attribute struct {
	db *gorm.DB
}

func (o *Attribute) NewTx(tx any) table.AttributeInterface {
	return &Attribute{db: tx.(*gorm.DB)}
}

func (o *Attribute) Create(ctx context.Context, attribute ...*table.Attribute) error {
	return errs.Wrap(o.db.WithContext(ctx).Create(attribute).Error)
}

func (o *Attribute) Update(ctx context.Context, userID string, data map[string]any) error {
	return errs.Wrap(o.db.WithContext(ctx).Model(&table.Attribute{}).Where("user_id = ?", userID).Updates(data).Error)
}

func (o *Attribute) Find(ctx context.Context, userIds []string) ([]*table.Attribute, error) {
	var a []*table.Attribute
	return a, errs.Wrap(o.db.WithContext(ctx).Where("user_id in (?)", userIds).Find(&a).Error)
}

func (o *Attribute) FindAccount(ctx context.Context, accounts []string) ([]*table.Attribute, error) {
	var a []*table.Attribute
	return a, errs.Wrap(o.db.WithContext(ctx).Where("account in (?)", accounts).Find(&a).Error)
}

func (o *Attribute) Search(ctx context.Context, keyword string, genders []int32, page int32, size int32) (uint32, []*table.Attribute, error) {
	db := o.db.WithContext(ctx)
	if len(genders) > 0 {
		db = db.Where("gender in ?", genders)
	}
	return ormutil.GormSearch[table.Attribute](db, []string{"user_id", "account", "nickname", "phone_number"}, keyword, page, size)
}

func (o *Attribute) TakePhone(ctx context.Context, areaCode string, phoneNumber string) (*table.Attribute, error) {
	var a table.Attribute
	return &a, errs.Wrap(o.db.WithContext(ctx).Where("area_code = ? and phone_number = ?", areaCode, phoneNumber).First(&a).Error)
}

func (o *Attribute) TakeAccount(ctx context.Context, account string) (*table.Attribute, error) {
	var a table.Attribute
	return &a, errs.Wrap(o.db.WithContext(ctx).Where("account = ?", account).Take(&a).Error)
}

func (o *Attribute) Take(ctx context.Context, userID string) (*table.Attribute, error) {
	var a table.Attribute
	return &a, errs.Wrap(o.db.WithContext(ctx).Where("user_id = ?", userID).Take(&a).Error)
}
