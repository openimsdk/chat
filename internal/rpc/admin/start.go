package admin

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"time"

	"github.com/openimsdk/chat/pkg/common/config"
	"github.com/openimsdk/chat/pkg/common/constant"
	"github.com/openimsdk/chat/pkg/common/db/database"
	"github.com/openimsdk/chat/pkg/common/db/dbutil"
	"github.com/openimsdk/chat/pkg/common/db/table/admin"
	"github.com/openimsdk/chat/pkg/common/tokenverify"
	adminpb "github.com/openimsdk/chat/pkg/protocol/admin"
	"github.com/openimsdk/chat/pkg/protocol/chat"
	chatClient "github.com/openimsdk/chat/pkg/rpclient/chat"
	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/db/redisutil"
	"github.com/openimsdk/tools/discovery"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/mw"
	"github.com/openimsdk/tools/utils/runtimeenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Config struct {
	RpcConfig     config.Admin
	RedisConfig   config.Redis
	MongodbConfig config.Mongo
	Discovery     config.Discovery
	Share         config.Share

	RuntimeEnv string
}

func Start(ctx context.Context, config *Config, client discovery.SvcDiscoveryRegistry, server *grpc.Server) error {
	config.RuntimeEnv = runtimeenv.PrintRuntimeEnvironment()

	if len(config.Share.ChatAdmin) == 0 {
		return errs.New("share chat admin not configured")
	}
	rand.Seed(time.Now().UnixNano())
	rdb, err := redisutil.NewRedisClient(ctx, config.RedisConfig.Build())
	if err != nil {
		return err
	}
	mgocli, err := mongoutil.NewMongoDB(ctx, config.MongodbConfig.Build())
	if err != nil {
		return err
	}
	var srv adminServer
	srv.Database, err = database.NewAdminDatabase(mgocli, rdb)
	if err != nil {
		return err
	}
	conn, err := client.GetConn(ctx, config.Discovery.RpcService.Chat, grpc.WithTransportCredentials(insecure.NewCredentials()), mw.GrpcClient())
	if err != nil {
		return err
	}
	srv.Chat = chatClient.NewChatClient(chat.NewChatClient(conn))
	srv.Token = &tokenverify.Token{
		Expires: time.Duration(config.RpcConfig.TokenPolicy.Expire) * time.Hour * 24,
		Secret:  config.RpcConfig.Secret,
	}
	if err := srv.initAdmin(ctx, config.Share.ChatAdmin, config.Share.OpenIM.AdminUserID); err != nil {
		return err
	}
	adminpb.RegisterAdminServer(server, &srv)
	return nil
}

type adminServer struct {
	adminpb.UnimplementedAdminServer
	Database database.AdminDatabaseInterface
	Chat     *chatClient.ChatClient
	Token    *tokenverify.Token
}

func (o *adminServer) initAdmin(ctx context.Context, admins []string, imUserID string) error {
	for _, account := range admins {
		if _, err := o.Database.GetAdmin(ctx, account); err == nil {
			continue
		} else if !dbutil.IsDBNotFound(err) {
			return err
		}
		sum := md5.Sum([]byte(account))
		a := admin.Admin{
			Account:    account,
			UserID:     imUserID,
			Password:   hex.EncodeToString(sum[:]),
			Level:      constant.DefaultAdminLevel,
			CreateTime: time.Now(),
		}
		if err := o.Database.AddAdminAccount(ctx, []*admin.Admin{&a}); err != nil {
			return err
		}
	}
	return nil
}
