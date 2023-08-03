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

	"github.com/OpenIMSDK/tools/errs"
	"gorm.io/gorm"

	"github.com/OpenIMSDK/chat/pkg/common/db/table/chat"
)

func NewRegister(db *gorm.DB) chat.RegisterInterface {
	return &Register{db: db}
}

type Register struct {
	db *gorm.DB
}

func (o *Register) NewTx(tx any) chat.RegisterInterface {
	return &Register{db: tx.(*gorm.DB)}
}

func (o *Register) Create(ctx context.Context, registers ...*chat.Register) error {
	return errs.Wrap(o.db.WithContext(ctx).Create(registers).Error)
}

func (o *Register) CountTotal(ctx context.Context, before *time.Time) (count int64, err error) {
	db := o.db.WithContext(ctx).Model(&chat.Register{})
	if before != nil {
		db.Where("create_time < ?", before)
	}
	if err := db.Count(&count).Error; err != nil {
		return 0, errs.Wrap(err)
	}
	return count, nil
}
