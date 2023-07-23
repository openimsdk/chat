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

	openIMConfig "github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	openKeeper "github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry/zookeeper"
	"gopkg.in/yaml.v3"
)

var (
	_, b, _, _ = runtime.Caller(0)
	// Root folder of this project.
	Root = filepath.Join(filepath.Dir(b), "../../..")
)

// read rpc config
func readConfig() ([]byte, error) {
	cfgName := os.Getenv("CONFIG_NAME")
	if len(cfgName) != 0 {
		data, err := os.ReadFile(filepath.Join(cfgName, "config", "config.yaml"))
		if err != nil {
			data, err = os.ReadFile(filepath.Join(Root, "config", "config.yaml"))
			if err != nil {
				return nil, err
			}
		} else {
			Root = cfgName
		}
		return data, nil
	} else {
		return os.ReadFile(fmt.Sprintf("../config/%s", "config.yaml"))
	}
}

// initialize config
func InitConfig() error {
	data, err := readConfig()
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
		return fmt.Errorf("conn zk error: %w", err)
	}
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
	if err := yaml.NewDecoder(bytes.NewReader(openIMConfigData)).Decode(&openIMConfig.Config); err != nil {
		return fmt.Errorf("parse zk openIMConfig: %w", err)
	}
	configFieldCopy(&Config.Mysql.Address, openIMConfig.Config.Mysql.Address)
	configFieldCopy(&Config.Mysql.Username, openIMConfig.Config.Mysql.Username)
	configFieldCopy(&Config.Mysql.Password, openIMConfig.Config.Mysql.Password)
	configFieldCopy(&Config.Mysql.Database, openIMConfig.Config.Mysql.Database)
	configFieldCopy(&Config.Mysql.MaxOpenConn, openIMConfig.Config.Mysql.MaxOpenConn)
	configFieldCopy(&Config.Mysql.MaxIdleConn, openIMConfig.Config.Mysql.MaxIdleConn)
	configFieldCopy(&Config.Mysql.MaxLifeTime, openIMConfig.Config.Mysql.MaxLifeTime)
	configFieldCopy(&Config.Mysql.LogLevel, openIMConfig.Config.Mysql.LogLevel)
	configFieldCopy(&Config.Mysql.SlowThreshold, openIMConfig.Config.Mysql.SlowThreshold)

	configFieldCopy(&Config.Log.StorageLocation, openIMConfig.Config.Log.StorageLocation)
	configFieldCopy(&Config.Log.RotationTime, openIMConfig.Config.Log.RotationTime)
	configFieldCopy(&Config.Log.RemainRotationCount, openIMConfig.Config.Log.RemainRotationCount)
	configFieldCopy(&Config.Log.RemainLogLevel, openIMConfig.Config.Log.RemainLogLevel)
	configFieldCopy(&Config.Log.IsStdout, openIMConfig.Config.Log.IsStdout)
	configFieldCopy(&Config.Log.WithStack, openIMConfig.Config.Log.WithStack)
	configFieldCopy(&Config.Log.IsJson, openIMConfig.Config.Log.IsJson)

	configFieldCopy(&Config.Secret, openIMConfig.Config.Secret)
	configFieldCopy(&Config.TokenPolicy.Expire, openIMConfig.Config.TokenPolicy.Expire)

	configData, err := yaml.Marshal(&Config)
	if err != nil {
		return err
	}
	fmt.Printf("%s\nconfig:\n%s\n", time.Now(), string(configData))

	return nil
}

// copy config field
func configFieldCopy[T any](local **T, remote T) {
	if *local == nil {
		*local = &remote
	}
}

// get default admin by im
func GetDefaultIMAdmin() string {
	return Config.AdminList[0].ImAdminID
}

// gey admin by im
func GetIMAdmin(chatAdminID string) string {
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
