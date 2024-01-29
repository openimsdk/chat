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
	"github.com/OpenIMSDK/chat/pkg/common/constant"
	"time"

	"github.com/OpenIMSDK/tools/log"

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

func (o *Admin) Create(ctx context.Context, admin *admin.Admin) error {
	return errs.Wrap(o.db.WithContext(ctx).Create(&admin).Error)
}

func (o *Admin) ChangePassword(ctx context.Context, userID string, newPassword string) error {
	return errs.Wrap(o.db.WithContext(ctx).Model(&admin.Admin{}).Where("user_id=?", userID).Update("password", newPassword).Error)
}

func (o *Admin) Delete(ctx context.Context, userIDs []string) error {
	return errs.Wrap(o.db.WithContext(ctx).Where("user_id in ?", userIDs).Delete(&admin.Admin{}).Error)
}

func (o *Admin) Search(ctx context.Context, page, size int32) (uint32, []*admin.Admin, error) {
	var count int64
	var admins []*admin.Admin
	if err := o.db.WithContext(ctx).Model(&admin.Admin{}).Where("level=?", constant.NormalAdmin).Count(&count).Error; err != nil {
		return 0, nil, errs.Wrap(err)
	}
	offset := (page - 1) * size
	if err := o.db.WithContext(ctx).Order("create_time desc").Offset(int(offset)).Where("level=?", constant.NormalAdmin).Limit(int(size)).Find(&admins).Error; err != nil {
		return 0, nil, errs.Wrap(err)
	}
	return uint32(count), admins, nil
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
	if len(config.Config.ChatAdmin) == 0 {
		log.ZInfo(ctx, "ChatAdmin is empty", "ChatAdmin", config.Config.ChatAdmin)
		return nil
	}

	admins := make([]*admin.Admin, 0, len(config.Config.ChatAdmin))
	o.createAdmins(&admins, config.Config.ChatAdmin)

	if err := o.db.WithContext(ctx).Create(&admins).Error; err != nil {
		return errs.Wrap(err)
	}
	return nil
}

func (o *Admin) createAdmins(adminList *[]*admin.Admin, registerList []config.Admin) {
	// chatAdmin set the level to 50, this account use for send notification.
	for _, adminChat := range registerList {
		table := admin.Admin{
			Account:    adminChat.AdminID,
			UserID:     adminChat.ImAdminID,
			Password:   o.passwordEncryption(adminChat.AdminID),
			Level:      100,
			CreateTime: time.Now(),
		}
		if adminChat.NickName != "" {
			table.Nickname = adminChat.NickName
		} else {
			table.Nickname = adminChat.AdminID
		}
		*adminList = append(*adminList, &table)
	}
}

func (o *Admin) passwordEncryption(password string) string {
	paswd := md5.Sum([]byte(password))
	return hex.EncodeToString(paswd[:])
}
