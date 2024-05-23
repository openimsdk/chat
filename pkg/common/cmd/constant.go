// Copyright © 2023 OpenIM. All rights reserved.
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

package cmd

import (
	"strings"
)

var (
	ShareFileName           = "share.yml"
	RedisConfigFileName     = "redis.yml"
	DiscoveryConfigFileName = "discovery.yml"
	MongodbConfigFileName   = "mongodb.yml"
	LogConfigFileName       = "log.yml"
	ChatAPIAdminCfgFileName = "chat-api-admin.yml"
	ChatAPIChatCfgFileName  = "chat-api-chat.yml"
	ChatRPCAdminCfgFileName = "chat-rpc-admin.yml"
	ChatRPCChatCfgFileName  = "chat-rpc-chat.yml"
)

var ConfigEnvPrefixMap map[string]string

func init() {
	ConfigEnvPrefixMap = make(map[string]string)
	fileNames := []string{
		ShareFileName,
		RedisConfigFileName,
		DiscoveryConfigFileName,
		MongodbConfigFileName,
		LogConfigFileName,
		ChatAPIAdminCfgFileName,
		ChatAPIChatCfgFileName,
		ChatRPCAdminCfgFileName,
		ChatRPCChatCfgFileName,
	}

	for _, fileName := range fileNames {
		envKey := strings.TrimSuffix(strings.TrimSuffix(fileName, ".yml"), ".yaml")
		envKey = "CHATENV_" + envKey
		envKey = strings.ToUpper(strings.ReplaceAll(envKey, "-", "_"))
		ConfigEnvPrefixMap[fileName] = envKey
	}
}

const (
	FlagConf          = "config_folder_path"
	FlagTransferIndex = "index"
)
