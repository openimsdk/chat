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
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/OpenIMSDK/chat/pkg/common/chatrpcstart"

	"github.com/OpenIMSDK/chat/internal/rpc/admin"
	"github.com/OpenIMSDK/chat/pkg/common/config"
)

func main() {
	if err := config.InitConfig(); err != nil {
		panic(err)
	}
	if err := log.InitFromConfig("chat.log", "admin-rpc", *config.Config.Log.RemainLogLevel, *config.Config.Log.IsStdout, *config.Config.Log.IsJson, *config.Config.Log.StorageLocation, *config.Config.Log.RemainRotationCount); err != nil {
		panic(err)
	}
	err := chatrpcstart.Start(config.Config.RpcPort.OpenImAdminPort[0], config.Config.RpcRegisterName.OpenImAdminName, 0, admin.Start)
	if err != nil {
		panic(err)
	}
}
