package main

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"

	"github.com/openimsdk/chat/internal/rpc/chat"
	"github.com/openimsdk/chat/pkg/common/config"
	"github.com/openimsdk/chat/pkg/common/constant"
	table "github.com/openimsdk/chat/pkg/common/db/table/chat"
	"github.com/openimsdk/chat/tools/dataversion"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/system/program"
	"github.com/openimsdk/tools/utils/runtimeenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	credentialKey     = "credential"
	credentialVersion = 1

	attributeCollection  = "attribute"
	credentialCollection = "credential"
	pageNum              = 1000
)

func initConfig(configDir string) (*config.Mongo, error) {
	var (
		mongoConfig = &config.Mongo{}
	)

	runtimeEnv := runtimeenv.PrintRuntimeEnvironment()

	err := config.Load(configDir, config.MongodbConfigFileName, config.EnvPrefixMap[config.MongodbConfigFileName], runtimeEnv, mongoConfig)
	if err != nil {
		return nil, err
	}

	return mongoConfig, nil
}

func pageGetAttribute(ctx context.Context, coll *mongo.Collection, pagination *sdkws.RequestPagination) (int64, []*table.Attribute, error) {
	return mongoutil.FindPage[*table.Attribute](ctx, coll, bson.M{}, pagination)
}

func doAttributeToCredential() error {
	var index int
	var configDir string
	flag.IntVar(&index, "i", 0, "Index number")
	defaultConfigDir := filepath.Join("..", "..", "..", "..", "..", "config")
	flag.StringVar(&configDir, "c", defaultConfigDir, "Configuration dir")
	flag.Parse()

	fmt.Printf("Index: %d, Config Path: %s\n", index, configDir)

	mongoConfig, err := initConfig(configDir)
	if err != nil {
		return err
	}

	ctx := context.Background()

	mgocli, err := mongoutil.NewMongoDB(ctx, mongoConfig.Build())
	if err != nil {
		return err
	}

	versionColl := mgocli.GetDB().Collection(dataversion.Collection)
	converted, err := dataversion.CheckVersion(versionColl, credentialKey, credentialVersion)
	if err != nil {
		return err
	}
	if converted {
		fmt.Println("[credential] credential data has been converted")
		return nil
	}

	attrColl := mgocli.GetDB().Collection(attributeCollection)
	credColl := mgocli.GetDB().Collection(credentialCollection)

	pagination := &sdkws.RequestPagination{
		PageNumber: 1,
		ShowNumber: pageNum,
	}
	tx := mgocli.GetTx()
	if err = tx.Transaction(ctx, func(ctx context.Context) error {
		for {
			_, attrs, err := pageGetAttribute(ctx, attrColl, pagination)
			if err != nil {
				return err
			}
			credentials := make([]*table.Credential, 0, pageNum*3)
			for _, attr := range attrs {
				if attr.Email != "" {
					credentials = append(credentials, &table.Credential{
						UserID:      attr.UserID,
						Account:     attr.Email,
						Type:        constant.CredentialEmail,
						AllowChange: true,
					})
				}
				if attr.Account != "" {
					credentials = append(credentials, &table.Credential{
						UserID:      attr.UserID,
						Account:     attr.Account,
						Type:        constant.CredentialAccount,
						AllowChange: true,
					})
				}
				if attr.PhoneNumber != "" && attr.AreaCode != "" {
					credentials = append(credentials, &table.Credential{
						UserID:      attr.UserID,
						Account:     chat.BuildCredentialPhone(attr.AreaCode, attr.PhoneNumber),
						Type:        constant.CredentialPhone,
						AllowChange: true,
					})
				}

			}
			for _, credential := range credentials {
				err = mongoutil.UpdateOne(ctx, credColl, bson.M{
					"user_id": credential.UserID,
					"type":    credential.Type,
				}, bson.M{
					"$set": credential,
				}, false, options.Update().SetUpsert(true))
				if err != nil {
					return err
				}
			}

			pagination.PageNumber++
			if len(attrs) < pageNum {
				break
			}
		}
		return nil
	}); err != nil {
		return err
	}
	if err := dataversion.SetVersion(versionColl, credentialKey, credentialVersion); err != nil {
		return fmt.Errorf("set mongodb credential version %w", err)
	}
	fmt.Println("[credential] update old data to credential success")
	return nil
}

func main() {
	if err := doAttributeToCredential(); err != nil {
		program.ExitWithError(err)
	}
}
