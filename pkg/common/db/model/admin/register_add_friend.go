package admin

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/ormutil"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/chat/pkg/common/db/table/admin"
	"gorm.io/gorm"
)

func NewRegisterAddFriend(db *gorm.DB) admin.RegisterAddFriendInterface {
	return &RegisterAddFriend{db: db}
}

type RegisterAddFriend struct {
	db *gorm.DB
}

func (o *RegisterAddFriend) Add(ctx context.Context, registerAddFriends []*admin.RegisterAddFriend) error {
	return errs.Wrap(o.db.WithContext(ctx).Create(registerAddFriends).Error)
}

func (o *RegisterAddFriend) Del(ctx context.Context, userIDs []string) error {
	return errs.Wrap(o.db.WithContext(ctx).Where("user_id in ?", userIDs).Delete(&admin.RegisterAddFriend{}).Error)
}

func (o *RegisterAddFriend) FindUserID(ctx context.Context, userIDs []string) ([]string, error) {
	db := o.db.WithContext(ctx).Model(&admin.RegisterAddFriend{})
	if len(userIDs) > 0 {
		db = db.Where("user_id in (?)", userIDs)
	}
	var ms []string
	if err := db.Pluck("user_id", &ms).Error; err != nil {
		return nil, errs.Wrap(err)
	}
	return ms, nil
}

func (o *RegisterAddFriend) Search(ctx context.Context, keyword string, page int32, size int32) (uint32, []*admin.RegisterAddFriend, error) {
	return ormutil.GormSearch[admin.RegisterAddFriend](o.db.WithContext(ctx), []string{"user_id"}, keyword, page, size)
}
