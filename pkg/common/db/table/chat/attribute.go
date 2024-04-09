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
	"github.com/openimsdk/tools/db/pagination"
	"time"
)

// Attribute 用户属性表.
type Attribute struct {
	UserID           string    `bson:"user_id"`
	Account          string    `bson:"account"`
	PhoneNumber      string    `bson:"phone_number"`
	AreaCode         string    `bson:"area_code"`
	Email            string    `bson:"email"`
	Nickname         string    `bson:"nickname"`
	FaceURL          string    `bson:"face_url"`
	Gender           int32     `bson:"gender"`
	CreateTime       time.Time `bson:"create_time"`
	ChangeTime       time.Time `bson:"change_time"`
	BirthTime        time.Time `bson:"birth_time"`
	Level            int32     `bson:"level"`
	AllowVibration   int32     `bson:"allow_vibration"`
	AllowBeep        int32     `bson:"allow_beep"`
	AllowAddFriend   int32     `bson:"allow_add_friend"`
	GlobalRecvMsgOpt int32     `bson:"global_recv_msg_opt"`
	RegisterType     int32     `bson:"register_type"`
}

func (Attribute) TableName() string {
	return "attributes"
}

type AttributeInterface interface {
	//NewTx(tx any) AttributeInterface
	Create(ctx context.Context, attribute ...*Attribute) error
	Update(ctx context.Context, userID string, data map[string]any) error
	Find(ctx context.Context, userIds []string) ([]*Attribute, error)
	FindAccount(ctx context.Context, accounts []string) ([]*Attribute, error)
	Search(ctx context.Context, keyword string, genders []int32, pagination pagination.Pagination) (int64, []*Attribute, error)
	TakePhone(ctx context.Context, areaCode string, phoneNumber string) (*Attribute, error)
	TakeEmail(ctx context.Context, email string) (*Attribute, error)
	TakeAccount(ctx context.Context, account string) (*Attribute, error)
	Take(ctx context.Context, userID string) (*Attribute, error)
	SearchNormalUser(ctx context.Context, keyword string, forbiddenID []string, gender int32, pagination pagination.Pagination) (int64, []*Attribute, error)
	SearchUser(ctx context.Context, keyword string, userIDs []string, genders []int32, pagination pagination.Pagination) (int64, []*Attribute, error)
}
