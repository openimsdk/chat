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

package config

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	Constant "github.com/OpenIMSDK/chat/pkg/common/constant"
	"github.com/OpenIMSDK/tools/utils"
	"gopkg.in/yaml.v3"
)

var (
	_, b, _, _ = runtime.Caller(0)
	// Root folder of this project.
	Root = filepath.Join(filepath.Dir(b), "../..")
)

func readConfig(configFile string) ([]byte, error) {
	b, err := os.ReadFile(configFile)
	if err != nil { // File exists and was read successfully
		return nil, utils.Wrap(err, configFile)
	}
	return b, nil

	//	//First, check the configFile argument
	//	if configFile != "" {
	//		b, err := os.ReadFile(configFile)
	//		if err == nil { // File exists and was read successfully
	//			return b, nil
	//		}
	//	}
	//
	//	// Second, check for OPENIMCHATCONFIG environment variable
	//	envConfigPath := os.Getenv("OPENIMCHATCONFIG")
	//	if envConfigPath != "" {
	//		b, err := os.ReadFile(envConfigPath)
	//		if err == nil { // File exists and was read successfully
	//			return b, nil
	//		}
	//		// Again, if there was an error, you can either log it or ignore.
	//	}
	//
	//	// If neither configFile nor environment variable provided a valid path, use default path
	//	defaultConfigPath := filepath.Join(Root, "config", "config.yaml")
	//	b, err := os.ReadFile(defaultConfigPath)
	//	if err != nil {
	//		return nil, utils.Wrap(err, defaultConfigPath)
	//	}
	//	return b, nil
}

func InitConfig(configFile string) error {
	data, err := readConfig(configFile)
	if err != nil {
		return fmt.Errorf("read loacl config file error: %w", err)
	}

	if err := yaml.NewDecoder(bytes.NewReader(data)).Decode(&Config); err != nil {
		return fmt.Errorf("parse loacl openIMConfig file error: %w", err)
	}
	if err != nil {
		return utils.Wrap(err, configFile)
	}

	if err := configGetEnv(); err != nil {
		return fmt.Errorf("get env error:%w", err)
	}

	configData, err := yaml.Marshal(&Config)
	fmt.Printf("debug: %s\nconfig:\n%s\n", time.Now(), string(configData))

	return nil
}

func configFieldCopy[T any](local **T, remote T) {
	if *local == nil {
		*local = &remote
	}
}

func GetDefaultIMAdmin() string {
	return Config.AdminList[0].ImAdminID
}

func GetIMAdmin(chatAdminID string) string {
	for _, admin := range Config.AdminList {
		if admin.ImAdminID == chatAdminID {
			return admin.ImAdminID
		}
	}
	for _, admin := range Config.AdminList {
		if admin.AdminID == chatAdminID {
			return admin.ImAdminID
		}
	}
	if id := getChatAdmin(chatAdminID); id != "" {
		return id
	}
	return ""
}

func getChatAdmin(chatAdminID string) string {
	for _, admin := range Config.ChatAdmin {
		if admin.ImAdminID == chatAdminID {
			return admin.ImAdminID
		}
	}
	return ""
}

type zkLogger struct{}

func (l *zkLogger) Printf(format string, a ...interface{}) {
	fmt.Printf("zk get config %s\n", fmt.Sprintf(format, a...))
}

