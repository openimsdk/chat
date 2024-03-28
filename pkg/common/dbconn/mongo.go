package dbconn

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/OpenIMSDK/chat/pkg/common/config"
	"github.com/OpenIMSDK/tools/errs"
	"github.com/OpenIMSDK/tools/mw/specialerror"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	maxRetry         = 10 // number of retries
	mongoConnTimeout = 10 * time.Second
)

// NewMongo Initialize MongoDB connection.
func NewMongo() (*mongo.Database, error) {
	specialerror.AddReplace(mongo.ErrNoDocuments, errs.ErrRecordNotFound)
	uri := buildMongoURI()
	database := buildMongoDatabase()
	var mongoClient *mongo.Client
	var err error

	// Retry connecting to MongoDB
	for i := 0; i <= maxRetry; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), mongoConnTimeout)
		mongoClient, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
		cancel()
		if err == nil {
			return mongoClient.Database(database), nil
		}
		if shouldRetry(err) {
			fmt.Printf("Failed to connect to MongoDB, retrying: %s\n", err)
			time.Sleep(time.Second) // exponential backoff could be implemented here
			continue
		}
		return nil, errs.Wrap(err)
	}
	return nil, errs.Wrap(err)
}

func buildMongoURI() string {
	uri := os.Getenv("MONGO_URI")
	if uri != "" {
		return uri
	}

	if config.Config.Mongo.Uri != "" {
		return config.Config.Mongo.Uri
	}

	username := os.Getenv("MONGO_OPENIM_USERNAME")
	password := os.Getenv("MONGO_OPENIM_PASSWORD")
	address := os.Getenv("MONGO_ADDRESS")
	port := os.Getenv("MONGO_PORT")
	database := os.Getenv("MONGO_DATABASE")
	maxPoolSize := os.Getenv("MONGO_MAX_POOL_SIZE")

	if username == "" {
		username = config.Config.Mongo.Username
	}
	if password == "" {
		password = config.Config.Mongo.Password
	}
	if address == "" {
		address = strings.Join(config.Config.Mongo.Address, ",")
	} else if port != "" {
		address = fmt.Sprintf("%s:%s", address, port)
	}
	if database == "" {
		database = config.Config.Mongo.Database
	}
	if maxPoolSize == "" {
		maxPoolSize = fmt.Sprint(config.Config.Mongo.MaxPoolSize)
	}

	if username != "" && password != "" {
		return fmt.Sprintf("mongodb://%s:%s@%s/%s?maxPoolSize=%s", username, password, address, database, maxPoolSize)
	}
	return fmt.Sprintf("mongodb://%s/%s?maxPoolSize=%s", address, database, maxPoolSize)
}

func buildMongoDatabase() string {
	return config.Config.Mongo.Database
}

func shouldRetry(err error) bool {
	if cmdErr, ok := err.(mongo.CommandError); ok {
		return cmdErr.Code != 13 && cmdErr.Code != 18
	}
	return true
}
