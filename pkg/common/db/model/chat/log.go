package chat

import (
	"context"
	"github.com/OpenIMSDK/chat/pkg/common/db/table/chat"
	"github.com/OpenIMSDK/tools/errs"
	"github.com/OpenIMSDK/tools/ormutil"
	"gorm.io/gorm"
	"time"
)

type Logs struct {
	db *gorm.DB
}

func (l *Logs) Create(ctx context.Context, log []*chat.Log) error {
	return errs.Wrap(l.db.WithContext(ctx).Create(log).Error)
}

func (l *Logs) Search(ctx context.Context, keyword string, start time.Time, end time.Time, pageNumber int32, showNumber int32) (uint32, []*chat.Log, error) {
	db := l.db.WithContext(ctx).Where("create_time >= ?", start)
	if end.UnixMilli() != 0 {
		db = l.db.WithContext(ctx).Where("create_time <= ?", end)
	}
	return ormutil.GormSearch[chat.Log](db, []string{"user_id"}, keyword, pageNumber, showNumber)
}

func (l *Logs) Delete(ctx context.Context, logIDs []string, userID string) error {
	if userID == "" {
		return errs.Wrap(l.db.WithContext(ctx).Where("log_id in ?", logIDs).Delete(&chat.Log{}).Error)
	}
	return errs.Wrap(l.db.WithContext(ctx).Where("log_id in ? and user_id=?", logIDs, userID).Delete(&chat.Log{}).Error)
}

func (l *Logs) Get(ctx context.Context, logIDs []string, userID string) ([]*chat.Log, error) {
	var logs []*chat.Log
	if userID == "" {
		return logs, errs.Wrap(l.db.WithContext(ctx).Where("log_id in ?", logIDs).Find(&logs).Error)
	}
	return logs, errs.Wrap(l.db.WithContext(ctx).Where("log_id in ? and user_id=?", logIDs, userID).Find(&logs).Error)
}
func NewLogs(db *gorm.DB) chat.LogInterface {
	return &Logs{db: db}
}
