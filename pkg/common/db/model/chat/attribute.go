package chat

import (
	"context"
	"errors"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/ormutil"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
	"github.com/OpenIMSDK/chat/pkg/common/db/table/chat"
	"gorm.io/gorm"
)

func NewAttribute(db *gorm.DB) chat.AttributeInterface {
	return &Attribute{db: db}
}

type Attribute struct {
	db *gorm.DB
}

func (o *Attribute) NewTx(tx any) chat.AttributeInterface {
	return &Attribute{db: tx.(*gorm.DB)}
}

func (o *Attribute) Create(ctx context.Context, attribute ...*chat.Attribute) error {
	return errs.Wrap(o.db.WithContext(ctx).Create(attribute).Error)
}

func (o *Attribute) Update(ctx context.Context, userID string, data map[string]any) error {
	return errs.Wrap(o.db.WithContext(ctx).Model(&chat.Attribute{}).Where("user_id = ?", userID).Updates(data).Error)
}

func (o *Attribute) Find(ctx context.Context, userIds []string) ([]*chat.Attribute, error) {
	var a []*chat.Attribute
	return a, errs.Wrap(o.db.WithContext(ctx).Where("user_id in (?)", userIds).Find(&a).Error)
}

func (o *Attribute) FindAccount(ctx context.Context, accounts []string) ([]*chat.Attribute, error) {
	var a []*chat.Attribute
	return a, errs.Wrap(o.db.WithContext(ctx).Where("account in (?)", accounts).Find(&a).Error)
}

func (o *Attribute) Search(ctx context.Context, keyword string, genders []int32, page int32, size int32) (uint32, []*chat.Attribute, error) {
	db := o.db.WithContext(ctx)
	if len(genders) > 0 {
		db = db.Where("gender in ?", genders)
	}
	return ormutil.GormSearch[chat.Attribute](db, []string{"user_id", "account", "nickname", "phone_number"}, keyword, page, size)
}

func (o *Attribute) TakePhone(ctx context.Context, areaCode string, phoneNumber string) (*chat.Attribute, error) {
	var a chat.Attribute
	return &a, errs.Wrap(o.db.WithContext(ctx).Where("area_code = ? and phone_number = ?", areaCode, phoneNumber).First(&a).Error)
}

func (o *Attribute) TakeAccount(ctx context.Context, account string) (*chat.Attribute, error) {
	var a chat.Attribute
	return &a, errs.Wrap(o.db.WithContext(ctx).Where("account = ?", account).Take(&a).Error)
}

func (o *Attribute) Take(ctx context.Context, userID string) (*chat.Attribute, error) {
	var a chat.Attribute
	return &a, errs.Wrap(o.db.WithContext(ctx).Where("user_id = ?", userID).Take(&a).Error)
}

func (tb *Attribute) GetAccountList(ctx context.Context, accountList []string) ([]*chat.Attribute, error) {
	if len(accountList) == 0 {
		return []*chat.Attribute{}, nil
	}
	var att []*chat.Attribute
	err := tb.db.WithContext(ctx).Model(&att).Where("account in (?)", accountList).Find(&att).Error
	return att, utils.Wrap(err, "")
}

func (tb *Attribute) ExistPhoneNumber(ctx context.Context, areaCode, phoneNumber string) (bool, error) {
	var m chat.Attribute
	err := tb.db.WithContext(ctx).Model(&chat.Attribute{}).Where("area_code = ? and phone_number = ?", areaCode, phoneNumber).First(&m).Error
	if err == nil {
		return true, nil
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	} else {
		return false, utils.Wrap(err, "")
	}
}
