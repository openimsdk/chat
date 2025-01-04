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

	"github.com/openimsdk/chat/internal/rpc/admin"
	"github.com/openimsdk/chat/pkg/common/config"
	"github.com/openimsdk/chat/pkg/common/startrpc"
	"github.com/openimsdk/tools/system/program"
	"github.com/spf13/cobra"
)

type AdminRpcCmd struct {
	*RootCmd
	ctx         context.Context
	configMap   map[string]any
	adminConfig admin.Config
}

func NewAdminRpcCmd() *AdminRpcCmd {
	var ret AdminRpcCmd
	ret.configMap = map[string]any{
		config.ChatRPCAdminCfgFileName: &ret.adminConfig.RpcConfig,
		config.RedisConfigFileName:     &ret.adminConfig.RedisConfig,
		config.DiscoveryConfigFileName: &ret.adminConfig.Discovery,
		config.MongodbConfigFileName:   &ret.adminConfig.MongodbConfig,
		config.ShareFileName:           &ret.adminConfig.Share,
	}
	ret.RootCmd = NewRootCmd(program.GetProcessName(), WithConfigMap(ret.configMap))
	ret.ctx = context.WithValue(context.Background(), "version", config.Version)
	ret.Command.RunE = func(cmd *cobra.Command, args []string) error {
		return ret.runE()
	}
	return &ret
}

func (a *AdminRpcCmd) Exec() error {
	return a.Execute()
}

func (a *AdminRpcCmd) runE() error {
	return startrpc.Start(a.ctx, &a.adminConfig.Discovery, a.adminConfig.RpcConfig.RPC.ListenIP,
		a.adminConfig.RpcConfig.RPC.RegisterIP, a.adminConfig.RpcConfig.RPC.Ports,
		a.Index(), a.adminConfig.Discovery.RpcService.Admin, &a.adminConfig.Share, &a.adminConfig,
		[]string{
			config.ChatRPCAdminCfgFileName,
			config.RedisConfigFileName,
			config.DiscoveryConfigFileName,
			config.MongodbConfigFileName,
			config.ShareFileName,
			config.LogConfigFileName,
		}, nil,
		admin.Start)
}