func checkFileExists(filePath string) bool {
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func findConfigFile(paths []string) (string, error) {
	for _, path := range paths {
		if checkFileExists(path) {
			return path, nil
		}
	}
	return "", fmt.Errorf("configPath not found")
}

func CreateCatalogPath(path string) []string {

	path1 := filepath.Dir(path)
	path1 = filepath.Dir(path1)
	// the parent of  binary file
	pa1 := filepath.Join(path1, Constant.ConfigPath)
	path2 := filepath.Dir(path1)
	path2 = filepath.Dir(path2)
	path2 = filepath.Dir(path2)
	// the parent is _output
	pa2 := filepath.Join(path2, Constant.ConfigPath)
	path3 := filepath.Dir(path2)
	// the parent is project(default)
	pa3 := filepath.Join(path3, Constant.ConfigPath)

	return []string{pa1, pa2, pa3}

}

func findConfigPath(configFile string) (string, error) {
	path := make([]string, 10)

	// First, check the configFile argument
	if configFile != "" {
		if _, err := findConfigFile([]string{configFile}); err != nil {
			return "", errors.New("the configFile argument path is error")
		}
		fmt.Println("configfile:", configFile)
		return configFile, nil
	}

	// Second, check for OPENIMCONFIG environment variable
	//envConfigPath := os.Getenv(Constant.OpenIMConfig)
	envConfigPath := os.Getenv(Constant.OpenIMConfig)
	if envConfigPath != "" {
		if _, err := findConfigFile([]string{envConfigPath}); err != nil {
			return "", errors.New("the environment path config path is error")
		}
		return envConfigPath, nil
	}
	// Third, check the catalog to find the config.yaml

	p1, err := os.Executable()
	if err != nil {
		return "", err
	}

	path = CreateCatalogPath(p1)
	pathFind, err := findConfigFile(path)
	if err == nil {
		return pathFind, nil
	}

	// Forth, use the Default path.
	return "", errors.New("the config.yaml path not found")
}

func FlagParse() (string, int, bool, bool, error) {
	var configFile string
	flag.StringVar(&configFile, "config_folder_path", "", "Config full path")

	var ginPort int
	flag.IntVar(&ginPort, "port", 10009, "get ginServerPort from cmd")

	var hide bool
	flag.BoolVar(&hide, "hide", false, "hide the ComponentCheck result")

	// Version flag
	var showVersion bool
	flag.BoolVar(&showVersion, "version", false, "show version and exit")

	flag.Parse()

	configFile, err := findConfigPath(configFile)
	if err != nil {
		return "", 0, false, false, err
	}
	return configFile, ginPort, hide, showVersion, nil
}

func configGetEnv() error {
	Config.Envs.Discovery = getEnv("ENVS_DISCOVERY", Config.Envs.Discovery)
	Config.Zookeeper.Schema = getEnv("ZOOKEEPER_SCHEMA", Config.Zookeeper.Schema)
	Config.Zookeeper.Username = getEnv("ZOOKEEPER_USERNAME", Config.Zookeeper.Username)
	Config.Zookeeper.Password = getEnv("ZOOKEEPER_PASSWORD", Config.Zookeeper.Password)

	Config.ChatApi.ListenIP = getEnv("CHAT_API_LISTEN_IP", Config.ChatApi.ListenIP)
	Config.AdminApi.ListenIP = getEnv("ADMIN_API_LISTEN_IP", Config.AdminApi.ListenIP)
	Config.Rpc.RegisterIP = getEnv("RPC_REGISTER_IP", Config.Rpc.RegisterIP)
	Config.Rpc.ListenIP = getEnv("RPC_LISTEN_IP", Config.Rpc.ListenIP)

	Config.Mysql.Username = getEnvStringPoint("MYSQL_USERNAME", Config.Mysql.Username)
	Config.Mysql.Password = getEnvStringPoint("MYSQL_PASSWORD", Config.Mysql.Password)
	Config.Mysql.Database = getEnvStringPoint("MYSQL_DATABASE", Config.Mysql.Database)
	Config.Mysql.Address = getArrPointEnv("MYSQL_ADDRESS", "MYSQL_PORT", Config.Mysql.Address)

	Config.Log.StorageLocation = getEnvStringPoint("LOG_STORAGE_LOCATION", Config.Log.StorageLocation)

	Config.Secret = getEnvStringPoint("SECRET", Config.Secret)
	Config.ProxyHeader = getEnv("PROXY_HEADER", Config.ProxyHeader)
	Config.OpenIMUrl = getStringEnv("OPENIM_SERVER_ADDRESS", "API_OPENIM_PORT", Config.OpenIMUrl)

	Config.Redis.Username = getEnv("REDIS_USERNAME", Config.Redis.Username)
	Config.Redis.Password = getEnv("REDIS_PASSWORD", Config.Redis.Password)
	Config.Redis.Address = getArrPointEnv("REDIS_ADDRESS", "REDIS_PORT", Config.Redis.Address)

	var err error
	Config.TokenPolicy.Expire, err = getEnvIntPoint("TOKEN_EXPIRE", Config.TokenPolicy.Expire)
	if err != nil {
		return err
	}
	getArrEnv("ZOOKEEPER_ADDRESS", "ZOOKEEPER_PORT", Config.Zookeeper.ZkAddr)
	return nil
}

func getArrEnv(key1, key2 string, fallback []string) {
	str1 := getEnv(key1, "")
	str2 := getEnv(key2, "")
	str := fmt.Sprintf("%s:%s", str1, str2)
	arr := make([]string, 1)
	if len(str) <= 1 {
		return
	}
	arr[0] = str
	fmt.Println("zookeeper Envirement valiable", "str", str)
	Config.Zookeeper.ZkAddr = arr
}

func getArrPointEnv(key1, key2 string, fallback *[]string) *[]string {
	str1 := getEnv(key1, "")
	str2 := getEnv(key2, "")
	str := fmt.Sprintf("%s:%s", str1, str2)
	if len(str) <= 1 {
		return fallback
	}
	return &[]string{str}
}

func getStringEnv(key1, key2 string, fallback string) string {
	str1 := getEnv(key1, "")
	str2 := getEnv(key2, "")
	str := fmt.Sprintf("%s:%s", str1, str2)
	if len(str) <= 2 {
		return fallback
	}
	return str
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvStringPoint(key string, fallback *string) *string {
	if value, exists := os.LookupEnv(key); exists {
		return &value
	}
	return fallback
}

func getEnvIntPoint(key string, fallback *int64) (*int64, error) {
	if value, exists := os.LookupEnv(key); exists {
		val, err := strconv.Atoi(value)
		temp := int64(val)
		if err != nil {
			return nil, err
		}
		return &temp, nil
	}
	return fallback, nil
}
