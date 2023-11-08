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
	"time"

	Constant "github.com/OpenIMSDK/chat/pkg/common/constant"
	"github.com/OpenIMSDK/protocol/constant"
	openKeeper "github.com/OpenIMSDK/tools/discoveryregistry/zookeeper"

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
	//			fmt.Println("这里aaaaaaaa")
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
	if Config.Envs.Discovery != "k8s" {
		zk, err := openKeeper.NewClient(Config.Zookeeper.ZkAddr, Config.Zookeeper.Schema,
			openKeeper.WithFreq(time.Hour), openKeeper.WithUserNameAndPassword(Config.Zookeeper.Username,
				Config.Zookeeper.Password), openKeeper.WithRoundRobin(), openKeeper.WithTimeout(10), openKeeper.WithLogger(&zkLogger{}))
		if err != nil {
			return utils.Wrap(err, "conn zk error ")
		}
		defer zk.Close()
		var openIMConfigData []byte
		for i := 0; i < 100; i++ {
			var err error
			configData, err := zk.GetConfFromRegistry(constant.OpenIMCommonConfigKey)
			if err != nil {
				fmt.Printf("get zk config [%d] error: %v\n;envs.descoery=%s", i, err, Config.Envs.Discovery)
				time.Sleep(time.Second)
				continue
			}
			if len(configData) == 0 {
				fmt.Printf("get zk config [%d] data is empty\n", i)
				time.Sleep(time.Second)
				continue
			}
			openIMConfigData = configData
		}
		if len(openIMConfigData) == 0 {
			return errors.New("get zk config data failed")
		}
		if err := yaml.NewDecoder(bytes.NewReader(openIMConfigData)).Decode(&imConfig); err != nil {
			return fmt.Errorf("parse zk openIMConfig: %w", err)
		}
		// 这里可以优化，可将其优化为结构体层面的赋值
		configFieldCopy(&Config.Mysql.Address, imConfig.Mysql.Address)
		configFieldCopy(&Config.Mysql.Username, imConfig.Mysql.Username)
		configFieldCopy(&Config.Mysql.Password, imConfig.Mysql.Password)
		configFieldCopy(&Config.Mysql.Database, imConfig.Mysql.Database)
		configFieldCopy(&Config.Mysql.MaxOpenConn, imConfig.Mysql.MaxOpenConn)
		configFieldCopy(&Config.Mysql.MaxIdleConn, imConfig.Mysql.MaxIdleConn)
		configFieldCopy(&Config.Mysql.MaxLifeTime, imConfig.Mysql.MaxLifeTime)
		configFieldCopy(&Config.Mysql.LogLevel, imConfig.Mysql.LogLevel)
		configFieldCopy(&Config.Mysql.SlowThreshold, imConfig.Mysql.SlowThreshold)

		configFieldCopy(&Config.Log.StorageLocation, imConfig.Log.StorageLocation)
		configFieldCopy(&Config.Log.RotationTime, imConfig.Log.RotationTime)
		configFieldCopy(&Config.Log.RemainRotationCount, imConfig.Log.RemainRotationCount)
		configFieldCopy(&Config.Log.RemainLogLevel, imConfig.Log.RemainLogLevel)
		configFieldCopy(&Config.Log.IsStdout, imConfig.Log.IsStdout)
		configFieldCopy(&Config.Log.WithStack, imConfig.Log.WithStack)
		configFieldCopy(&Config.Log.IsJson, imConfig.Log.IsJson)

		configFieldCopy(&Config.Secret, imConfig.Secret)
		configFieldCopy(&Config.TokenPolicy.Expire, imConfig.TokenPolicy.Expire)

		// Redis
		configFieldCopy(&Config.Redis.Address, imConfig.Redis.Address)
		configFieldCopy(&Config.Redis.Password, imConfig.Redis.Password)
		configFieldCopy(&Config.Redis.Username, imConfig.Redis.Username)
	}

	configData, err := yaml.Marshal(&Config)
	fmt.Printf("debug: %s\nconfig:\n%s\n", time.Now(), string(configData))
	if err != nil {
		return utils.Wrap(err, configFile)
	}
	fmt.Printf("%s\nconfig:\n%s\n", time.Now(), string(configData))

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
	path2 := filepath.Join(path1, Constant.ConfigPath)
	path2 = filepath.Dir(path1)
	path2 = filepath.Dir(path2)
	path2 = filepath.Dir(path2)
	// the parent is _output
	path3 := filepath.Join(path2, Constant.ConfigPath)
	path2 = filepath.Dir(path2)
	// the parent is project(default)
	path4 := filepath.Join(path2, Constant.ConfigPath)

	return []string{path1, path3, path4}

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
	return Constant.ConfigPath, nil
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
