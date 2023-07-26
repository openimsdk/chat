// Copyright © 2023 OpenIM open source community. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package admin

import (
	"context"

	"github.com/OpenIMSDK/tools/errs"
	"github.com/OpenIMSDK/tools/ormutil"
	"gorm.io/gorm"

	"github.com/OpenIMSDK/chat/pkg/common/db/table/admin"
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
