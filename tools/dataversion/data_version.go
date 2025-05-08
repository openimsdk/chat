package dataversion

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/openimsdk/tools/db/mongoutil"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	Collection = "data_version"
)

func CheckVersion(coll *mongo.Collection, key string, currentVersion int) (converted bool, err error) {
	type VersionTable struct {
		Key   string `bson:"key"`
		Value string `bson:"value"`
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	res, err := mongoutil.FindOne[VersionTable](ctx, coll, bson.M{"key": key})
	if err == nil {
		ver, err := strconv.Atoi(res.Value)
		if err != nil {
			return false, fmt.Errorf("version %s parse error %w", res.Value, err)
		}
		if ver >= currentVersion {
			return true, nil
		}
		return false, nil
	} else if errors.Is(err, mongo.ErrNoDocuments) {
		return false, nil
	} else {
		return false, err
	}
}

func SetVersion(coll *mongo.Collection, key string, version int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	option := options.Update().SetUpsert(true)
	filter := bson.M{"key": key}
	update := bson.M{"$set": bson.M{"key": key, "value": strconv.Itoa(version)}}
	return mongoutil.UpdateOne(ctx, coll, filter, update, false, option)
}
