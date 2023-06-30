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
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/chat/pkg/common/db/table/admin"
	"gorm.io/gorm"
)

func NewAdmin(db *gorm.DB) *Admin {
	return &Admin{
		db: db,
	}
}

type Admin struct {
	db *gorm.DB
}

func (o *Admin) Take(ctx context.Context, account string) (*admin.Admin, error) {
	var a admin.Admin
	return &a, errs.Wrap(o.db.WithContext(ctx).Where("account = ?", account).Take(&a).Error)
}

func (o *Admin) TakeUserID(ctx context.Context, userID string) (*admin.Admin, error) {
	var a admin.Admin
	return &a, errs.Wrap(o.db.WithContext(ctx).Where("user_id = ?", userID).Take(&a).Error)
}

func (o *Admin) Update(ctx context.Context, account string, update map[string]any) error {
	return errs.Wrap(o.db.WithContext(ctx).Model(&admin.Admin{}).Where("user_id = ?", account).Updates(update).Error)
}
