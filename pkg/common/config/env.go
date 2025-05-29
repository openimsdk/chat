package config

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
	ChatAPIBotCfgFileName   = "chat-api-bot.yml"
	ChatRPCAdminCfgFileName = "chat-rpc-admin.yml"
	ChatRPCChatCfgFileName  = "chat-rpc-chat.yml"
	ChatRPCBotCfgFileName   = "chat-rpc-bot.yml"
)

var EnvPrefixMap map[string]string

func init() {
	EnvPrefixMap = make(map[string]string)
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
		EnvPrefixMap[fileName] = envKey
	}
}

const (
	FlagConf          = "config_folder_path"
	FlagTransferIndex = "index"
)
