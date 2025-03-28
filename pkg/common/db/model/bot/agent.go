package bot

import (
	"context"

	"github.com/openimsdk/chat/pkg/common/db/table/bot"
	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/db/pagination"
	"github.com/openimsdk/tools/errs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewAgent(db *mongo.Database) (bot.AgentInterface, error) {
	coll := db.Collection("agent")
	_, err := coll.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{
			{Key: "user_id", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return &Agent{coll: coll}, nil
}

type Agent struct {
	coll *mongo.Collection
}

func (o *Agent) Create(ctx context.Context, elems ...*bot.Agent) error {
	return mongoutil.InsertMany(ctx, o.coll, elems)
}

func (o *Agent) Take(ctx context.Context, userId string) (*bot.Agent, error) {
	return mongoutil.FindOne[*bot.Agent](ctx, o.coll, bson.M{"user_id": userId})
}

func (o *Agent) Find(ctx context.Context, userIDs []string) ([]*bot.Agent, error) {
	return mongoutil.Find[*bot.Agent](ctx, o.coll, bson.M{"user_id": bson.M{"$in": userIDs}})
}

func (o *Agent) Update(ctx context.Context, userID string, data map[string]any) error {
	if len(data) == 0 {
		return nil
	}
	return mongoutil.UpdateOne(ctx, o.coll, bson.M{"user_id": userID}, bson.M{"$set": data}, false)
}

func (o *Agent) Delete(ctx context.Context, userIDs []string) error {
	if len(userIDs) == 0 {
		return nil
	}
	return mongoutil.DeleteMany(ctx, o.coll, bson.M{"user_id": bson.M{"$in": userIDs}})
}

func (o *Agent) Page(ctx context.Context, userIDs []string, pagination pagination.Pagination) (int64, []*bot.Agent, error) {
	filter := bson.M{}
	if len(userIDs) > 0 {
		filter["user_id"] = bson.M{"$in": userIDs}
	}
	return mongoutil.FindPage[*bot.Agent](ctx, o.coll, filter, pagination)
}
