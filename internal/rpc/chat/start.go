package chat

import (
	"github.com/openimsdk/chat/pkg/common/apicall"
	"github.com/openimsdk/tools/discoveryregistry"
	"github.com/openimsdk/tools/errs"
	"google.golang.org/grpc"

	"github.com/openimsdk/chat/pkg/common/config"
	"github.com/openimsdk/chat/pkg/common/db/database"
	"github.com/openimsdk/chat/pkg/common/dbconn"
	"github.com/openimsdk/chat/pkg/email"
	"github.com/openimsdk/chat/pkg/proto/chat"
	chatClient "github.com/openimsdk/chat/pkg/rpclient/chat"
	"github.com/openimsdk/chat/pkg/sms"
)

func Start(discov discoveryregistry.SvcDiscoveryRegistry, server *grpc.Server) error {
	mgodb, err := dbconn.NewMongo()
	if err != nil {
		return err
	}
	s, err := sms.New()
	if err != nil {
		return errs.Wrap(err)
	}
	db, err := database.NewChatDatabase(mgodb)
	if err != nil {
		return err
	}
	if err := discov.CreateRpcRootNodes([]string{config.Config.RpcRegisterName.OpenImAdminName, config.Config.RpcRegisterName.OpenImChatName}); err != nil {
		return err
	}
	chat.RegisterChatServer(server, &chatSvr{
		Database:    db,
		Admin:       chatClient.NewAdminClient(discov),
		SMS:         s,
		Mail:        email.NewMail(),
		imApiCaller: apicall.NewCallerInterface(),
	})
	return nil
}

type chatSvr struct {
	Database    database.ChatDatabaseInterface
	Admin       *chatClient.AdminClient
	SMS         sms.SMS
	Mail        email.Mail
	imApiCaller apicall.CallerInterface
}
