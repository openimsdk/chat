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

type Account struct {
	UserID         string    `bson:"user_id"`
	Password       string    `bson:"password"`
	CreateTime     time.Time `bson:"create_time"`
	ChangeTime     time.Time `bson:"change_time"`
	OperatorUserID string    `bson:"operator_user_id"`
}

func (Account) TableName() string {
	return "accounts"
}

type AccountInterface interface {
	Create(ctx context.Context, accounts ...*Account) error
	Take(ctx context.Context, userId string) (*Account, error)
	Update(ctx context.Context, userID string, data map[string]any) error
	UpdatePassword(ctx context.Context, userId string, password string) error
	Delete(ctx context.Context, userIDs []string) error
}
