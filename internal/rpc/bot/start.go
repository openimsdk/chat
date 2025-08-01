package bot

import (
	"context"
	"net/http"
	"time"

	"github.com/openimsdk/chat/pkg/common/config"
	"github.com/openimsdk/chat/pkg/common/db/database"
	"github.com/openimsdk/chat/pkg/common/imapi"
	"github.com/openimsdk/chat/pkg/protocol/bot"
	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/discovery"
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

	srv.database, err = database.NewBotDatabase(mgocli)
	if err != nil {
		return err
	}
	srv.timeout = config.RpcConfig.Timeout
	srv.httpClient = &http.Client{
		Timeout: time.Duration(config.RpcConfig.Timeout) * time.Second,
	}
	im := imapi.New(config.Share.OpenIM.ApiURL, config.Share.OpenIM.Secret, config.Share.OpenIM.AdminUserID)
	srv.imCaller = im
	bot.RegisterBotServer(server, &srv)
	return nil
}

type botSvr struct {
	bot.UnimplementedBotServer
	database   database.BotDatabase
	httpClient *http.Client
	timeout    int
	imCaller   imapi.CallerInterface
	//Admin           *chatClient.AdminClient
}
