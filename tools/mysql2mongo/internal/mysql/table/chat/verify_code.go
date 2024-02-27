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
)

type VerifyCode struct {
	ID         uint      `gorm:"column:id;primary_key;autoIncrement"`
	Account    string    `gorm:"column:account;type:char(64)"`
	Platform   string    `gorm:"column:platform;type:varchar(32)"`
	Code       string    `gorm:"column:verify_code;type:varchar(16)"`
	Duration   uint      `gorm:"column:duration;type:int(11)"`
	Count      int       `gorm:"column:count;type:int(11)"`
	Used       bool      `gorm:"column:used"`
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime"`
}

func (VerifyCode) TableName() string {
	return "verify_codes"
}

type VerifyCodeInterface interface {
	NewTx(tx any) VerifyCodeInterface
	Add(ctx context.Context, ms []*VerifyCode) error
	RangeNum(ctx context.Context, account string, start time.Time, end time.Time) (uint32, error)
	TakeLast(ctx context.Context, account string) (*VerifyCode, error)
	Incr(ctx context.Context, id uint) error
	Delete(ctx context.Context, id uint) error
}
