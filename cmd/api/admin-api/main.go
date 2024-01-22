// Copyright © 2023 OpenIM open source community. All rights reserved.
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

package main

import (
	"fmt"
	"github.com/OpenIMSDK/tools/errs"
	"math/rand"
	"net"
	"strconv"
	"time"

	"github.com/OpenIMSDK/chat/pkg/discovery_register"

	"github.com/OpenIMSDK/chat/tools/component"
	"github.com/OpenIMSDK/tools/discoveryregistry"

	mw2 "github.com/OpenIMSDK/chat/pkg/common/mw"
	"github.com/OpenIMSDK/chat/pkg/common/version"

	"github.com/OpenIMSDK/chat/internal/api"
	"github.com/OpenIMSDK/chat/pkg/common/config"
	"github.com/OpenIMSDK/tools/log"
	"github.com/OpenIMSDK/tools/mw"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/gin-gonic/gin"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	configFile, ginPort, hide, showVersion, err := config.FlagParse()
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

	if err := config.InitConfig(configFile); err != nil {
		panic(err)
	}
	err = component.ComponentCheck(configFile, hide)
	if err != nil {
		panic(err)
	}
	if err := log.InitFromConfig("chat.log", "admin-api", *config.Config.Log.RemainLogLevel, *config.Config.Log.IsStdout, *config.Config.Log.IsJson, *config.Config.Log.StorageLocation, *config.Config.Log.RemainRotationCount, *config.Config.Log.RotationTime); err != nil {
		panic(fmt.Errorf("InitFromConfig failed:%w", err))
	}
	if config.Config.Envs.Discovery == "k8s" {
		ginPort = 80
	}
	var zk discoveryregistry.SvcDiscoveryRegistry
	zk, err = discovery_register.NewDiscoveryRegister(config.Config.Envs.Discovery)
	//zk, err = openKeeper.NewClient(config.Config.Zookeeper.ZkAddr, config.Config.Zookeeper.Schema,
	//		openKeeper.WithFreq(time.Hour), openKeeper.WithUserNameAndPassword(config.Config.Zookeeper.Username, config.Config.Zookeeper.Password), openKeeper.WithRoundRobin(), openKeeper.WithTimeout(10), openKeeper.WithLogger(log.NewZkLogger()))
	if err != nil {
		panic(err)
	}

	if err := zk.CreateRpcRootNodes([]string{config.Config.RpcRegisterName.OpenImAdminName, config.Config.RpcRegisterName.OpenImChatName}); err != nil {
		panic(errs.Wrap(err, "CreateRpcRootNodes error"))
	}
	zk.AddOption(mw.GrpcClient(), grpc.WithTransportCredentials(insecure.NewCredentials())) // 默认RPC中间件
	engine := gin.Default()
	engine.Use(mw.CorsHandler(), mw.GinParseOperationID(), mw2.GinLog())
	api.NewAdminRoute(engine, zk)

	address := net.JoinHostPort(config.Config.AdminApi.ListenIP, strconv.Itoa(ginPort))
	if err := engine.Run(address); err != nil {
		panic(err)
	}
}
