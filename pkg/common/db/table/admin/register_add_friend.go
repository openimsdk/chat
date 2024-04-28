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

package admin

import (
	"context"
	"github.com/openimsdk/tools/db/pagination"
	"time"
)

type RegisterAddFriend struct {
	UserID     string    `bson:"user_id"`
	CreateTime time.Time `bson:"create_time"`
}

func (RegisterAddFriend) TableName() string {
	return "register_add_friends"
}

type RegisterAddFriendInterface interface {
	Add(ctx context.Context, registerAddFriends []*RegisterAddFriend) error
	Del(ctx context.Context, userIDs []string) error
	FindUserID(ctx context.Context, userIDs []string) ([]string, error)
	Search(ctx context.Context, keyword string, pagination pagination.Pagination) (int64, []*RegisterAddFriend, error)
}
