package gormimpl

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/chat/pkg/common/db/table"
	"gorm.io/gorm"
	"time"
)

func NewVerifyCode(db *gorm.DB) *VerifyCode {
	return &VerifyCode{
		db: db,
	}
}

type VerifyCode struct {
	db *gorm.DB
}

func (o *VerifyCode) NewTx(tx any) table.VerifyCodeInterface {
	return &VerifyCode{db: tx.(*gorm.DB)}
}

func (o *VerifyCode) Add(ctx context.Context, ms []*table.VerifyCode) error {
	return errs.Wrap(o.db.WithContext(ctx).Create(&ms).Error)
}

func (o *VerifyCode) RangeNum(ctx context.Context, account string, start time.Time, end time.Time) (uint32, error) {
	var count int64
	if err := o.db.WithContext(ctx).Model(&table.VerifyCode{}).Where("account = ?", account).Where("create_time BETWEEN ? AND ?", start, end).Count(&count).Error; err != nil {
		return 0, errs.Wrap(err)
	}
	return uint32(count), nil
}

func (o *VerifyCode) TakeLast(ctx context.Context, account string) (*table.VerifyCode, error) {
	var m table.VerifyCode
	return &m, errs.Wrap(o.db.WithContext(ctx).Where("account = ?", account).Order("id DESC").Take(&m).Error)
}

func (o *VerifyCode) Incr(ctx context.Context, id uint) error {
	return errs.Wrap(o.db.WithContext(ctx).Model(&table.VerifyCode{}).Where("id = ?", id).Updates(map[string]any{"count": gorm.Expr("count + 1")}).Error)
}

func (o *VerifyCode) Delete(ctx context.Context, id uint) error {
	return errs.Wrap(o.db.WithContext(ctx).Where("id = ?", id).Delete(&table.VerifyCode{}).Error)
}
