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

func NewLimitUserLoginIP(db *gorm.DB) admin.LimitUserLoginIPInterface {
	return &LimitUserLoginIP{db: db}
}

type LimitUserLoginIP struct {
	db *gorm.DB
}

func (o *LimitUserLoginIP) Create(ctx context.Context, ms []*admin.LimitUserLoginIP) error {
	return errs.Wrap(o.db.WithContext(ctx).Create(&ms).Error)
}

func (o *LimitUserLoginIP) Delete(ctx context.Context, ms []*admin.LimitUserLoginIP) error {
	return errs.Wrap(o.db.WithContext(ctx).Delete(&ms).Error)
}

func (o *LimitUserLoginIP) Count(ctx context.Context, userID string) (uint32, error) {
	var count int64
	if err := o.db.WithContext(ctx).Model(&admin.LimitUserLoginIP{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return 0, errs.Wrap(err)
	}
	return uint32(count), nil
}

func (o *LimitUserLoginIP) Take(ctx context.Context, userID string, ip string) (*admin.LimitUserLoginIP, error) {
	var f admin.LimitUserLoginIP
	return &f, errs.Wrap(o.db.WithContext(ctx).Where("user_id = ? and ip = ?", userID, ip).Take(&f).Error)
}

func (o *LimitUserLoginIP) Search(ctx context.Context, keyword string, page int32, size int32) (uint32, []*admin.LimitUserLoginIP, error) {
	return ormutil.GormSearch[admin.LimitUserLoginIP](o.db.WithContext(ctx), []string{"user_id", "ip"}, keyword, page, size)
}
