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

	"github.com/OpenIMSDK/tools/errs"
	"github.com/OpenIMSDK/tools/mgoutil"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/OpenIMSDK/chat/pkg/common/db/table/chat"
)

func NewAccount(db *mongo.Database) (chat.AccountInterface, error) {
	coll := db.Collection("account")
	_, err := coll.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{
			{Key: "user_id", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return &Account{coll: coll}, nil
}

type Account struct {
	coll *mongo.Collection
}

func (o *Account) Create(ctx context.Context, accounts ...*chat.Account) error {
	return mgoutil.InsertMany(ctx, o.coll, accounts)
}

func (o *Account) Take(ctx context.Context, userId string) (*chat.Account, error) {
	return mgoutil.FindOne[*chat.Account](ctx, o.coll, bson.M{"user_id": userId})
}

func (o *Account) Update(ctx context.Context, userID string, data map[string]any) error {
	if len(data) == 0 {
		return nil
	}
	return mgoutil.UpdateOne(ctx, o.coll, bson.M{"user_id": userID}, bson.M{"$set": data}, false)
}

func (o *Account) UpdatePassword(ctx context.Context, userId string, password string) error {
	return mgoutil.UpdateOne(ctx, o.coll, bson.M{"user_id": userId}, bson.M{"$set": bson.M{"password": password, "change_time": time.Now()}}, false)
}
