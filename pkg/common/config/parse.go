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
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/OpenIMSDK/protocol/constant"
	openKeeper "github.com/OpenIMSDK/tools/discoveryregistry/zookeeper"
	"github.com/OpenIMSDK/tools/utils"
	"gopkg.in/yaml.v3"
)

var (
	_, b, _, _ = runtime.Caller(0)
	// Root folder of this project.
	Root = filepath.Join(filepath.Dir(b), "../../..")
)

func readConfig(configFile string) ([]byte, error) {
	b, err := os.ReadFile(configFile)
	if err != nil {
		return nil, utils.Wrap(err, configFile)
	}
	return b, nil
	// cfgName := os.Getenv("CONFIG_NAME")
	// if len(cfgName) != 0 {
	// 	data, err := os.ReadFile(filepath.Join(cfgName, "config", "config.yaml"))
	// 	if err != nil {
	// 		data, err = os.ReadFile(filepath.Join(Root, "config", "config.yaml"))
	// 		if err != nil {
	// 			return nil, err
	// 		}
	// 	} else {
	// 		Root = cfgName
	// 	}
	// 	return data, nil
	// } else {
	// 	return os.ReadFile(fmt.Sprintf("../config/%s", "config.yaml"))
	// }
}

func InitConfig(configFile string) error {
	data, err := readConfig(configFile)
	if err != nil {
		return fmt.Errorf("read loacl config file error: %w", err)
	}
	if err := yaml.NewDecoder(bytes.NewReader(data)).Decode(&Config); err != nil {
		return fmt.Errorf("parse loacl openIMConfig file error: %w", err)
	}
	zk, err := openKeeper.NewClient(Config.Zookeeper.ZkAddr, Config.Zookeeper.Schema,
		openKeeper.WithFreq(time.Hour), openKeeper.WithUserNameAndPassword(Config.Zookeeper.Username,
			Config.Zookeeper.Password), openKeeper.WithRoundRobin(), openKeeper.WithTimeout(10), openKeeper.WithLogger(&zkLogger{}))
	if err != nil {
		return utils.Wrap(err, "conn zk error ")
	}
	defer zk.CloseZK()
	var openIMConfigData []byte
	for i := 0; i < 100; i++ {
		var err error
		configData, err := zk.GetConfFromRegistry(constant.OpenIMCommonConfigKey)
		if err != nil {
			fmt.Printf("get zk config [%d] error: %v\n", i, err)
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
