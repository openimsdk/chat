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
	return mongoutil.InsertMany(ctx, o.coll, ms)
}

func (o *ForbiddenAccount) Take(ctx context.Context, userID string) (*admin.ForbiddenAccount, error) {
	return mongoutil.FindOne[*admin.ForbiddenAccount](ctx, o.coll, bson.M{"user_id": userID})
}

func (o *ForbiddenAccount) Delete(ctx context.Context, userIDs []string) error {
	if len(userIDs) == 0 {
		return nil
	}
	return mongoutil.DeleteMany(ctx, o.coll, bson.M{"user_id": bson.M{"$in": userIDs}})
}

func (o *ForbiddenAccount) Find(ctx context.Context, userIDs []string) ([]*admin.ForbiddenAccount, error) {
	return mongoutil.Find[*admin.ForbiddenAccount](ctx, o.coll, bson.M{"user_id": bson.M{"$in": userIDs}})
}

func (o *ForbiddenAccount) Search(ctx context.Context, keyword string, pagination pagination.Pagination) (int64, []*admin.ForbiddenAccount, error) {
	filter := bson.M{}

	if keyword != "" {
		filter = bson.M{
			"$or": []bson.M{
				{"user_id": bson.M{"$regex": keyword, "$options": "i"}},
				{"reason": bson.M{"$regex": keyword, "$options": "i"}},
				{"operator_user_id": bson.M{"$regex": keyword, "$options": "i"}},
			},
		}
	}
	return mongoutil.FindPage[*admin.ForbiddenAccount](ctx, o.coll, filter, pagination)
}

func (o *ForbiddenAccount) FindAllIDs(ctx context.Context) ([]string, error) {
	return mongoutil.Find[string](ctx, o.coll, bson.M{}, options.Find().SetProjection(bson.M{"_id": 0, "user_id": 1}))
}
