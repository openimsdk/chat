package bot

import (
	"context"
	"time"

	"github.com/openimsdk/chat/pkg/common/config"
	"github.com/openimsdk/chat/pkg/common/db/database"
	"github.com/openimsdk/chat/pkg/common/imapi"
	"github.com/openimsdk/chat/pkg/protocol/bot"
	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/db/redisutil"
	"github.com/openimsdk/tools/discovery"
	"github.com/openimsdk/tools/utils/httputil"
	"google.golang.org/grpc"
)

type Config struct {
	RpcConfig     config.Bot
	RedisConfig   config.Redis
	MongodbConfig config.Mongo
	Discovery     config.Discovery
	Share         config.Share
}

func Start(ctx context.Context, config *Config, client discovery.SvcDiscoveryRegistry, server *grpc.Server) error {
	mgocli, err := mongoutil.NewMongoDB(ctx, config.MongodbConfig.Build())
	if err != nil {
		return err
	}
	var srv botSvr
	rdb, err := redisutil.NewRedisClient(ctx, config.RedisConfig.Build())
	if err != nil {
		return err
	}

	srv.database, err = database.NewBotDatabase(mgocli)
	if err != nil {
		return err
	}
	httpCfg := httputil.NewClientConfig()
	httpCfg.Timeout = time.Duration(config.RpcConfig.Timeout) * time.Second
	srv.timeout = config.RpcConfig.Timeout
	srv.httpClient = httputil.NewHTTPClient(httpCfg)
	im := imapi.New(config.Share.OpenIM.ApiURL, config.Share.OpenIM.Secret, config.Share.OpenIM.AdminUserID, rdb, config.Share.OpenIM.TokenRefreshInterval)
	srv.imCaller = im
	bot.RegisterBotServer(server, &srv)
	return nil
}

type botSvr struct {
	bot.UnimplementedBotServer
	database   database.BotDatabase
	httpClient *httputil.HTTPClient
	timeout    int
	imCaller   imapi.CallerInterface
	//Admin           *chatClient.AdminClient
}
