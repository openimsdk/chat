package gormimpl

import (
	"context"
	"github.com/OpenIMSDK/chat/pkg/common/db/table"
	"gorm.io/gorm"
)

func NewUserLoginRecord(db *gorm.DB) table.UserLoginRecordInterface {
	return &UserLoginRecord{
		db: db,
	}
}

type UserLoginRecord struct {
	db *gorm.DB
}

func (o *UserLoginRecord) NewTx(tx any) table.UserLoginRecordInterface {
	return &UserLoginRecord{db: tx.(*gorm.DB)}
}

func (o *UserLoginRecord) Create(ctx context.Context, records ...*table.UserLoginRecord) error {
	return o.db.WithContext(ctx).Create(&records).Error
}
