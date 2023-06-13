package admin

import (
	"context"
	"time"
)

type Applet struct {
	ID         string    `gorm:"column:id;primary_key;size:64"`
	Name       string    `gorm:"column:name;size:64"`
	AppID      string    `gorm:"column:app_id;uniqueIndex;size:255"`
	Icon       string    `gorm:"column:icon;size:255"`
	URL        string    `gorm:"column:url;size:255"`
	MD5        string    `gorm:"column:md5;size:255"`
	Size       int64     `gorm:"column:size"`
	Version    string    `gorm:"column:version;size:64"`
	Priority   uint32    `gorm:"column:priority;size:64"`
	Status     uint8     `gorm:"column:status"`
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime;size:64"`
}

func (Applet) TableName() string {
	return "applets"
}

type AppletInterface interface {
	Create(ctx context.Context, applets ...*Applet) error
	Del(ctx context.Context, ids []string) error
	Update(ctx context.Context, id string, data map[string]any) error
	Take(ctx context.Context, id string) (*Applet, error)
	Search(ctx context.Context, keyword string, page int32, size int32) (uint32, []*Applet, error)
	FindOnShelf(ctx context.Context) ([]*Applet, error)
	FindID(ctx context.Context, ids []string) ([]*Applet, error)
}
