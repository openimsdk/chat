// Copyright © 2023 OpenIM. All rights reserved.
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

package cmd

import (
	"context"
	"github.com/openimsdk/chat/internal/rpc/chat"
	"github.com/openimsdk/chat/pkg/common/config"
	"github.com/openimsdk/chat/pkg/common/startrpc"
	"github.com/openimsdk/tools/system/program"
	"github.com/spf13/cobra"
)

type ChatRpcCmd struct {
	*RootCmd
	ctx        context.Context
	configMap  map[string]any
	chatConfig chat.Config
}

func NewChatRpcCmd() *ChatRpcCmd {
	var chatConfig chat.Config
	ret := &ChatRpcCmd{chatConfig: chatConfig}
	ret.configMap = map[string]any{
		OpenIMRPCChatCfgFileName: &chatConfig.RpcConfig,
		RedisConfigFileName:      &chatConfig.RedisConfig,
		ZookeeperConfigFileName:  &chatConfig.ZookeeperConfig,
		MongodbConfigFileName:    &chatConfig.MongodbConfig,
		ShareFileName:            &chatConfig.Share,
		NotificationFileName:     &chatConfig.NotificationConfig,
		WebhooksConfigFileName:   &chatConfig.WebhooksConfig,
		LocalCacheConfigFileName: &chatConfig.LocalCacheConfig,
	}
	ret.RootCmd = NewRootCmd(program.GetProcessName(), WithConfigMap(ret.configMap))
	ret.ctx = context.WithValue(context.Background(), "version", config.Version)
	ret.Command.PreRunE = func(cmd *cobra.Command, args []string) error {
		return ret.preRunE()
	}
	return ret
}

func (a *ChatRpcCmd) Exec() error {
	return a.Execute()
}

func (a *ChatRpcCmd) preRunE() error {
	return startrpc.Start(a.ctx, &a.chatConfig.ZookeeperConfig, a.chatConfig.RpcConfig.RPC.ListenIP,
		a.chatConfig.RpcConfig.RPC.RegisterIP, a.chatConfig.RpcConfig.RPC.Ports,
		a.Index(), a.chatConfig.Share.RpcRegisterName.Auth, &a.chatConfig.Share, &a.chatConfig, chat.Start)
}
