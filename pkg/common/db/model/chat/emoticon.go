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
	return errs.Wrap(e.db.WithContext(ctx).Create(emoticon).Error)
}

func (e *Emoticons) DeleteEmoticon(ctx context.Context, userID, emoticonID string) error {
	return errs.Wrap(e.db.WithContext(ctx).Where("id = ? AND owner_id = ?", emoticonID, userID).Delete(&chat.Emoticon{}).Error)
}

func NewEmoticons(db *gorm.DB) chat.EmoticonInterface {
	return &Emoticons{db: db}
}
