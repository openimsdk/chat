package chat

import (
	"context"
)

type Image struct {
	ID       int64  `gorm:"column:id;primary_key;"`
	ImageURL string `gorm:"column:image_url;type:varchar(255)"`
	//EmoticonID string `gorm:"column:owner_id;type:char(64)"`
	OwnerID string `gorm:"column:owner_id;type:char(64)"`
}

func (Image) TableName() string {
	return "images"
}

type Emoticon struct {
	EmoticonID int64  `gorm:"column:id;primary_key;"`
	OwnerID    string `gorm:"column:owner_id;type:char(64)"`
}

func (Emoticon) TableName() string {
	return "emoticons"
}

type EmoticonInterface interface {
	AddEmoticon(ctx context.Context, emoticon *Emoticon) error
	DeleteEmoticon(ctx context.Context, userID string, emoticonID int64) error
	GetEmoticon(ctx context.Context, userID string, emoticonID int64) (*Emoticon, error)

	AddImage(ctx context.Context, image *Image) error
	DeleteImage(ctx context.Context, userID string, imageID int64) error
	GetImages(ctx context.Context, userID string) ([]*Image, error)
}
