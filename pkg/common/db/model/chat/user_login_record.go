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

func (o *UserLoginRecord) CountTotal(ctx context.Context, before *time.Time) (count int64, err error) {
	db := o.db.WithContext(ctx).Model(&chat.UserLoginRecord{})
	if before != nil {
		db.Where("create_time < ?", before)
	}
	if err := db.Count(&count).Error; err != nil {
		return 0, errs.Wrap(err)
	}
	return count, nil
}

func (o *UserLoginRecord) CountRangeEverydayTotal(ctx context.Context, start *time.Time, end *time.Time) (map[string]int64, int64, error) {
	var res []struct {
		Date  time.Time `gorm:"column:date"`
		Count int64     `gorm:"column:count"`
	}
	var loginCount int64
	err := o.db.WithContext(ctx).
		Model(&chat.UserLoginRecord{}).
		Select("DATE(login_time) AS date, count(distinct(user_id)) AS count").
		Where("login_time >= ? and login_time < ?", start, end).
		Group("date").
		Find(&res).
		Error
	if err != nil {
		return nil, 0, errs.Wrap(err)
	}
	v := make(map[string]int64)
	for _, r := range res {
		loginCount += r.Count
		v[r.Date.Format("2006-01-02")] = r.Count
	}
	return v, loginCount, nil
}
