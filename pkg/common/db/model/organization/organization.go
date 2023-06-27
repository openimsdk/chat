package organization

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
	table "github.com/OpenIMSDK/chat/pkg/common/db/table/organization"
	"gorm.io/gorm"
	"time"
)

func NewOrganization(db *gorm.DB) *Organization {
	return &Organization{
		db: db,
	}
}

type Organization struct {
	db *gorm.DB
}

func (o *Organization) Set(ctx context.Context, m *table.Organization) error {
	var org Organization
	if err := o.db.First(&org).Error; err == nil {
		m.CreateTime = time.Time{}
		m.ChangeTime = time.Now()
		return utils.Wrap(o.db.WithContext(ctx).Where("1 = 1").Updates(m).Error, "")
	} else if err == gorm.ErrRecordNotFound {
		m.CreateTime = time.Now()
		m.ChangeTime = time.Now()
		return utils.Wrap(o.db.WithContext(ctx).Create(m).Error, "")
	} else {
		return utils.Wrap(err, "")
	}
}

func (o *Organization) Get(ctx context.Context) (*table.Organization, error) {
	var m table.Organization
	if err := o.db.WithContext(ctx).First(&m).Error; err == gorm.ErrRecordNotFound {
		m.CreateTime = time.UnixMilli(0)
	} else if err != nil {
		return nil, utils.Wrap(err, "")
	}
	return &m, nil
}
