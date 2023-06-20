package config

import (
	"bytes"
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

//func init()  {
//	cfgName := os.Getenv("CONFIG_NAME")
//	if len(cfgName) != 0 {
//		bytes, err := os.ReadFile(filepath.Join(cfgName, "openIMConfig", "openIMConfig.yaml"))
//		if err != nil {
//			bytes, err = os.ReadFile(filepath.Join(Root, "openIMConfig", "openIMConfig.yaml"))
//			if err != nil {
//				panic(err.Error() + " openIMConfig: " + filepath.Join(cfgName, "openIMConfig", "openIMConfig.yaml"))
//			}
//		} else {
//			Root = cfgName
//		}
//		if err = yaml.Unmarshal(bytes, &Config); err != nil {
//			panic(err.Error())
//		}
//	} else {
//		bytes, err := os.ReadFile(fmt.Sprintf("../openIMConfig/%s", "openIMConfig.yaml"))
//		if err != nil {
//			panic(err.Error())
//		}
//		if err = yaml.Unmarshal(bytes, &Config); err != nil {
//			panic(err.Error())
//		}
//	}
//
//	syncEtcd()
//
//}

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

//func readConfig() ([]byte, error) {
//	path := filepath.Join("./config", "config.yaml")
//	absPath, err := filepath.Abs(path)
//	if err != nil {
//		return nil, err
//	}
//	return os.ReadFile(absPath)
//}

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
			Config.Zookeeper.Password), openKeeper.WithRoundRobin(), openKeeper.WithTimeout(10))
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
	//Config.Mysql.DBDatabaseName = openIMConfig.Config.Mysql.DBDatabaseName
	Config.Mysql.DBMaxOpenConns = openIMConfig.Config.Mysql.DBMaxOpenConns
	Config.Mysql.DBMaxIdleConns = openIMConfig.Config.Mysql.DBMaxIdleConns
	Config.Mysql.DBMaxLifeTime = openIMConfig.Config.Mysql.DBMaxLifeTime
	Config.Mysql.LogLevel = openIMConfig.Config.Mysql.LogLevel
	Config.Mysql.SlowThreshold = openIMConfig.Config.Mysql.SlowThreshold

	//Config.Log.StorageLocation = openIMConfig.Config.Log.StorageLocation
	//Config.Log.RotationTime = openIMConfig.Config.Log.RotationTime
	//Config.Log.RemainRotationCount = openIMConfig.Config.Log.RemainRotationCount
	//Config.Log.RemainLogLevel = openIMConfig.Config.Log.RemainLogLevel

	Config.Log.StorageLocation = openIMConfig.Config.Log.StorageLocation
	Config.Log.RotationTime = openIMConfig.Config.Log.RotationTime
	Config.Log.RemainRotationCount = openIMConfig.Config.Log.RemainRotationCount
	Config.Log.RemainLogLevel = openIMConfig.Config.Log.RemainLogLevel
	Config.Log.IsStdout = openIMConfig.Config.Log.IsStdout
	Config.Log.WithStack = openIMConfig.Config.Log.WithStack
	Config.Log.IsJson = openIMConfig.Config.Log.IsJson

	return nil
}

