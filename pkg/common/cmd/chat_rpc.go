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
