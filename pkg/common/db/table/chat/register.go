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

type Register struct {
	UserID      string    `bson:"user_id"`
	DeviceID    string    `bson:"device_id"`
	IP          string    `bson:"ip"`
	Platform    string    `bson:"platform"`
	AccountType string    `bson:"account_type"`
	Mode        string    `bson:"mode"`
	CreateTime  time.Time `bson:"create_time"`
}

func (Register) TableName() string {
	return "registers"
}

type RegisterInterface interface {
	// NewTx(tx any) RegisterInterface
	Create(ctx context.Context, registers ...*Register) error
	CountTotal(ctx context.Context, before *time.Time) (int64, error)
	Delete(ctx context.Context, userIDs []string) error
}
