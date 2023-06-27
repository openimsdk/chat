package organization

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Organization struct {
	LogoURL         string    `gorm:"column:logo_url;size:255" json:"logoURL"`
	Name            string    `gorm:"column:name;size:256" json:"name" binding:"required"`
	Homepage        string    `gorm:"column:homepage" json:"homepage" `
	RelatedGroupID  string    `gorm:"column:related_group_id;size:64" json:"relatedGroupID"`
	Introduction    string    `gorm:"column:introduction;size:255" json:"introduction"`
	DefaultPassword string    `gorm:"column:default_password" json:"defaultPassword"`
	CreateTime      time.Time `gorm:"column:create_time" json:"createTime"`
	ChangeTime      time.Time `gorm:"column:change_time" json:"changeTime"`
}

type OrganizationInterface interface {
	Set(ctx context.Context, m *Organization) error
	Get(ctx context.Context) (*Organization, error)
	BeginTransaction(ctx context.Context) (*gorm.DB, error)
}
