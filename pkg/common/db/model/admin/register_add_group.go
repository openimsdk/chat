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

func NewRegisterAddGroup(db *gorm.DB) admin.RegisterAddGroupInterface {
	return &RegisterAddGroup{db: db}
}

type RegisterAddGroup struct {
	db *gorm.DB
}

func (o *RegisterAddGroup) Add(ctx context.Context, registerAddGroups []*admin.RegisterAddGroup) error {
	return errs.Wrap(o.db.WithContext(ctx).Create(registerAddGroups).Error)
}

func (o *RegisterAddGroup) Del(ctx context.Context, userIDs []string) error {
	return errs.Wrap(o.db.WithContext(ctx).Where("group_id in ?", userIDs).Delete(&admin.RegisterAddGroup{}).Error)
}

func (o *RegisterAddGroup) FindGroupID(ctx context.Context, userIDs []string) ([]string, error) {
	db := o.db.WithContext(ctx).Model(&admin.RegisterAddGroup{})
	if len(userIDs) > 0 {
		db = db.Where("group_id in ?", userIDs)
	}
	var ms []string
	if err := db.Pluck("group_id", &ms).Error; err != nil {
		return nil, errs.Wrap(err)
	}
	return ms, nil
}

func (o *RegisterAddGroup) Search(ctx context.Context, keyword string, page int32, size int32) (uint32, []*admin.RegisterAddGroup, error) {
	return ormutil.GormSearch[admin.RegisterAddGroup](o.db.WithContext(ctx), []string{"group_id"}, keyword, page, size)
}
