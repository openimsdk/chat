package admin

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"github.com/openimsdk/chat/pkg/common/config"
	"github.com/openimsdk/chat/pkg/common/constant"
	"github.com/openimsdk/chat/pkg/common/db/database"
	"github.com/openimsdk/chat/pkg/common/db/dbutil"
	"github.com/openimsdk/chat/pkg/common/db/table/admin"
	"github.com/openimsdk/chat/pkg/common/tokenverify"
	pbadmin "github.com/openimsdk/chat/pkg/protocol/admin"
	"github.com/openimsdk/chat/pkg/protocol/chat"
	chatClient "github.com/openimsdk/chat/pkg/rpclient/chat"
	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/db/redisutil"
	"github.com/openimsdk/tools/discovery"
	"github.com/openimsdk/tools/errs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"math/rand"
	"time"
)

type Config struct {
	RpcConfig       config.Admin
	RedisConfig     config.Redis
	MongodbConfig   config.Mongo
	ZookeeperConfig config.ZooKeeper
	Share           config.Share
}

func Start(ctx context.Context, config *Config, client discovery.SvcDiscoveryRegistry, server *grpc.Server) error {
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
	conn, err := client.GetConn(ctx, config.Share.RpcRegisterName.Chat, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	srv.Chat = chatClient.NewChatClient(chat.NewChatClient(conn))
	srv.Token = &tokenverify.Token{
		Expires: time.Duration(config.RpcConfig.TokenPolicy.Expire) * time.Hour * 24,
		Secret:  config.RpcConfig.Secret,
	}
	if err := srv.initAdmin(ctx, config.Share.ChatAdmin); err != nil {
		return err
	}
	pbadmin.RegisterAdminServer(server, &srv)
	return nil
}

type adminServer struct {
	Database database.AdminDatabaseInterface
	Chat     *chatClient.ChatClient
	Token    *tokenverify.Token
}

func (o *adminServer) initAdmin(ctx context.Context, users []config.AdminUser) error {
	for _, user := range users {
		if _, err := o.Database.GetAdmin(ctx, user.AdminID); err == nil {
			continue
		} else if !dbutil.IsDBNotFound(err) {
			return err
		}
		sum := md5.Sum([]byte(user.AdminID))
		a := admin.Admin{
			Account:    user.AdminID,
			UserID:     user.IMUserID,
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
