package chat

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/openimsdk/chat/pkg/common/db/table/chat"
	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/db/pagination"
	"github.com/openimsdk/tools/errs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewCredential(db *mongo.Database) (chat.CredentialInterface, error) {
	coll := db.Collection("credential")
	_, err := coll.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "user_id", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "account", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
	})
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return &Credential{coll: coll}, nil
}

type Credential struct {
	coll *mongo.Collection
}

func (o *Credential) Create(ctx context.Context, attribute ...*chat.Credential) error {
	return mongoutil.InsertMany(ctx, o.coll, attribute)
}

func (o *Credential) Update(ctx context.Context, userID string, data map[string]any) error {
	if len(data) == 0 {
		return nil
	}
	return mongoutil.UpdateOne(ctx, o.coll, bson.M{"user_id": userID}, bson.M{"$set": data}, false)
}

func (o *Credential) Find(ctx context.Context, userIds []string) ([]*chat.Credential, error) {
	return mongoutil.Find[*chat.Credential](ctx, o.coll, bson.M{"user_id": bson.M{"$in": userIds}})
}

func (o *Credential) FindAccount(ctx context.Context, accounts []string) ([]*chat.Credential, error) {
	return mongoutil.Find[*chat.Credential](ctx, o.coll, bson.M{"account": bson.M{"$in": accounts}})
}

func (o *Credential) Search(ctx context.Context, keyword string, pagination pagination.Pagination) (int64, []*chat.Credential, error) {
	return o.SearchUser(ctx, keyword, nil, pagination)
}

func (o *Credential) TakeAccount(ctx context.Context, account string) (*chat.Credential, error) {
	return mongoutil.FindOne[*chat.Credential](ctx, o.coll, bson.M{"account": account})
}

func (o *Credential) Take(ctx context.Context, userID string) (*chat.Credential, error) {
	return mongoutil.FindOne[*chat.Credential](ctx, o.coll, bson.M{"user_id": userID})
}

func (o *Credential) SearchNormalUser(ctx context.Context, keyword string, forbiddenIDs []string, pagination pagination.Pagination) (int64, []*chat.Credential, error) {
	filter := bson.M{}

	if len(forbiddenIDs) > 0 {
		filter["user_id"] = bson.M{
			"$nin": forbiddenIDs,
		}
	}
	if keyword != "" {
		filter["$or"] = []bson.M{
			{"user_id": bson.M{"$regex": keyword, "$options": "i"}},
			{"account": bson.M{"$regex": keyword, "$options": "i"}},
		}
	}
	return mongoutil.FindPage[*chat.Credential](ctx, o.coll, filter, pagination)
}

func (o *Credential) SearchUser(ctx context.Context, keyword string, userIDs []string, pagination pagination.Pagination) (int64, []*chat.Credential, error) {
	filter := bson.M{}

	if len(userIDs) > 0 {
		filter["user_id"] = bson.M{
			"$in": userIDs,
		}
	}
	if keyword != "" {
		filter["$or"] = []bson.M{
			{"user_id": bson.M{"$regex": keyword, "$options": "i"}},
			{"account": bson.M{"$regex": keyword, "$options": "i"}},
		}
	}
	return mongoutil.FindPage[*chat.Credential](ctx, o.coll, filter, pagination)
}

func (o *Credential) Delete(ctx context.Context, userIDs []string) error {
	if len(userIDs) == 0 {
		return nil
	}
	return mongoutil.DeleteMany(ctx, o.coll, bson.M{"user_id": bson.M{"$in": userIDs}})
}
