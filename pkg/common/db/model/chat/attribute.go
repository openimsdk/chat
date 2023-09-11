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

package chat

import (
	"context"

	"github.com/OpenIMSDK/tools/errs"
	"github.com/OpenIMSDK/tools/ormutil"
	"gorm.io/gorm"

	"github.com/OpenIMSDK/chat/pkg/common/db/table/chat"
)

func NewAttribute(db *gorm.DB) chat.AttributeInterface {
	return &Attribute{db: db}
}

type Attribute struct {
	db *gorm.DB
}

func (o *Attribute) NewTx(tx any) chat.AttributeInterface {
	return &Attribute{db: tx.(*gorm.DB)}
}

func (o *Attribute) Create(ctx context.Context, attribute ...*chat.Attribute) error {
	return errs.Wrap(o.db.WithContext(ctx).Create(attribute).Error)
}

func (o *Attribute) Update(ctx context.Context, userID string, data map[string]any) error {
	return errs.Wrap(o.db.WithContext(ctx).Model(&chat.Attribute{}).Where("user_id = ?", userID).Updates(data).Error)
}

func (o *Attribute) Find(ctx context.Context, userIds []string) ([]*chat.Attribute, error) {
	var a []*chat.Attribute
	return a, errs.Wrap(o.db.WithContext(ctx).Where("user_id in (?)", userIds).Find(&a).Error)
}

func (o *Attribute) FindAccount(ctx context.Context, accounts []string) ([]*chat.Attribute, error) {
	var a []*chat.Attribute
	return a, errs.Wrap(o.db.WithContext(ctx).Where("account in (?)", accounts).Find(&a).Error)
}

func (o *Attribute) Search(ctx context.Context, keyword string, genders []int32, page int32, size int32) (uint32, []*chat.Attribute, error) {
	db := o.db.WithContext(ctx)
	if len(genders) > 0 {
		db = db.Where("gender in ?", genders)
	}
	return ormutil.GormSearch[chat.Attribute](db, []string{"user_id", "account", "nickname", "phone_number"}, keyword, page, size)
}

func (o *Attribute) TakePhone(ctx context.Context, areaCode string, phoneNumber string) (*chat.Attribute, error) {
	var a chat.Attribute
	return &a, errs.Wrap(o.db.WithContext(ctx).Where("area_code = ? and phone_number = ?", areaCode, phoneNumber).First(&a).Error)
}

func (o *Attribute) TakeAccount(ctx context.Context, account string) (*chat.Attribute, error) {
	var a chat.Attribute
	return &a, errs.Wrap(o.db.WithContext(ctx).Where("account = ?", account).Take(&a).Error)
}

func (o *Attribute) Take(ctx context.Context, userID string) (*chat.Attribute, error) {
	var a chat.Attribute
	return &a, errs.Wrap(o.db.WithContext(ctx).Where("user_id = ?", userID).Take(&a).Error)
}

func (o *Attribute) SearchNormalUser(ctx context.Context, keyword string, forbiddenIDs []string, gender int32, page int32, size int32) (uint32, []*chat.Attribute, error) {
	db := o.db.WithContext(ctx)
	var genders []int32
	if gender == 0 {
		genders = append(genders, 0, 1, 2)
	} else {
		genders = append(genders, gender)
	}
	db = db.Where("gender in ?", genders)
	if len(forbiddenIDs) > 0 {
		db = db.Where("user_id not in ?", forbiddenIDs)
	}
	return ormutil.GormSearch[chat.Attribute](db, []string{"user_id", "account", "nickname", "phone_number"}, keyword, page, size)
}

func (o *Attribute) SearchUser(ctx context.Context, keyword string, userIDs []string, genders []int32, pageNumber int32, showNumber int32) (uint32, []*chat.Attribute, error) {
	db := o.db.WithContext(ctx)
	ormutil.GormIn(&db, "user_id", userIDs)
	ormutil.GormIn(&db, "gender", genders)
	return ormutil.GormSearch[chat.Attribute](db, []string{"user_id", "nickname", "phone_number"}, keyword, pageNumber, showNumber)
}
