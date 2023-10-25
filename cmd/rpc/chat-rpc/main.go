package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"

	"github.com/OpenIMSDK/chat/pkg/common/chatrpcstart"
	"github.com/OpenIMSDK/chat/tools/component"
	"github.com/OpenIMSDK/tools/log"

	"github.com/OpenIMSDK/chat/internal/rpc/chat"
	"github.com/OpenIMSDK/chat/pkg/common/config"
	"github.com/OpenIMSDK/chat/pkg/common/version"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	var configFile string
	flag.StringVar(&configFile, "config_folder_path", "../../../../../config/config.yaml", "Config full path")

	var rpcPort int
	flag.IntVar(&rpcPort, "port", 30300, "get rpc ServerPort from cmd")

	var hide bool
	flag.BoolVar(&hide, "hide", true, "hide the ComponentCheck result")

	// Version flag
	var showVersion bool
	flag.BoolVar(&showVersion, "version", false, "show version and exit")

	flag.Parse()

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

	err := component.ComponentCheck(configFile, hide)
	if err != nil {
		return
	}
	if err := config.InitConfig(configFile); err != nil {
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
