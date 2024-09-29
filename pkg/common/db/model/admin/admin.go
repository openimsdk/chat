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

	"github.com/openimsdk/chat/pkg/common/constant"
	admindb "github.com/openimsdk/chat/pkg/common/db/table/admin"
	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/db/pagination"
	"github.com/openimsdk/tools/errs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewAdmin(db *mongo.Database) (admindb.AdminInterface, error) {
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

func (o *Admin) Take(ctx context.Context, account string) (*admindb.Admin, error) {
	return mongoutil.FindOne[*admindb.Admin](ctx, o.coll, bson.M{"account": account})
}

func (o *Admin) TakeUserID(ctx context.Context, userID string) (*admindb.Admin, error) {
	return mongoutil.FindOne[*admindb.Admin](ctx, o.coll, bson.M{"user_id": userID})
}

func (o *Admin) Update(ctx context.Context, account string, update map[string]any) error {
	if len(update) == 0 {
		return nil
	}
	return mongoutil.UpdateOne(ctx, o.coll, bson.M{"user_id": account}, bson.M{"$set": update}, false)
}

func (o *Admin) Create(ctx context.Context, admins []*admindb.Admin) error {
	return mongoutil.InsertMany(ctx, o.coll, admins)
}

func (o *Admin) ChangePassword(ctx context.Context, userID string, newPassword string) error {
	return mongoutil.UpdateOne(ctx, o.coll, bson.M{"user_id": userID}, bson.M{"$set": bson.M{"password": newPassword}}, false)
}

func (o *Admin) Delete(ctx context.Context, userIDs []string) error {
	if len(userIDs) == 0 {
		return nil
	}
	return mongoutil.DeleteMany(ctx, o.coll, bson.M{"user_id": bson.M{"$in": userIDs}})
}

func (o *Admin) Search(ctx context.Context, pagination pagination.Pagination) (int64, []*admindb.Admin, error) {
	opt := options.Find().SetSort(bson.D{{Key: "create_time", Value: -1}})
	filter := bson.M{"level": constant.NormalAdmin}
	return mongoutil.FindPage[*admindb.Admin](ctx, o.coll, filter, pagination, opt)
}
