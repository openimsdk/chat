package bot

import (
	"context"

	"github.com/openimsdk/chat/pkg/common/db/table/bot"
	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/errs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewConversationRespID(db *mongo.Database) (bot.ConversationRespIDInterface, error) {
	coll := db.Collection("conversation_resp_id")
	_, err := coll.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{
			{Key: "conversation_id", Value: 1},
			{Key: "agent_id", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return &ConversationRespID{coll: coll}, nil
}

type ConversationRespID struct {
	coll *mongo.Collection
}

func (o *ConversationRespID) Create(ctx context.Context, elems ...*bot.ConversationRespID) error {
	return mongoutil.InsertMany(ctx, o.coll, elems)
}

func (o *ConversationRespID) Take(ctx context.Context, convID, agentID string) (*bot.ConversationRespID, error) {
	return mongoutil.FindOne[*bot.ConversationRespID](ctx, o.coll, bson.M{"conversation_id": convID, "agent_id": agentID})
}

func (o *ConversationRespID) Update(ctx context.Context, convID, agentID string, data map[string]any) error {
	if len(data) == 0 {
		return nil
	}
	return mongoutil.UpdateOne(ctx, o.coll, bson.M{"conversation_id": convID, "agent_id": agentID}, bson.M{"$set": data}, false, options.Update().SetUpsert(true))
}

func (o *ConversationRespID) Delete(ctx context.Context, convID, agentID string) error {
	return mongoutil.DeleteMany(ctx, o.coll, bson.M{"conversation_id": convID, "agent_id": agentID})
}
