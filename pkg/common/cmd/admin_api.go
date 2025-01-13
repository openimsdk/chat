package cmd

import (
	"context"

	"github.com/openimsdk/chat/internal/api/admin"
	"github.com/openimsdk/chat/pkg/common/config"
	"github.com/openimsdk/tools/system/program"
	"github.com/spf13/cobra"
)

type AdminApiCmd struct {
	*RootCmd
	ctx       context.Context
	configMap map[string]any
	apiConfig admin.Config
}

func NewAdminApiCmd() *AdminApiCmd {
	ret := AdminApiCmd{apiConfig: admin.Config{
		AllConfig: &config.AllConfig{},
	}}
	ret.configMap = map[string]any{
		config.DiscoveryConfigFileName: &ret.apiConfig.Discovery,
		config.LogConfigFileName:       &ret.apiConfig.Log,
		config.MongodbConfigFileName:   &ret.apiConfig.Mongo,
		config.ChatAPIAdminCfgFileName: &ret.apiConfig.AdminAPI,
		config.ChatAPIChatCfgFileName:  &ret.apiConfig.ChatAPI,
		config.ChatRPCAdminCfgFileName: &ret.apiConfig.Admin,
		config.ChatRPCChatCfgFileName:  &ret.apiConfig.Chat,
		config.RedisConfigFileName:     &ret.apiConfig.Redis,
		config.ShareFileName:           &ret.apiConfig.Share,
	}
	ret.RootCmd = NewRootCmd(program.GetProcessName(), WithConfigMap(ret.configMap))
	ret.ctx = context.WithValue(context.Background(), "version", config.Version)
	ret.Command.RunE = func(cmd *cobra.Command, args []string) error {
		ret.apiConfig.ConfigPath = ret.configPath
		return ret.runE()
	}
	return &ret
}

func (a *AdminApiCmd) Exec() error {
	return a.Execute()
}

func (a *AdminApiCmd) runE() error {
	return admin.Start(a.ctx, a.Index(), &a.apiConfig)
}
