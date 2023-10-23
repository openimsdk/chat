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
	"crypto/md5"
	"encoding/hex"
	"github.com/OpenIMSDK/tools/log"
	"time"

	"github.com/OpenIMSDK/chat/pkg/common/config"
	"github.com/OpenIMSDK/chat/pkg/common/db/table/admin"
	"github.com/OpenIMSDK/tools/errs"
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

func (o *Admin) InitAdmin(ctx context.Context) error {
	var count int64
	if err := o.db.WithContext(ctx).Model(&admin.Admin{}).Count(&count).Error; err != nil {
		return errs.Wrap(err)
	}
	if count > 0 {
		log.ZInfo(ctx, "Admins are already registered in database", "admin count", count)
		return nil
	}
	if len(config.Config.AdminList) == 0 {
		log.ZInfo(ctx, "AdminList is empty", "adminList", config.Config.AdminList)
		return nil
	}
	now := time.Now()
	admins := make([]*admin.Admin, 0, len(config.Config.AdminList))
	for _, adminChat := range config.Config.AdminList {
		password := md5.Sum([]byte(adminChat.AdminID))
		table := admin.Admin{
			Account:    adminChat.AdminID,
			UserID:     adminChat.ImAdminID,
			Password:   hex.EncodeToString(password[:]),
			Level:      100,
			CreateTime: now,
		}
		if adminChat.NickName != "" {
			table.Nickname = adminChat.NickName
		} else {
			table.Nickname = adminChat.AdminID
		}
		admins = append(admins, &table)
	}
	if err := o.db.WithContext(ctx).Create(&admins).Error; err != nil {
		return errs.Wrap(err)
	}
	return nil
}
