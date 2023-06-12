package gormimpl

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/ormutil"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/chat/pkg/common/constant"
	"github.com/OpenIMSDK/chat/pkg/common/db/table"
	"gorm.io/gorm"
)

func NewIPForbidden(db *gorm.DB) table.IPForbiddenInterface {
	return &IPForbidden{db: db}
}

type IPForbidden struct {
	db *gorm.DB
}

func (o *IPForbidden) NewTx(tx any) table.IPForbiddenInterface {
	return &IPForbidden{db: tx.(*gorm.DB)}
}

func (o *IPForbidden) Take(ctx context.Context, ip string) (*table.IPForbidden, error) {
	var f table.IPForbidden
	return &f, errs.Wrap(o.db.WithContext(ctx).Where("ip = ?", ip).Take(&f).Error)
}

func (o *IPForbidden) Find(ctx context.Context, ips []string) ([]*table.IPForbidden, error) {
	var forbiddens []*table.IPForbidden
	return forbiddens, errs.Wrap(o.db.WithContext(ctx).Where("ip in ?", ips).Find(&forbiddens).Error)
}

func (o *IPForbidden) Search(ctx context.Context, keyword string, state int32, page int32, size int32) (uint32, []*table.IPForbidden, error) {
	db := o.db.WithContext(ctx)
	switch state {
	case constant.LimitNil:
	case constant.LimitEmpty:
		db = db.Where("limit_register = 0 and limit_login = 0")
	case constant.LimitOnlyRegisterIP:
		db = db.Where("limit_register = 1 and limit_login = 0")
	case constant.LimitOnlyLoginIP:
		db = db.Where("limit_register = 0 and limit_login = 1")
	case constant.LimitRegisterIP:
		db = db.Where("limit_register = 1")
	case constant.LimitLoginIP:
		db = db.Where("limit_login = 1")
	case constant.LimitLoginRegisterIP:
		db = db.Where("limit_register = 1 and limit_login = 1")
	}
	return ormutil.GormSearch[table.IPForbidden](db, []string{"ip"}, keyword, page, size)
}

func (o *IPForbidden) Create(ctx context.Context, ms []*table.IPForbidden) error {
	return errs.Wrap(o.db.WithContext(ctx).Create(&ms).Error)
}

func (o *IPForbidden) Delete(ctx context.Context, ips []string) error {
	return errs.Wrap(o.db.WithContext(ctx).Where("ip in ?", ips).Delete(&table.IPForbidden{}).Error)
}
