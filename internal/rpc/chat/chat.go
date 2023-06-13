package chat

import (
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/chat/pkg/common/db/database"
	chat2 "github.com/OpenIMSDK/chat/pkg/common/db/table/chat"
	"github.com/OpenIMSDK/chat/pkg/common/dbconn"
	"github.com/OpenIMSDK/chat/pkg/proto/chat"
	chatClient "github.com/OpenIMSDK/chat/pkg/rpclient/chat"
	"github.com/OpenIMSDK/chat/pkg/rpclient/openim"
	"github.com/OpenIMSDK/chat/pkg/sms"
	"google.golang.org/grpc"
)

func Start(zk discoveryregistry.SvcDiscoveryRegistry, server *grpc.Server) error {
	db, err := dbconn.NewGormDB()
	if err != nil {
		return err
	}
	tables := []any{
		chat2.Account{},
		chat2.Register{},
		chat2.Attribute{},
		chat2.VerifyCode{},
		chat2.UserLoginRecord{},
	}
	if err := db.AutoMigrate(tables...); err != nil {
		return err
	}
	s, err := sms.New()
	if err != nil {
		return err
	}
	chat.RegisterChatServer(server, &chatSvr{
		Database: database.NewChatDatabase(db),
		Admin:    chatClient.NewAdmin(zk),
		OpenIM:   openim.NewOpenIM(zk),
		SMS:      s,
	})
	return nil
}

type chatSvr struct {
	Database database.ChatDatabaseInterface
	Admin    *chatClient.Admin
	OpenIM   *openim.OpenIM
	SMS      sms.SMS
}
