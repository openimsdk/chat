package chat

import (
	"context"
	"time"
)

type Log struct {
	LogID      string    `gorm:"column:log_id;primary_key;type:char(64)"`
	Platform   string    `gorm:"column:platform;type:varchar(32)"`
	UserID     string    `gorm:"column:user_id;type:char(64)"`
	CreateTime time.Time `gorm:"column:create_time"`
	Url        string    `gorm:"column:url;type varchar(255)"`
}

type LogInterface interface {
	Create(ctx context.Context, log *Log) error
	Search(ctx context.Context, keyword string, start time.Time, end time.Time, pageNumber int32, showNumber int32) (uint32, []*Log, error)
	Delete(ctx context.Context, logID []string, userID string) error
	Get(ctx context.Context, logIDs []string, userID string) ([]*Log, error)
}
