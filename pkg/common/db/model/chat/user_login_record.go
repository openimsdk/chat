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
	"github.com/openimsdk/tools/db/mongoutil"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"

	"github.com/openimsdk/chat/pkg/common/db/table/chat"
)

func NewUserLoginRecord(db *mongo.Database) (chat.UserLoginRecordInterface, error) {
	coll := db.Collection("user_login_record")
	_, err := coll.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "create_time", Value: 1},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return &UserLoginRecord{
		coll: coll,
	}, nil
}

type UserLoginRecord struct {
	coll *mongo.Collection
}

func (o *UserLoginRecord) Create(ctx context.Context, records ...*chat.UserLoginRecord) error {
	return mongoutil.InsertMany(ctx, o.coll, records)
}

func (o *UserLoginRecord) CountTotal(ctx context.Context, before *time.Time) (count int64, err error) {
	filter := bson.M{}
	if before != nil {
		filter["create_time"] = bson.M{"$lt": before}
	}
	return mongoutil.Count(ctx, o.coll, filter)
}

func (o *UserLoginRecord) CountRangeEverydayTotal(ctx context.Context, start *time.Time, end *time.Time) (map[string]int64, int64, error) {
	pipeline := make([]bson.M, 0, 4)
	if start != nil || end != nil {
		filter := bson.M{}
		if start != nil {
			filter["$gte"] = start
		}
		if end != nil {
			filter["$lt"] = end
		}
		pipeline = append(pipeline, bson.M{"$match": bson.M{"login_time": filter}})
	}
	pipeline = append(pipeline,
		bson.M{
			"$project": bson.M{
				"_id":     0,
				"user_id": 1,
				"login_time": bson.M{
					"$dateToString": bson.M{
						"format": "%Y-%m-%d",
						"date":   "$login_time",
					},
				},
			},
		},

		bson.M{
			"$group": bson.M{
				"_id": bson.M{
					"user_id":    "$user_id",
					"login_time": "$login_time",
				},
			},
		},

		bson.M{
			"$group": bson.M{
				"_id": "$_id.login_time",
				"count": bson.M{
					"$sum": 1,
				},
			},
		},
	)

	type Temp struct {
		ID    string `bson:"_id"`
		Count int64  `bson:"count"`
	}
	res, err := mongoutil.Aggregate[Temp](ctx, o.coll, pipeline)
	if err != nil {
		return nil, 0, err
	}
	var loginCount int64
	countMap := make(map[string]int64, len(res))
	for _, r := range res {
		loginCount += r.Count
		countMap[r.ID] = r.Count
	}
	return countMap, loginCount, nil
}
