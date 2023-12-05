package chat

import (
	"context"
)

type Emoticon struct {
	ID       int64  `gorm:"column:id;primary_key;"`
	ImageURL string `gorm:"column:image_url;type:varchar(255)"`
	OwnerID  string `gorm:"column:owner_id;type:char(64)"`
}
type EmoticonInterface interface {
	AddEmoticon(ctx context.Context, emoticon *Emoticon) error
	DeleteEmoticon(ctx context.Context, userId string, id string) error
}
