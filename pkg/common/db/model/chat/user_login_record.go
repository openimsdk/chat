package chat

import (
	"context"
	"github.com/OpenIMSDK/chat/pkg/common/db/table/chat"
	"gorm.io/gorm"
)

func NewUserLoginRecord(db *gorm.DB) chat.UserLoginRecordInterface {
	return &UserLoginRecord{
		db: db,
	}
}

type UserLoginRecord struct {
	db *gorm.DB
}

func (o *UserLoginRecord) NewTx(tx any) chat.UserLoginRecordInterface {
	return &UserLoginRecord{db: tx.(*gorm.DB)}
}

func (o *UserLoginRecord) Create(ctx context.Context, records ...*chat.UserLoginRecord) error {
	return o.db.WithContext(ctx).Create(&records).Error
}
