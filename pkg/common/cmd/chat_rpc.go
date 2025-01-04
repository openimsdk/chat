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
	var ret ChatRpcCmd
	ret.configMap = map[string]any{
		config.ChatRPCChatCfgFileName:  &ret.chatConfig.RpcConfig,
		config.RedisConfigFileName:     &ret.chatConfig.RedisConfig,
		config.DiscoveryConfigFileName: &ret.chatConfig.Discovery,
		config.MongodbConfigFileName:   &ret.chatConfig.MongodbConfig,
		config.ShareFileName:           &ret.chatConfig.Share,
	}
	ret.RootCmd = NewRootCmd(program.GetProcessName(), WithConfigMap(ret.configMap))
	ret.ctx = context.WithValue(context.Background(), "version", config.Version)
	ret.Command.RunE = func(cmd *cobra.Command, args []string) error {
		return ret.runE()
	}
	return &ret
}

func (a *ChatRpcCmd) Exec() error {
	return a.Execute()
}

func (a *ChatRpcCmd) runE() error {
	return startrpc.Start(a.ctx, &a.chatConfig.Discovery, a.chatConfig.RpcConfig.RPC.ListenIP,
		a.chatConfig.RpcConfig.RPC.RegisterIP, a.chatConfig.RpcConfig.RPC.Ports,
		a.Index(), a.chatConfig.Discovery.RpcService.Chat, &a.chatConfig.Share, &a.chatConfig,
		[]string{
			config.ChatRPCChatCfgFileName,
			config.RedisConfigFileName,
			config.DiscoveryConfigFileName,
			config.MongodbConfigFileName,
			config.ShareFileName,
			config.LogConfigFileName,
		}, nil,
		chat.Start)
}
