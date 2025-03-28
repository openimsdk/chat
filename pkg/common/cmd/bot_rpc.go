package cmd

import (
	"context"

	"github.com/openimsdk/chat/internal/rpc/bot"
	"github.com/openimsdk/chat/pkg/common/config"
	"github.com/openimsdk/chat/pkg/common/startrpc"
	"github.com/openimsdk/tools/system/program"
	"github.com/spf13/cobra"
)

type BotRpcCmd struct {
	*RootCmd
	ctx       context.Context
	configMap map[string]any
	botConfig bot.Config
}

func NewBotRpcCmd() *BotRpcCmd {
	var ret BotRpcCmd
	ret.configMap = map[string]any{
		config.ChatRPCBotCfgFileName:   &ret.botConfig.RpcConfig,
		config.RedisConfigFileName:     &ret.botConfig.RedisConfig,
		config.DiscoveryConfigFileName: &ret.botConfig.Discovery,
		config.MongodbConfigFileName:   &ret.botConfig.MongodbConfig,
		config.ShareFileName:           &ret.botConfig.Share,
	}
	ret.RootCmd = NewRootCmd(program.GetProcessName(), WithConfigMap(ret.configMap))
	ret.ctx = context.WithValue(context.Background(), "version", config.Version)
	ret.Command.RunE = func(cmd *cobra.Command, args []string) error {
		return ret.runE()
	}
	return &ret
}

func (a *BotRpcCmd) Exec() error {
	return a.Execute()
}

func (a *BotRpcCmd) runE() error {
	return startrpc.Start(a.ctx, &a.botConfig.Discovery, a.botConfig.RpcConfig.RPC.ListenIP,
		a.botConfig.RpcConfig.RPC.RegisterIP, a.botConfig.RpcConfig.RPC.Ports,
		a.Index(), a.botConfig.Discovery.RpcService.Bot, &a.botConfig.Share, &a.botConfig,
		[]string{
			config.ChatRPCBotCfgFileName,
			config.RedisConfigFileName,
			config.DiscoveryConfigFileName,
			config.MongodbConfigFileName,
			config.ShareFileName,
			config.LogConfigFileName,
		}, nil,
		bot.Start)
}
