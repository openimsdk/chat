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
	"time"

	"github.com/OpenIMSDK/chat/pkg/common/db/table/chat"
	"github.com/OpenIMSDK/tools/errs"
	"gorm.io/gorm"
)

func NewAccount(db *gorm.DB) chat.AccountInterface {
	return &Account{db: db}
}

type Account struct {
	db *gorm.DB
}

func (o *Account) NewTx(tx any) chat.AccountInterface {
	return &Account{db: tx.(*gorm.DB)}
}

func (o *Account) Create(ctx context.Context, accounts ...*chat.Account) error {
	return errs.Wrap(o.db.WithContext(ctx).Create(&accounts).Error)
}

func (o *Account) Take(ctx context.Context, userId string) (*chat.Account, error) {
	var a chat.Account
	return &a, errs.Wrap(o.db.WithContext(ctx).Where("user_id = ?", userId).First(&a).Error)
}

func (o *Account) Update(ctx context.Context, userID string, data map[string]any) error {
	return errs.Wrap(o.db.WithContext(ctx).Model(&chat.Account{}).Where("user_id = ?", userID).Updates(data).Error)
}

func (o *Account) UpdatePassword(ctx context.Context, userId string, password string) error {
	return errs.Wrap(o.db.WithContext(ctx).Model(&chat.Account{}).Where("user_id = ?", userId).Updates(map[string]interface{}{"password": password, "change_time": time.Now()}).Error)
}
