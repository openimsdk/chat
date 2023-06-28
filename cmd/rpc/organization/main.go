package main

import (
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/startrpc"
	"github.com/OpenIMSDK/chat/internal/rpc/organization"
	"github.com/OpenIMSDK/chat/pkg/common/config"
)

func main() {
	if err := config.InitConfig(); err != nil {
		panic(err)
	}
	if err := log.InitFromConfig("organization.log", "organization-rpc", *config.Config.Log.RemainLogLevel, *config.Config.Log.IsStdout, *config.Config.Log.IsJson, *config.Config.Log.StorageLocation, *config.Config.Log.RemainRotationCount); err != nil {
		panic(err)
	}
	err := startrpc.Start(config.Config.RpcPort.OpenImOrganizationPort[0], config.Config.RpcRegisterName.OpenImOrganizationName, 0, organization.Start)
	if err != nil {
		panic(err)
	}
}
