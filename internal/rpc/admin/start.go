package admin

import (
	"context"
	"github.com/openimsdk/chat/pkg/common/config"
	"github.com/openimsdk/chat/pkg/common/db/database"
	"github.com/openimsdk/chat/pkg/proto/admin"
	"github.com/openimsdk/chat/pkg/rpclient/chat"
	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/db/redisutil"
	"github.com/openimsdk/tools/discovery"
	"github.com/openimsdk/tools/errs"
	"google.golang.org/grpc"
)

type Config struct {
	RpcConfig       config.Admin
	RedisConfig     config.Redis
	MongodbConfig   config.Mongo
	ZookeeperConfig config.ZooKeeper
	Share           config.Share
}

func Start(ctx context.Context, config *Config, client discovery.SvcDiscoveryRegistry, server *grpc.Server) error {
	return nil
}

func start(discov discoveryregistry.SvcDiscoveryRegistry, server *grpc.Server) error {
	mgocli, err := mongoutil.NewMongoDB(ctx, config.MongodbConfig.Build())
	if err != nil {
		return err
	}
	rdb, err := redisutil.NewRedisClient(ctx, config.RedisConfig.Build())
	if err != nil {
		return err
	}
	adminDatabase, err := database.NewAdminDatabase(mgocli, rdb)
	if err != nil {
		return err
	}

	if err := adminDatabase.InitAdmin(context.Background()); err != nil {
		return err
	}
	if err := discov.CreateRpcRootNodes([]string{config.Config.RpcRegisterName.OpenImAdminName, config.Config.RpcRegisterName.OpenImChatName}); err != nil {
		return errs.Wrap(err, "CreateRpcRootNodes error")
	}

	admin.RegisterAdminServer(server, &adminServer{Database: adminDatabase,
		Chat: chat.NewChatClient(discov),
	})
	return nil
}

type adminServer struct {
	Database database.AdminDatabaseInterface
	Chat     *chat.ChatClient
}
