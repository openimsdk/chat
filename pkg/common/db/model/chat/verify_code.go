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

func NewVerifyCode(db *gorm.DB) *VerifyCode {
	return &VerifyCode{
		db: db,
	}
}

type VerifyCode struct {
	db *gorm.DB
}

func (o *VerifyCode) NewTx(tx any) chat.VerifyCodeInterface {
	return &VerifyCode{db: tx.(*gorm.DB)}
}

func (o *VerifyCode) Add(ctx context.Context, ms []*chat.VerifyCode) error {
	return errs.Wrap(o.db.WithContext(ctx).Create(&ms).Error)
}

func (o *VerifyCode) RangeNum(ctx context.Context, account string, start time.Time, end time.Time) (uint32, error) {
	var count int64
	if err := o.db.WithContext(ctx).Model(&chat.VerifyCode{}).Where("account = ?", account).Where("create_time BETWEEN ? AND ?", start, end).Count(&count).Error; err != nil {
		return 0, errs.Wrap(err)
	}
	return uint32(count), nil
}

func (o *VerifyCode) TakeLast(ctx context.Context, account string) (*chat.VerifyCode, error) {
	var m chat.VerifyCode
	return &m, errs.Wrap(o.db.WithContext(ctx).Where("account = ?", account).Order("id DESC").Take(&m).Error)
}

func (o *VerifyCode) Incr(ctx context.Context, id uint) error {
	return errs.Wrap(o.db.WithContext(ctx).Model(&chat.VerifyCode{}).Where("id = ?", id).Updates(map[string]any{"count": gorm.Expr("count + 1")}).Error)
}

func (o *VerifyCode) Delete(ctx context.Context, id uint) error {
	return errs.Wrap(o.db.WithContext(ctx).Where("id = ?", id).Delete(&chat.VerifyCode{}).Error)
}
