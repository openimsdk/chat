package cmd

import (
	"context"

	"github.com/openimsdk/chat/internal/api/bot"
	"github.com/openimsdk/chat/pkg/common/config"
	"github.com/openimsdk/tools/system/program"
	"github.com/spf13/cobra"
)

type BotApiCmd struct {
	*RootCmd
	ctx       context.Context
	configMap map[string]any
	apiConfig bot.Config
}

func NewBotApiCmd() *BotApiCmd {
	ret := BotApiCmd{apiConfig: bot.Config{}}
	ret.configMap = map[string]any{
		config.DiscoveryConfigFileName: &ret.apiConfig.Discovery,
		config.ChatAPIBotCfgFileName:   &ret.apiConfig.ApiConfig,
		config.ShareFileName:           &ret.apiConfig.Share,
		config.RedisConfigFileName:     &ret.apiConfig.Redis,
	}
	ret.RootCmd = NewRootCmd(program.GetProcessName(), WithConfigMap(ret.configMap))
	ret.ctx = context.WithValue(context.Background(), "version", config.Version)
	ret.Command.RunE = func(cmd *cobra.Command, args []string) error {
		return ret.runE()
	}
	return &ret
}

func (a *BotApiCmd) Exec() error {
	return a.Execute()
}

func (a *BotApiCmd) runE() error {
	return bot.Start(a.ctx, a.Index(), &a.apiConfig)
}
