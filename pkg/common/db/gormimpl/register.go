package gormimpl

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/chat/pkg/common/db/table"
	"gorm.io/gorm"
)

func NewRegister(db *gorm.DB) table.RegisterInterface {
	return &Register{db: db}
}

type Register struct {
	db *gorm.DB
}

func (o *Register) NewTx(tx any) table.RegisterInterface {
	return &Register{db: tx.(*gorm.DB)}
}

func (o *Register) Create(ctx context.Context, registers ...*table.Register) error {
	return errs.Wrap(o.db.WithContext(ctx).Create(registers).Error)
}
