// Copyright © 2023 OpenIM open source community. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package chat

import (
	"github.com/OpenIMSDK/chat/pkg/common/apicall"
	"github.com/OpenIMSDK/tools/discoveryregistry"
	"github.com/OpenIMSDK/tools/errs"
	"google.golang.org/grpc"

	"github.com/OpenIMSDK/chat/pkg/common/config"
	"github.com/OpenIMSDK/chat/pkg/common/db/database"
	chat2 "github.com/OpenIMSDK/chat/pkg/common/db/table/chat"
	"github.com/OpenIMSDK/chat/pkg/common/dbconn"
	"github.com/OpenIMSDK/chat/pkg/email"
	"github.com/OpenIMSDK/chat/pkg/proto/chat"
	chatClient "github.com/OpenIMSDK/chat/pkg/rpclient/chat"
	"github.com/OpenIMSDK/chat/pkg/sms"
)

func Start(discov discoveryregistry.SvcDiscoveryRegistry, server *grpc.Server) error {
	db, err := dbconn.NewGormDB()
	if err != nil {
		return errs.Wrap(err)
	}
	tables := []any{
		chat2.Account{},
		chat2.Register{},
		chat2.Attribute{},
		chat2.VerifyCode{},
		chat2.UserLoginRecord{},
		chat2.Log{},
	}
	if err := db.AutoMigrate(tables...); err != nil {
		return errs.Wrap(err)
	}
	s, err := sms.New()
	if err != nil {
		return errs.Wrap(err)
	}
	email := email.NewMail()
	if err := discov.CreateRpcRootNodes([]string{config.Config.RpcRegisterName.OpenImAdminName, config.Config.RpcRegisterName.OpenImChatName}); err != nil {
		panic(errs.Wrap(err, "CreateRpcRootNodes error"))
	}
	chat.RegisterChatServer(server, &chatSvr{
		Database:    database.NewChatDatabase(db),
		Admin:       chatClient.NewAdminClient(discov),
		SMS:         s,
		Mail:        email,
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
