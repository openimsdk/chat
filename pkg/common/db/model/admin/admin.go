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
	"crypto/md5"
	"encoding/hex"
	"github.com/OpenIMSDK/chat/pkg/common/constant"
	"github.com/OpenIMSDK/tools/mgoutil"
	"github.com/OpenIMSDK/tools/pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"

	"github.com/OpenIMSDK/chat/pkg/common/config"
	"github.com/OpenIMSDK/chat/pkg/common/db/table/admin"
	"github.com/OpenIMSDK/tools/errs"
)

func NewAdmin(db *mongo.Database) (admin.AdminInterface, error) {
	coll := db.Collection("admin")
	_, err := coll.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{
			{Key: "account", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return &Admin{
		coll: coll,
	}, nil
}

type Admin struct {
	coll *mongo.Collection
}

func (o *Admin) Take(ctx context.Context, account string) (*admin.Admin, error) {
	return mgoutil.FindOne[*admin.Admin](ctx, o.coll, bson.M{"account": account})
}

func (o *Admin) TakeUserID(ctx context.Context, userID string) (*admin.Admin, error) {
	return mgoutil.FindOne[*admin.Admin](ctx, o.coll, bson.M{"user_id": userID})
}

func (o *Admin) Update(ctx context.Context, account string, update map[string]any) error {
	if len(update) == 0 {
		return nil
	}
	return mgoutil.UpdateOne(ctx, o.coll, bson.M{"user_id": account}, bson.M{"$set": update}, false)
}

func (o *Admin) Create(ctx context.Context, admins []*admin.Admin) error {
	return mgoutil.InsertMany(ctx, o.coll, admins)
}

func (o *Admin) ChangePassword(ctx context.Context, userID string, newPassword string) error {
	return mgoutil.UpdateOne(ctx, o.coll, bson.M{"user_id": userID}, bson.M{"$set": bson.M{"password": newPassword}}, false)

}

func (o *Admin) Delete(ctx context.Context, userIDs []string) error {
	if len(userIDs) == 0 {
		return nil
	}
	return mgoutil.DeleteMany(ctx, o.coll, bson.M{"user_id": bson.M{"$in": userIDs}})
}

func (o *Admin) Search(ctx context.Context, pagination pagination.Pagination) (int64, []*admin.Admin, error) {
	opt := options.Find().SetSort(bson.D{{"create_time", -1}})
	filter := bson.M{"level": constant.NormalAdmin}
	return mgoutil.FindPage[*admin.Admin](ctx, o.coll, filter, pagination, opt)
}

func (o *Admin) InitAdmin(ctx context.Context) error {
	filter := bson.M{}
	count, err := mgoutil.Count(ctx, o.coll, filter)
	if err != nil {
		return errs.Wrap(err)
	}
	if count > 0 {
		return nil
	}
	if len(config.Config.ChatAdmin) == 0 {
		return nil
	}

	admins := make([]*admin.Admin, 0, len(config.Config.ChatAdmin))
	o.createAdmins(&admins, config.Config.ChatAdmin)

	return mgoutil.InsertMany(ctx, o.coll, admins)
}

func (o *Admin) createAdmins(adminList *[]*admin.Admin, registerList []config.Admin) {
	// chatAdmin set the level to 50, this account use for send notification.
	for _, adminChat := range registerList {
		table := admin.Admin{
			Account:    adminChat.AdminID,
			UserID:     adminChat.ImAdminID,
			Password:   o.passwordEncryption(adminChat.AdminID),
			Level:      100,
			CreateTime: time.Now(),
		}
		if adminChat.NickName != "" {
			table.Nickname = adminChat.NickName
		} else {
			table.Nickname = adminChat.AdminID
		}
		*adminList = append(*adminList, &table)
	}
}

func (o *Admin) passwordEncryption(password string) string {
	paswd := md5.Sum([]byte(password))
	return hex.EncodeToString(paswd[:])
}
