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
	"github.com/OpenIMSDK/tools/pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/OpenIMSDK/chat/pkg/common/db/table/chat"
)

func NewLogs(db *mongo.Database) (chat.LogInterface, error) {
	coll := db.Collection("chat_logs")
	_, err := coll.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "log_id", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{
				{Key: "user_id", Value: 1},
			},
		},
	})
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return &Logs{coll: coll}, nil
}

type Logs struct {
	coll *mongo.Collection
}

func (l *Logs) Create(ctx context.Context, log []*chat.Log) error {
	return mgoutil.InsertMany(ctx, l.coll, log)
}

func (l *Logs) Search(ctx context.Context, keyword string, start time.Time, end time.Time, pagination pagination.Pagination) (int64, []*chat.Log, error) {
	filter := bson.M{}
	if end.UnixMilli() == 0 {
		filter["create_time"] = bson.M{
			"$gte": start,
		}
	} else {
		filter["create_time"] = bson.M{
			"$gte": start,
			"$lte": end,
		}
	}
	if keyword != "" {
		filter["user_id"] = bson.M{"$regex": keyword, "$options": "i"}
	}
	return mgoutil.FindPage[*chat.Log](ctx, l.coll, filter, pagination)
}

func (l *Logs) Delete(ctx context.Context, logIDs []string, userID string) error {
	if len(logIDs) == 0 {
		return nil
	}
	if userID == "" {
		return mgoutil.DeleteMany(ctx, l.coll, bson.M{"log_id": bson.M{"$in": logIDs}})
	}
	return mgoutil.DeleteMany(ctx, l.coll, bson.M{"log_id": bson.M{"$in": logIDs}, "user_id": userID})
}

func (l *Logs) Get(ctx context.Context, logIDs []string, userID string) ([]*chat.Log, error) {
	if userID == "" {
		return mgoutil.Find[*chat.Log](ctx, l.coll, bson.M{"log_id": bson.M{"$in": logIDs}})
	}
	return mgoutil.Find[*chat.Log](ctx, l.coll, bson.M{"log_id": bson.M{"$in": logIDs}, "user_id": userID})
}
