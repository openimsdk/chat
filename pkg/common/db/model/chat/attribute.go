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
	"github.com/openimsdk/tools/db/pagination"
	"github.com/openimsdk/tools/errs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/openimsdk/chat/pkg/common/db/table/chat"
)

func NewAttribute(db *mongo.Database) (chat.AttributeInterface, error) {
	coll := db.Collection("attribute")
	_, err := coll.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "user_id", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{
				{Key: "account", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "email", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "area_code", Value: 1},
				{Key: "phone_number", Value: 1},
			},
		},
	})
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return &Attribute{coll: coll}, nil
}

type Attribute struct {
	coll *mongo.Collection
}

func (o *Attribute) Create(ctx context.Context, attribute ...*chat.Attribute) error {
	return mongoutil.InsertMany(ctx, o.coll, attribute)
}

func (o *Attribute) Update(ctx context.Context, userID string, data map[string]any) error {
	if len(data) == 0 {
		return nil
	}
	return mongoutil.UpdateOne(ctx, o.coll, bson.M{"user_id": userID}, bson.M{"$set": data}, false)
}

func (o *Attribute) Find(ctx context.Context, userIds []string) ([]*chat.Attribute, error) {
	return mongoutil.Find[*chat.Attribute](ctx, o.coll, bson.M{"user_id": bson.M{"$in": userIds}})
}

func (o *Attribute) FindAccount(ctx context.Context, accounts []string) ([]*chat.Attribute, error) {
	return mongoutil.Find[*chat.Attribute](ctx, o.coll, bson.M{"account": bson.M{"$in": accounts}})
}

func (o *Attribute) FindPhone(ctx context.Context, phoneNumbers []string) ([]*chat.Attribute, error) {
	return mongoutil.Find[*chat.Attribute](ctx, o.coll, bson.M{"phone_number": bson.M{"$in": phoneNumbers}})
}

func (o *Attribute) Search(ctx context.Context, keyword string, genders []int32, pagination pagination.Pagination) (int64, []*chat.Attribute, error) {
	filter := bson.M{}
	if len(genders) > 0 {
		filter["gender"] = bson.M{
			"$in": genders,
		}
	}
	if keyword != "" {
		filter["$or"] = []bson.M{
			{"user_id": bson.M{"$regex": keyword, "$options": "i"}},
			{"account": bson.M{"$regex": keyword, "$options": "i"}},
			{"nickname": bson.M{"$regex": keyword, "$options": "i"}},
			{"phone_number": bson.M{"$regex": keyword, "$options": "i"}},
		}
	}
	return mongoutil.FindPage[*chat.Attribute](ctx, o.coll, filter, pagination)
}

func (o *Attribute) TakePhone(ctx context.Context, areaCode string, phoneNumber string) (*chat.Attribute, error) {
	return mongoutil.FindOne[*chat.Attribute](ctx, o.coll, bson.M{"area_code": areaCode, "phone_number": phoneNumber})
}

func (o *Attribute) TakeEmail(ctx context.Context, email string) (*chat.Attribute, error) {
	return mongoutil.FindOne[*chat.Attribute](ctx, o.coll, bson.M{"email": email})
}

func (o *Attribute) TakeAccount(ctx context.Context, account string) (*chat.Attribute, error) {
	return mongoutil.FindOne[*chat.Attribute](ctx, o.coll, bson.M{"account": account})
}

func (o *Attribute) Take(ctx context.Context, userID string) (*chat.Attribute, error) {
	return mongoutil.FindOne[*chat.Attribute](ctx, o.coll, bson.M{"user_id": userID})
}

func (o *Attribute) SearchNormalUser(ctx context.Context, keyword string, forbiddenIDs []string, gender int32, pagination pagination.Pagination) (int64, []*chat.Attribute, error) {
	filter := bson.M{}
	if gender == 0 {
		filter["gender"] = bson.M{
			"$in": []int32{0, 1, 2},
		}
	} else {
		filter["gender"] = gender
	}
	if len(forbiddenIDs) > 0 {
		filter["user_id"] = bson.M{
			"$nin": forbiddenIDs,
		}
	}
	if keyword != "" {
		filter["$or"] = []bson.M{
			{"user_id": bson.M{"$regex": keyword, "$options": "i"}},
			{"account": bson.M{"$regex": keyword, "$options": "i"}},
			{"nickname": bson.M{"$regex": keyword, "$options": "i"}},
			{"phone_number": bson.M{"$regex": keyword, "$options": "i"}},
			{"email": bson.M{"$regex": keyword, "$options": "i"}},
		}
	}
	return mongoutil.FindPage[*chat.Attribute](ctx, o.coll, filter, pagination)
}

func (o *Attribute) SearchUser(ctx context.Context, keyword string, userIDs []string, genders []int32, pagination pagination.Pagination) (int64, []*chat.Attribute, error) {
	filter := bson.M{}
	if len(genders) > 0 {
		filter["gender"] = bson.M{
			"$in": genders,
		}
	}
	if len(userIDs) > 0 {
		filter["user_id"] = bson.M{
			"$in": userIDs,
		}
	}
	if keyword != "" {
		filter["$or"] = []bson.M{
			{"user_id": bson.M{"$regex": keyword, "$options": "i"}},
			{"account": bson.M{"$regex": keyword, "$options": "i"}},
			{"nickname": bson.M{"$regex": keyword, "$options": "i"}},
			{"phone_number": bson.M{"$regex": keyword, "$options": "i"}},
			{"email": bson.M{"$regex": keyword, "$options": "i"}},
		}
	}
	return mongoutil.FindPage[*chat.Attribute](ctx, o.coll, filter, pagination)
}

func (o *Attribute) Delete(ctx context.Context, userIDs []string) error {
	if len(userIDs) == 0 {
		return nil
	}
	return mongoutil.DeleteMany(ctx, o.coll, bson.M{"user_id": bson.M{"$in": userIDs}})
}
