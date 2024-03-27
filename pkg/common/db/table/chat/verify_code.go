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
	ID         string    `bson:"_id"`
	Account    string    `bson:"account"`
	Platform   string    `bson:"platform"`
	Code       string    `bson:"code"`
	Duration   uint      `bson:"duration"`
	Count      int       `bson:"count"`
	Used       bool      `bson:"used"`
	CreateTime time.Time `bson:"create_time"`
}

func (VerifyCode) TableName() string {
	return "verify_codes"
}

type VerifyCodeInterface interface {
	Add(ctx context.Context, ms []*VerifyCode) error
	RangeNum(ctx context.Context, account string, start time.Time, end time.Time) (int64, error)
	TakeLast(ctx context.Context, account string) (*VerifyCode, error)
	Incr(ctx context.Context, id string) error
	Delete(ctx context.Context, id string) error
}
