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
	"fmt"
	"math/rand"
	"time"

	"github.com/OpenIMSDK/chat/pkg/common/chatrpcstart"
	"github.com/OpenIMSDK/chat/pkg/common/version"
	"github.com/OpenIMSDK/chat/tools/component"
	"github.com/OpenIMSDK/tools/log"

	"github.com/OpenIMSDK/chat/internal/rpc/admin"
	"github.com/OpenIMSDK/chat/pkg/common/config"
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

	flag.Parse()

	if err := config.InitConfig(configFile); err != nil {
		panic(err)
	}
	err = component.ComponentCheck(configFile, hide)
	if err != nil {
		panic(err)
	}
	if config.Config.Envs.Discovery == "k8s" {
		rpcPort = 80
	}
	if err := log.InitFromConfig("chat.log", "admin-rpc", *config.Config.Log.RemainLogLevel, *config.Config.Log.IsStdout, *config.Config.Log.IsJson, *config.Config.Log.StorageLocation, *config.Config.Log.RemainRotationCount, *config.Config.Log.RotationTime); err != nil {
		panic(fmt.Errorf("InitFromConfig failed:%w", err))
	}
	err = chatrpcstart.Start(rpcPort, config.Config.RpcRegisterName.OpenImAdminName, 0, admin.Start)
	if err != nil {
		panic(err)
	}
}