//func syncEtcd() {
//	zk, err := openKeeper.NewClient(openIMConfig.Config.Zookeeper.ZkAddr, openIMConfig.Config.Zookeeper.Schema, 10, openIMConfig.Config.Zookeeper.UserName, openIMConfig.Config.Zookeeper.Password)
//	if err != nil {
//		return err
//	}
//	if err := configEtcd.Load(Config.Etcd.EtcdSchema, Config.Etcd.EtcdAddr, Config.Etcd.UserName, Config.Etcd.Password, Config.Etcd.Secret); err != nil {
//		panic(err.Error() + " load etcd openIMConfig error ")
//	}
//
//	// #################### MYSQL ####################
//	Config.Mysql.DBAddress = configEtcd.Config.Mysql.DBAddress
//	Config.Mysql.DBUserName = configEtcd.Config.Mysql.DBUserName
//	Config.Mysql.DBPassword = configEtcd.Config.Mysql.DBPassword
//	//Config.Mysql.DBDatabaseName = configEtcd.Config.Mysql.DBDatabaseName
//	Config.Mysql.DBMaxOpenConns = configEtcd.Config.Mysql.DBMaxOpenConns
//	Config.Mysql.DBMaxIdleConns = configEtcd.Config.Mysql.DBMaxIdleConns
//	Config.Mysql.DBMaxLifeTime = configEtcd.Config.Mysql.DBMaxLifeTime
//	Config.Mysql.LogLevel = configEtcd.Config.Mysql.LogLevel
//	Config.Mysql.SlowThreshold = configEtcd.Config.Mysql.SlowThreshold
//
//	// #################### REDIS ####################
//	Config.Redis.DBAddress = configEtcd.Config.Redis.DBAddress
//	Config.Redis.DBUserName = configEtcd.Config.Redis.DBUserName
//	Config.Redis.DBPassWord = configEtcd.Config.Redis.DBPassWord
//	Config.Redis.EnableCluster = configEtcd.Config.Redis.EnableCluster
//
//	Config.RpcRegisterIP = configEtcd.Config.RpcRegisterIP
//	Config.RpcListenIP = "0.0.0.0"
//	// #################### Config.RpcPort ####################
//
//	Config.RpcPort.OpenImUserPort = configEtcd.Config.RpcPort.OpenImUserPort
//	Config.RpcPort.OpenImFriendPort = configEtcd.Config.RpcPort.OpenImFriendPort
//	Config.RpcPort.OpenImMessagePort = configEtcd.Config.RpcPort.OpenImMessagePort
//	Config.RpcPort.OpenImMessageGatewayPort = configEtcd.Config.RpcPort.OpenImMessageGatewayPort
//	Config.RpcPort.OpenImGroupPort = configEtcd.Config.RpcPort.OpenImGroupPort
//	Config.RpcPort.OpenImAuthPort = configEtcd.Config.RpcPort.OpenImAuthPort
//	Config.RpcPort.OpenImPushPort = configEtcd.Config.RpcPort.OpenImPushPort
//	//Config.RpcPort.OpenImAdminPort = configEtcd.Config.RpcPort.OpenImAdminPort
//	//Config.RpcPort.OpenImChatPort = configEtcd.Config.RpcPort.OpenImChatPort
//	Config.RpcPort.OpenImOfficePort = configEtcd.Config.RpcPort.OpenImOfficePort
//	Config.RpcPort.OpenImOrganizationPort = configEtcd.Config.RpcPort.OpenImOrganizationPort
//	Config.RpcPort.OpenImConversationPort = configEtcd.Config.RpcPort.OpenImConversationPort
//	Config.RpcPort.OpenImCachePort = configEtcd.Config.RpcPort.OpenImCachePort
//	Config.RpcPort.OpenImRealTimeCommPort = configEtcd.Config.RpcPort.OpenImRealTimeCommPort
//
//	// #################### Config.RpcRegisterName ####################
//
//	Config.RpcRegisterName.OpenImUserName = configEtcd.Config.RpcRegisterName.OpenImUserName
//	Config.RpcRegisterName.OpenImFriendName = configEtcd.Config.RpcRegisterName.OpenImFriendName
//	Config.RpcRegisterName.OpenImMsgName = configEtcd.Config.RpcRegisterName.OpenImMsgName
//	Config.RpcRegisterName.OpenImPushName = configEtcd.Config.RpcRegisterName.OpenImPushName
//	Config.RpcRegisterName.OpenImRelayName = configEtcd.Config.RpcRegisterName.OpenImRelayName
//	Config.RpcRegisterName.OpenImGroupName = configEtcd.Config.RpcRegisterName.OpenImGroupName
//	Config.RpcRegisterName.OpenImAuthName = configEtcd.Config.RpcRegisterName.OpenImAuthName
//	//Config.RpcRegisterName.OpenImAdminCMSName = configEtcd.Config.RpcRegisterName.OpenImAdminCMSName
//	//Config.RpcRegisterName.OpenImChatName = configEtcd.Config.RpcRegisterName.OpenImChatName
//	Config.RpcRegisterName.OpenImOfficeName = configEtcd.Config.RpcRegisterName.OpenImOfficeName
//	Config.RpcRegisterName.OpenImOrganizationName = configEtcd.Config.RpcRegisterName.OpenImOrganizationName
//	Config.RpcRegisterName.OpenImConversationName = configEtcd.Config.RpcRegisterName.OpenImConversationName
//	Config.RpcRegisterName.OpenImCacheName = configEtcd.Config.RpcRegisterName.OpenImCacheName
//	Config.RpcRegisterName.OpenImRealTimeCommName = configEtcd.Config.RpcRegisterName.OpenImRealTimeCommName
//
//	// #################### LOG ####################
//	Config.Log.StorageLocation = configEtcd.Config.Log.StorageLocation
//	Config.Log.RotationTime = configEtcd.Config.Log.RotationTime
//	Config.Log.RemainRotationCount = configEtcd.Config.Log.RemainRotationCount
//	Config.Log.RemainLogLevel = configEtcd.Config.Log.RemainLogLevel
//
//	// #################### TokenPolicy ####################
//	Config.TokenPolicy.AccessSecret = configEtcd.Config.TokenPolicy.AccessSecret
//	Config.TokenPolicy.AccessExpire = configEtcd.Config.TokenPolicy.AccessExpire
//
//	// #################### MINIO ####################
//	// old version
//	Config.Credential.Minio.Bucket = configEtcd.Config.Credential.Minio.Bucket
//	Config.Credential.Minio.Endpoint = configEtcd.Config.Credential.Minio.Endpoint
//	Config.Credential.Minio.AccessKeyID = configEtcd.Config.Credential.Minio.AccessKeyID
//	Config.Credential.Minio.SecretAccessKey = configEtcd.Config.Credential.Minio.SecretAccessKey
//	Config.Credential.Minio.EndpointInner = configEtcd.Config.Credential.Minio.EndpointInner
//	Config.Credential.Minio.EndpointInnerEnable = configEtcd.Config.Credential.Minio.EndpointInnerEnable
//
//	// new version
//	Config.OSS.Minio.Endpoint = configEtcd.Config.Credential.Minio.Endpoint
//	Config.OSS.Minio.AccessKeyID = configEtcd.Config.Credential.Minio.AccessKeyID
//	Config.OSS.Minio.SecretAccessKey = configEtcd.Config.Credential.Minio.SecretAccessKey
//	Config.OSS.Minio.Secure = false
//	Config.OSS.Minio.AccessAddress = configEtcd.Config.Credential.Minio.Endpoint
//	Config.OSS.Minio.Bucket = configEtcd.Config.Credential.Minio.Bucket
//
//	if err := minioConf(); err != nil {
//		panic(err)
//	}
//
//	temp, _ := json.Marshal(Config)
//	fmt.Printf("openIMConfig: %s\n", string(temp))
//
//}
//
//func minioConf() error {
//	if strings.Index(Config.OSS.Minio.Endpoint, "http://") == 0 || strings.Index(Config.OSS.Minio.Endpoint, "https://") == 0 {
//		u, err := url.Parse(Config.OSS.Minio.Endpoint)
//		if err != nil {
//			return fmt.Errorf("minio endpoint url parse error: %w", err)
//		}
//		if u.Port() == "" {
//			Config.OSS.Minio.Endpoint = net.JoinHostPort(u.Host, "10005")
//		} else {
//			Config.OSS.Minio.Endpoint = u.Host
//		}
//	}
//
//	//fmt.Printf("remote minio openIMConfig: %v\n", configEtcd.Config.Credential.Minio)
//	//fmt.Println()
//	//fmt.Printf("local minio openIMConfig: %v\n", Config.OSS.Minio)
//
//	return nil
//}
