package chat

import (
	"context"
	"github.com/OpenIMSDK/chat/pkg/common/db/table/chat"
	"github.com/OpenIMSDK/tools/errs"
	"gorm.io/gorm"
)

type Emoticons struct {
	db *gorm.DB
}

func (e *Emoticons) AddEmoticon(ctx context.Context, emoticon *chat.Emoticon) error {
	result := e.db.WithContext(ctx).Create(emoticon)
	return errs.Wrap(result.Error) // Wrap the error using the errs package
}

func (e *Emoticons) DeleteEmoticon(ctx context.Context, userID string, emoticonID int64) error {
	result := e.db.WithContext(ctx).Where("id = ? AND owner_id = ?", emoticonID, userID).Delete(&chat.Emoticon{})
	return errs.Wrap(result.Error)
}

func (e *Emoticons) GetEmoticon(ctx context.Context, userID string, emoticonID int64) (*chat.Emoticon, error) {
	var emoticon chat.Emoticon
	result := e.db.WithContext(ctx).Where("id = ? AND owner_id = ?", emoticonID, userID).First(&emoticon)
	if result.Error != nil {
		return nil, errs.Wrap(result.Error)
	}
	return &emoticon, nil
}

func (e *Emoticons) AddImage(ctx context.Context, image *chat.Image) error {
	result := e.db.WithContext(ctx).Create(image)
	return errs.Wrap(result.Error)
}
func (e *Emoticons) DeleteImage(ctx context.Context, userID string, imageID int64) error {
	result := e.db.WithContext(ctx).Where("id = ? AND owner_id = ?", imageID, userID).Delete(&chat.Image{})
	return errs.Wrap(result.Error)
}
func (e *Emoticons) GetImages(ctx context.Context, userID string) ([]*chat.Image, error) {
	var images []*chat.Image
	result := e.db.WithContext(ctx).Where("owner_id = ?", userID).Find(&images)
	if result.Error != nil {
		return nil, errs.Wrap(result.Error)
	}
	return images, nil
}

func NewEmoticons(db *gorm.DB) chat.EmoticonInterface {
	return &Emoticons{db: db}
}
