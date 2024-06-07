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

	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/db/pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/openimsdk/chat/pkg/common/db/table/admin"
	"github.com/openimsdk/tools/errs"
)

func NewLimitUserLoginIP(db *mongo.Database) (admin.LimitUserLoginIPInterface, error) {
	coll := db.Collection("limit_user_login_ip")
	_, err := coll.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{
			{Key: "user_id", Value: 1},
			{Key: "ip", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return &LimitUserLoginIP{
		coll: coll,
	}, nil
}

type LimitUserLoginIP struct {
	coll *mongo.Collection
}

func (o *LimitUserLoginIP) Create(ctx context.Context, ms []*admin.LimitUserLoginIP) error {
	return mongoutil.InsertMany(ctx, o.coll, ms)
}

func (o *LimitUserLoginIP) Delete(ctx context.Context, ms []*admin.LimitUserLoginIP) error {
	return mongoutil.DeleteMany(ctx, o.coll, o.limitUserLoginIPFilter(ms))
}

func (o *LimitUserLoginIP) Count(ctx context.Context, userID string) (uint32, error) {
	count, err := mongoutil.Count(ctx, o.coll, bson.M{"user_id": userID})
	if err != nil {
		return 0, err
	}
	return uint32(count), nil
}

func (o *LimitUserLoginIP) Take(ctx context.Context, userID string, ip string) (*admin.LimitUserLoginIP, error) {
	return mongoutil.FindOne[*admin.LimitUserLoginIP](ctx, o.coll, bson.M{"user_id": userID, "ip": ip})
}

func (o *LimitUserLoginIP) Search(ctx context.Context, keyword string, pagination pagination.Pagination) (int64, []*admin.LimitUserLoginIP, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"user_id": bson.M{"$regex": keyword, "$options": "i"}},
			{"ip": bson.M{"$regex": keyword, "$options": "i"}},
		},
	}
	return mongoutil.FindPage[*admin.LimitUserLoginIP](ctx, o.coll, filter, pagination)
}

func (o *LimitUserLoginIP) limitUserLoginIPFilter(ips []*admin.LimitUserLoginIP) bson.M {
	if len(ips) == 0 {
		return nil
	}
	or := make(bson.A, 0, len(ips))
	for _, ip := range ips {
		or = append(or, bson.M{
			"user_id": ip.UserID,
			"ip":      ip.IP,
		})
	}
	return bson.M{"$or": or}
}
