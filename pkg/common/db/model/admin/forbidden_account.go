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

func NewForbiddenAccount(db *mongo.Database) (admin.ForbiddenAccountInterface, error) {
	coll := db.Collection("forbidden_account")
	_, err := coll.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{
			{Key: "user_id", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return &ForbiddenAccount{
		coll: coll,
	}, nil
}

type ForbiddenAccount struct {
	coll *mongo.Collection
}

func (o *ForbiddenAccount) Create(ctx context.Context, ms []*admin.ForbiddenAccount) error {
	return mgoutil.InsertMany(ctx, o.coll, ms)
}

func (o *ForbiddenAccount) Take(ctx context.Context, userID string) (*admin.ForbiddenAccount, error) {
	return mgoutil.FindOne[*admin.ForbiddenAccount](ctx, o.coll, bson.M{"user_id": userID})
}

func (o *ForbiddenAccount) Delete(ctx context.Context, userIDs []string) error {
	if len(userIDs) == 0 {
		return nil
	}
	return mgoutil.DeleteMany(ctx, o.coll, bson.M{"user_id": bson.M{"$in": userIDs}})
}

func (o *ForbiddenAccount) Find(ctx context.Context, userIDs []string) ([]*admin.ForbiddenAccount, error) {
	return mgoutil.Find[*admin.ForbiddenAccount](ctx, o.coll, bson.M{"user_id": bson.M{"$in": userIDs}})
}

func (o *ForbiddenAccount) Search(ctx context.Context, keyword string, pagination pagination.Pagination) (int64, []*admin.ForbiddenAccount, error) {
	var filter bson.M

	if keyword != "" {
		filter = bson.M{
			"$or": []bson.M{
				{"user_id": bson.M{"$regex": keyword, "$options": "i"}},
				{"reason": bson.M{"$regex": keyword, "$options": "i"}},
				{"operator_user_id": bson.M{"$regex": keyword, "$options": "i"}},
			},
		}
	}
	return mgoutil.FindPage[*admin.ForbiddenAccount](ctx, o.coll, filter, pagination)
}

func (o *ForbiddenAccount) FindAllIDs(ctx context.Context) ([]string, error) {
	return mgoutil.Find[string](ctx, o.coll, bson.M{}, options.Find().SetProjection(bson.M{"_id": 0, "user_id": 1}))
}
