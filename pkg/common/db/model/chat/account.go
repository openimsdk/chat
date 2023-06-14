package chat

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/chat/pkg/common/db/table/chat"
	"gorm.io/gorm"
	"time"
)

func NewAccount(db *gorm.DB) chat.AccountInterface {
	return &Account{db: db}
}

type Account struct {
	db *gorm.DB
}

func (o *Account) NewTx(tx any) chat.AccountInterface {
	return &Account{db: tx.(*gorm.DB)}
}

func (o *Account) Create(ctx context.Context, accounts ...*chat.Account) error {
	return errs.Wrap(o.db.WithContext(ctx).Create(&accounts).Error)
}

func (o *Account) Take(ctx context.Context, userId string) (*chat.Account, error) {
	var a chat.Account
	return &a, errs.Wrap(o.db.WithContext(ctx).Where("user_id = ?", userId).First(&a).Error)
}

func (o *Account) Update(ctx context.Context, userID string, data map[string]any) error {
	return errs.Wrap(o.db.WithContext(ctx).Model(&chat.Account{}).Where("user_id = ?", userID).Updates(data).Error)
}

func (o *Account) UpdatePassword(ctx context.Context, userId string, password string) error {
	return errs.Wrap(o.db.WithContext(ctx).Model(&chat.Account{}).Where("user_id = ?", userId).Updates(map[string]interface{}{"password": password, "change_time": time.Now()}).Error)
}
