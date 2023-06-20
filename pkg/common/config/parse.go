package config

import (
	"bytes"
	"fmt"
	openIMConfig "github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	openKeeper "github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry/zookeeper"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

var (
	_, b, _, _ = runtime.Caller(0)
	// Root folder of this project
	Root = filepath.Join(filepath.Dir(b), "../../..")
)

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

func InitConfig() error {
	data, err := readConfig()
	if err != nil {
		return fmt.Errorf("read loacl config file error: %w", err)
	}
	if err := yaml.NewDecoder(bytes.NewReader(data)).Decode(&Config); err != nil {
		return fmt.Errorf("parse loacl openIMConfig file error: %w", err)
	}
	zk, err := openKeeper.NewClient(openIMConfig.Config.Zookeeper.ZkAddr, openIMConfig.Config.Zookeeper.Schema,
		openKeeper.WithFreq(time.Hour), openKeeper.WithUserNameAndPassword(openIMConfig.Config.Zookeeper.UserName,
			openIMConfig.Config.Zookeeper.Password), openKeeper.WithRoundRobin(), openKeeper.WithTimeout(10), openKeeper.WithLogger(log.NewZkLogger()))
	if err != nil {
		return fmt.Errorf("conn zk error: %w", err)
	}
	openIMConfigData, err := zk.GetConfFromRegistry(constant.OpenIMCommonConfigKey)
	if err != nil {
		return fmt.Errorf("get zk config: %w", err)
	}
	if err := yaml.NewDecoder(bytes.NewReader(openIMConfigData)).Decode(&openIMConfig.Config); err != nil {
		return fmt.Errorf("parse zk openIMConfig: %w", err)
	}
	Config.Mysql.DBAddress = openIMConfig.Config.Mysql.DBAddress
	Config.Mysql.DBUserName = openIMConfig.Config.Mysql.DBUserName
	Config.Mysql.DBPassword = openIMConfig.Config.Mysql.DBPassword
	Config.Mysql.DBMaxOpenConns = openIMConfig.Config.Mysql.DBMaxOpenConns
	Config.Mysql.DBMaxIdleConns = openIMConfig.Config.Mysql.DBMaxIdleConns
	Config.Mysql.DBMaxLifeTime = openIMConfig.Config.Mysql.DBMaxLifeTime
	Config.Mysql.LogLevel = openIMConfig.Config.Mysql.LogLevel
	Config.Mysql.SlowThreshold = openIMConfig.Config.Mysql.SlowThreshold

	Config.Log.StorageLocation = openIMConfig.Config.Log.StorageLocation
	Config.Log.RotationTime = openIMConfig.Config.Log.RotationTime
	Config.Log.RemainRotationCount = openIMConfig.Config.Log.RemainRotationCount
	Config.Log.RemainLogLevel = openIMConfig.Config.Log.RemainLogLevel
	Config.Log.IsStdout = openIMConfig.Config.Log.IsStdout
	Config.Log.WithStack = openIMConfig.Config.Log.WithStack
	Config.Log.IsJson = openIMConfig.Config.Log.IsJson

	return nil
}
