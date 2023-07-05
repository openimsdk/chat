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
	"flag"
	"math/rand"
	"net"
	"strconv"
	"time"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/mw"
	openKeeper "github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry/zookeeper"
	"github.com/OpenIMSDK/chat/internal/api"
	"github.com/OpenIMSDK/chat/pkg/common/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/gin-gonic/gin"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	err := config.InitConfig()
	if err != nil {
		panic(err)
	}
	if err := log.InitFromConfig("chat.log", "chat-api", *config.Config.Log.RemainLogLevel, *config.Config.Log.IsStdout, *config.Config.Log.IsJson, *config.Config.Log.StorageLocation, *config.Config.Log.RemainRotationCount); err != nil {
		panic(err)
	}
	zk, err := openKeeper.NewClient(config.Config.Zookeeper.ZkAddr, config.Config.Zookeeper.Schema,
		openKeeper.WithFreq(time.Hour), openKeeper.WithUserNameAndPassword(config.Config.Zookeeper.UserName,
			config.Config.Zookeeper.Password), openKeeper.WithRoundRobin(), openKeeper.WithTimeout(10), openKeeper.WithLogger(log.NewZkLogger()))
	if err != nil {
		panic(err)
	}
	if err := zk.CreateRpcRootNodes([]string{config.Config.RpcRegisterName.OpenImAdminName, config.Config.RpcRegisterName.OpenImChatName}); err != nil {
		panic(err)
	}
	zk.AddOption(mw.GrpcClient(), grpc.WithTransportCredentials(insecure.NewCredentials())) // 默认RPC中间件
	engine := gin.Default()
	engine.Use(mw.CorsHandler(), mw.GinParseOperationID())
	api.NewChatRoute(engine, zk)
	defaultPorts := config.Config.ChatApi.GinPort
	ginPort := flag.Int("port", defaultPorts[0], "get ginServerPort from cmd")
	flag.Parse()
	address := net.JoinHostPort(config.Config.ChatApi.ListenIP, strconv.Itoa(*ginPort))
	if err := engine.Run(address); err != nil {
		panic(err)
	}
}
