package admin

import (
	"context"

	"github.com/openimsdk/tools/db/mongoutil"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/openimsdk/chat/pkg/common/db/table/admin"
	"github.com/openimsdk/tools/errs"
)

func NewClientConfig(db *mongo.Database) (admin.ClientConfigInterface, error) {
	coll := db.Collection("client_config")
	_, err := coll.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{
			{Key: "key", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return &ClientConfig{
		coll: coll,
	}, nil
}

type ClientConfig struct {
	coll *mongo.Collection
}

func (o *ClientConfig) Set(ctx context.Context, config map[string]string) error {
	for key, value := range config {
		filter := bson.M{"key": key}
		update := bson.M{
			"value": value,
		}
		err := mongoutil.UpdateOne(ctx, o.coll, filter, bson.M{"$set": update}, false, options.Update().SetUpsert(true))
		if err != nil {
			return err
		}
	}
	return nil
}

func (o *ClientConfig) Del(ctx context.Context, keys []string) error {
	return mongoutil.DeleteMany(ctx, o.coll, bson.M{"key": bson.M{"$in": keys}})
}

func (o *ClientConfig) Get(ctx context.Context) (map[string]string, error) {
	cs, err := mongoutil.Find[*admin.ClientConfig](ctx, o.coll, bson.M{})
	if err != nil {
		return nil, err
	}
	cm := make(map[string]string)
	for _, config := range cs {
		cm[config.Key] = config.Value
	}
	return cm, nil
}
