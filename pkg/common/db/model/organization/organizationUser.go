package organization

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
	table "github.com/OpenIMSDK/chat/pkg/common/db/table/organization"
	"gorm.io/gorm"
	"time"
)

func NewOrganizationUser(db *gorm.DB) *OrganizationUser {
	return &OrganizationUser{
		db: db,
	}
}

type OrganizationUser struct {
	db *gorm.DB
}

func (tb *OrganizationUser) Create(ctx context.Context, m *table.OrganizationUser) error {
	m.CreateTime = time.Now()
	m.ChangeTime = time.Now()
	return utils.Wrap(tb.db.WithContext(ctx).Create(m).Error, "")
}

func (tb *OrganizationUser) Update(ctx context.Context, m *table.OrganizationUser) error {
	m.ChangeTime = time.Now()
	return utils.Wrap(tb.db.WithContext(ctx).Where("user_id = ?", m.UserID).Updates(&m).Error, "")
}

func (tb *OrganizationUser) Delete(ctx context.Context, userID string) error {
	return utils.Wrap(tb.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&table.OrganizationUser{}).Error, "")
}
