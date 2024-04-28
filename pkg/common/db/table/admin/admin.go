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

// Admin user
type Admin struct {
	Account    string    `bson:"account"`
	Password   string    `bson:"password"`
	FaceURL    string    `bson:"face_url"`
	Nickname   string    `bson:"nickname"`
	UserID     string    `bson:"user_id"`
	Level      int32     `bson:"level"`
	CreateTime time.Time `bson:"create_time"`
}

func (Admin) TableName() string {
	return "admins"
}

type AdminInterface interface {
	Create(ctx context.Context, admins []*Admin) error
	Take(ctx context.Context, account string) (*Admin, error)
	TakeUserID(ctx context.Context, userID string) (*Admin, error)
	Update(ctx context.Context, account string, update map[string]any) error
	ChangePassword(ctx context.Context, userID string, newPassword string) error
	Delete(ctx context.Context, userIDs []string) error
	Search(ctx context.Context, pagination pagination.Pagination) (int64, []*Admin, error)
}
