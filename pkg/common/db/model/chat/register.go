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

	"github.com/openimsdk/tools/db/mongoutil"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/openimsdk/chat/pkg/common/db/table/chat"
	"github.com/openimsdk/tools/errs"
)

func NewRegister(db *mongo.Database) (chat.RegisterInterface, error) {
	coll := db.Collection("register")
	_, err := coll.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "user_id", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
	})
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return &Register{coll: coll}, nil
}

type Register struct {
	coll *mongo.Collection
}

func (o *Register) Create(ctx context.Context, registers ...*chat.Register) error {
	return mongoutil.InsertMany(ctx, o.coll, registers)
}

func (o *Register) CountTotal(ctx context.Context, before *time.Time) (int64, error) {
	filter := bson.M{}
	if before != nil {
		filter["create_time"] = bson.M{"$lt": before}
	}
	return mongoutil.Count(ctx, o.coll, filter)
}

func (o *Register) Delete(ctx context.Context, userIDs []string) error {
	if len(userIDs) == 0 {
		return nil
	}
	return mongoutil.DeleteMany(ctx, o.coll, bson.M{"user_id": bson.M{"$in": userIDs}})
}
