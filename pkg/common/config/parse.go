package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	openIMConfig "github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
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
	zk, err := openKeeper.NewClient(Config.Zookeeper.ZkAddr, Config.Zookeeper.Schema,
		openKeeper.WithFreq(time.Hour), openKeeper.WithUserNameAndPassword(Config.Zookeeper.UserName,
			Config.Zookeeper.Password), openKeeper.WithRoundRobin(), openKeeper.WithTimeout(10), openKeeper.WithLogger(&zkLogger{}))
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

	configFieldCopy(&Config.Mysql.DBAddress, openIMConfig.Config.Mysql.DBAddress)
	configFieldCopy(&Config.Mysql.DBUserName, openIMConfig.Config.Mysql.DBUserName)
	configFieldCopy(&Config.Mysql.DBPassword, openIMConfig.Config.Mysql.DBPassword)
	configFieldCopy(&Config.Mysql.DBDatabaseName, openIMConfig.Config.Mysql.DBDatabaseName)
	configFieldCopy(&Config.Mysql.DBMaxOpenConns, openIMConfig.Config.Mysql.DBMaxOpenConns)
	configFieldCopy(&Config.Mysql.DBMaxIdleConns, openIMConfig.Config.Mysql.DBMaxIdleConns)
	configFieldCopy(&Config.Mysql.DBMaxLifeTime, openIMConfig.Config.Mysql.DBMaxLifeTime)
	configFieldCopy(&Config.Mysql.LogLevel, openIMConfig.Config.Mysql.LogLevel)
	configFieldCopy(&Config.Mysql.SlowThreshold, openIMConfig.Config.Mysql.SlowThreshold)

	configFieldCopy(&Config.Log.StorageLocation, openIMConfig.Config.Log.StorageLocation)
	configFieldCopy(&Config.Log.RotationTime, openIMConfig.Config.Log.RotationTime)
	configFieldCopy(&Config.Log.RemainRotationCount, openIMConfig.Config.Log.RemainRotationCount)
	configFieldCopy(&Config.Log.RemainLogLevel, openIMConfig.Config.Log.RemainLogLevel)
	configFieldCopy(&Config.Log.IsStdout, openIMConfig.Config.Log.IsStdout)
	configFieldCopy(&Config.Log.WithStack, openIMConfig.Config.Log.WithStack)
	configFieldCopy(&Config.Log.IsJson, openIMConfig.Config.Log.IsJson)

	jsonData, err := json.Marshal(Config)
	if err != nil {
		return err
	}
	fmt.Println(string(jsonData))

	return nil
}

func configFieldCopy[T any](local **T, remote T) {
	if *local == nil {
		*local = &remote
	}
}

type zkLogger struct{}

func (l *zkLogger) Printf(format string, a ...interface{}) {
	fmt.Printf("zk get config %s\n", fmt.Sprintf(format, a...))
}
