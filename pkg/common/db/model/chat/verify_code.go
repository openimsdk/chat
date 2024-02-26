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
	"github.com/OpenIMSDK/tools/mgoutil"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"

	"github.com/OpenIMSDK/chat/pkg/common/db/table/chat"
	"github.com/OpenIMSDK/tools/errs"
)

func NewVerifyCode(db *mongo.Database) (chat.VerifyCodeInterface, error) {
	coll := db.Collection("verify_code")
	_, err := coll.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "account", Value: 1},
			},
		},
	})
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return &VerifyCode{
		coll: coll,
	}, nil
}

type VerifyCode struct {
	coll *mongo.Collection
}

func (o *VerifyCode) parseID(s string) (primitive.ObjectID, error) {
	objID, err := primitive.ObjectIDFromHex(s)
	if err != nil {
		var zero primitive.ObjectID
		return zero, errs.Wrap(err)
	}
	return objID, nil
}

func (o *VerifyCode) Add(ctx context.Context, ms []*chat.VerifyCode) error {
	type MongoVerifyCode struct {
		ID         primitive.ObjectID `bson:"_id"`
		Account    string             `bson:"account"`
		Platform   string             `bson:"platform"`
		Code       string             `bson:"code"`
		Duration   uint               `bson:"duration"`
		Count      int                `bson:"count"`
		Used       bool               `bson:"used"`
		CreateTime time.Time          `bson:"create_time"`
	}
	tmp := make([]MongoVerifyCode, 0, len(ms))
	for _, m := range ms {
		objID, err := o.parseID(m.ID)
		if err != nil {
			return err
		}
		tmp = append(tmp, MongoVerifyCode{
			ID:         objID,
			Account:    m.Account,
			Platform:   m.Platform,
			Code:       m.Code,
			Duration:   m.Duration,
			Count:      m.Count,
			Used:       m.Used,
			CreateTime: m.CreateTime,
		})
	}
	return mgoutil.InsertMany(ctx, o.coll, tmp)
}

func (o *VerifyCode) RangeNum(ctx context.Context, account string, start time.Time, end time.Time) (int64, error) {
	filter := bson.M{
		"account": account,
		"create_time": bson.M{
			"$gte": start,
			"$lte": end,
		},
	}
	return mgoutil.Count(ctx, o.coll, filter)
}

func (o *VerifyCode) TakeLast(ctx context.Context, account string) (*chat.VerifyCode, error) {
	filter := bson.M{
		"account": account,
	}
	opt := options.FindOne().SetSort(bson.M{"id": -1})
	return mgoutil.FindOne[*chat.VerifyCode](ctx, o.coll, filter, opt)
}

func (o *VerifyCode) Incr(ctx context.Context, id string) error {
	objID, err := o.parseID(id)
	if err != nil {
		return err
	}
	return mgoutil.UpdateOne(ctx, o.coll, bson.M{"_id": objID}, bson.M{"$inc": bson.M{"count": 1}}, false)
}

func (o *VerifyCode) Delete(ctx context.Context, id string) error {
	objID, err := o.parseID(id)
	if err != nil {
		return err
	}
	return mgoutil.DeleteOne(ctx, o.coll, bson.M{"_id": objID})
}
