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

	"github.com/OpenIMSDK/tools/mgoutil"
	"github.com/OpenIMSDK/tools/pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/OpenIMSDK/chat/pkg/common/db/table/admin"
	"github.com/OpenIMSDK/tools/errs"
)

func NewRegisterAddFriend(db *mongo.Database) (admin.RegisterAddFriendInterface, error) {
	coll := db.Collection("register_add_friend")
	_, err := coll.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{
			{Key: "user_id", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return &RegisterAddFriend{
		coll: coll,
	}, nil
}

type RegisterAddFriend struct {
	coll *mongo.Collection
}

func (o *RegisterAddFriend) Add(ctx context.Context, registerAddFriends []*admin.RegisterAddFriend) error {
	return mgoutil.InsertMany(ctx, o.coll, registerAddFriends)
}

func (o *RegisterAddFriend) Del(ctx context.Context, userIDs []string) error {
	if len(userIDs) == 0 {
		return nil
	}
	return mgoutil.DeleteMany(ctx, o.coll, bson.M{"user_id": bson.M{"$in": userIDs}})
}

func (o *RegisterAddFriend) FindUserID(ctx context.Context, userIDs []string) ([]string, error) {
	filter := bson.M{}
	if len(userIDs) > 0 {
		filter["user_id"] = bson.M{"$in": userIDs}
	}
	return mgoutil.Find[string](ctx, o.coll, filter, options.Find().SetProjection(bson.M{"_id": 0, "user_id": 1}))
}

func (o *RegisterAddFriend) Search(ctx context.Context, keyword string, pagination pagination.Pagination) (int64, []*admin.RegisterAddFriend, error) {
	filter := bson.M{"user_id": bson.M{"$regex": keyword, "$options": "i"}}
	return mgoutil.FindPage[*admin.RegisterAddFriend](ctx, o.coll, filter, pagination)
}
