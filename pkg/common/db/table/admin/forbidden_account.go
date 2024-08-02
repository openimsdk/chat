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

// ForbiddenAccount table
type ForbiddenAccount struct {
	UserID         string    `bson:"user_id"`
	Reason         string    `bson:"reason"`
	OperatorUserID string    `bson:"operator_user_id"`
	CreateTime     time.Time `bson:"create_time"`
}

func (ForbiddenAccount) TableName() string {
	return "forbidden_accounts"
}

type ForbiddenAccountInterface interface {
	Create(ctx context.Context, ms []*ForbiddenAccount) error
	Take(ctx context.Context, userID string) (*ForbiddenAccount, error)
	Delete(ctx context.Context, userIDs []string) error
	Find(ctx context.Context, userIDs []string) ([]*ForbiddenAccount, error)
	Search(ctx context.Context, keyword string, pagination pagination.Pagination) (int64, []*ForbiddenAccount, error)
	FindAllIDs(ctx context.Context) ([]string, error)
}
