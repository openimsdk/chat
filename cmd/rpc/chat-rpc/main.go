package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/OpenIMSDK/chat/pkg/common/chatrpcstart"
	"github.com/OpenIMSDK/tools/log"

	"github.com/OpenIMSDK/chat/internal/rpc/chat"
	"github.com/OpenIMSDK/chat/pkg/common/config"
	"github.com/OpenIMSDK/chat/pkg/common/version"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	configFile, rpcPort, hide, showVersion, err := config.FlagParse()
	if err != nil {
		panic(err)
	}

	// Check if the version flag was set
	if showVersion {
		ver := version.Get()
		fmt.Println("Version:", ver.GitVersion)
		fmt.Println("Git Commit:", ver.GitCommit)
		fmt.Println("Build Date:", ver.BuildDate)
		fmt.Println("Go Version:", ver.GoVersion)
		fmt.Println("Compiler:", ver.Compiler)
		fmt.Println("Platform:", ver.Platform)
		return
	}
	if err := config.InitConfig(configFile, hide); err != nil {
		panic(err)
	}
	if config.Config.Envs.Discovery == "k8s" {
		rpcPort = 80
	}
	if err := log.InitFromConfig("chat.log", "chat-rpc", *config.Config.Log.RemainLogLevel, *config.Config.Log.IsStdout, *config.Config.Log.IsJson, *config.Config.Log.StorageLocation, *config.Config.Log.RemainRotationCount, *config.Config.Log.RotationTime); err != nil {
		panic(err)
	}
	err = chatrpcstart.Start(rpcPort, config.Config.RpcRegisterName.OpenImChatName, 0, chat.Start)
	if err != nil {
		panic(err)
	}
}
