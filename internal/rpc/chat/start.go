package chat

import (
	"context"
	"github.com/openimsdk/chat/pkg/proto/admin"
	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/discovery"
	"google.golang.org/grpc"

	"github.com/openimsdk/chat/pkg/common/config"
	"github.com/openimsdk/chat/pkg/common/db/database"
	"github.com/openimsdk/chat/pkg/email"
	chatClient "github.com/openimsdk/chat/pkg/rpclient/chat"
	"github.com/openimsdk/chat/pkg/sms"
)

type Config struct {
	RpcConfig       config.Chat
	RedisConfig     config.Redis
	MongodbConfig   config.Mongo
	ZookeeperConfig config.ZooKeeper
	Share           config.Share
}

func Start(ctx context.Context, config *Config, client discovery.SvcDiscoveryRegistry, server *grpc.Server) error {
	mgocli, err := mongoutil.NewMongoDB(ctx, config.MongodbConfig.Build())
	if err != nil {
		return err
	}
	var srv chatSvr
	switch config.RpcConfig.VerifyCode.Phone.Use {
	case "ali":
		ali := config.RpcConfig.VerifyCode.Phone.Ali
		srv.SMS, err = sms.NewAli(ali.Endpoint, ali.AccessKeyID, ali.AccessKeySecret, ali.SignName, ali.VerificationCodeTemplateCode)
		if err != nil {
			return err
		}
	}
	if mail := config.RpcConfig.VerifyCode.Mail; mail.Enable {
		srv.Mail = email.NewMail(mail.SMTPAddr, mail.SMTPPort, mail.SenderMail, mail.SenderAuthorizationCode, mail.Title)
	}
	srv.Database, err = database.NewChatDatabase(mgocli.GetDB())
	if err != nil {
		return err
	}
	conn, err := client.GetConn(ctx, config.Share.RpcRegisterName.Admin)
	if err != nil {
		return err
	}
	srv.Admin = chatClient.NewAdminClient(admin.NewAdminClient(conn))
	return nil
}

type chatSvr struct {
	Database database.ChatDatabaseInterface
	Admin    *chatClient.AdminClient
	SMS      sms.SMS
	Mail     email.Mail
}
