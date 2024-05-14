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
	var ret AdminApiCmd
	ret.configMap = map[string]any{
		ShareFileName:           &ret.apiConfig.Share,
		ChatAPIAdminCfgFileName: &ret.apiConfig.ApiConfig,
		DiscoveryConfigFileName: &ret.apiConfig.Discovery,
	}
	ret.RootCmd = NewRootCmd(program.GetProcessName(), WithConfigMap(ret.configMap))
	ret.ctx = context.WithValue(context.Background(), "version", config.Version)
	ret.Command.RunE = func(cmd *cobra.Command, args []string) error {
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
