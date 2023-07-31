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
	"gorm.io/gorm"

	"github.com/OpenIMSDK/chat/pkg/common/db/table/chat"
)

func NewUserLoginRecord(db *gorm.DB) chat.UserLoginRecordInterface {
	return &UserLoginRecord{
		db: db,
	}
}

type UserLoginRecord struct {
	db *gorm.DB
}

func (o *UserLoginRecord) NewTx(tx any) chat.UserLoginRecordInterface {
	return &UserLoginRecord{db: tx.(*gorm.DB)}
}

func (o *UserLoginRecord) Create(ctx context.Context, records ...*chat.UserLoginRecord) error {
	return o.db.WithContext(ctx).Create(&records).Error
}
