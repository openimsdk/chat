package config

import _ "embed"

//go:embed version
var Version string

var Config struct {
	Zookeeper struct {
		Schema   string   `yaml:"schema"`
		ZkAddr   []string `yaml:"zkAddr"`
		UserName string   `yaml:"userName"`
		Password string   `yaml:"password"`
	} `yaml:"zookeeper"`
	ChatApi struct {
		GinPort  []int  `yaml:"openImChatApiPort"`
		ListenIP string `yaml:"listenIP"`
	}
	AdminApi struct {
		GinPort  []int  `yaml:"openImAdminApiPort"`
		ListenIP string `yaml:"listenIP"`
	}
	RpcPort struct {
		OpenImAdminPort []int `yaml:"openImAdminPort"`
		OpenImChatPort  []int `yaml:"openImChatPort"`
	} `yaml:"rpcport"`
	RpcRegisterName struct {
		OpenImAdminName string `yaml:"openImAdminName"`
		OpenImChatName  string `yaml:"openImChatName"`
	} `yaml:"rpcregistername"`
	Mysql struct {
		DBAddress      []string `yaml:"dbMysqlAddress"`
		DBUserName     string   `yaml:"dbMysqlUserName"`
		DBPassword     string   `yaml:"dbMysqlPassword"`
		DBDatabaseName string   `yaml:"dbMysqlDatabaseName"`
		DBMaxOpenConns int      `yaml:"dbMaxOpenConns"`
		DBMaxIdleConns int      `yaml:"dbMaxIdleConns"`
		DBMaxLifeTime  int      `yaml:"dbMaxLifeTime"`
		LogLevel       int      `yaml:"logLevel"`
		SlowThreshold  int      `yaml:"slowThreshold"`
	}
	Log struct {
		StorageLocation     string `yaml:"storageLocation"`
		RotationTime        int    `yaml:"rotationTime"`
		RemainRotationCount uint   `yaml:"remainRotationCount"`
		RemainLogLevel      int    `yaml:"remainLogLevel"`
		IsStdout            bool   `yaml:"isStdout"`
		WithStack           bool   `yaml:"withStack"`
		IsJson              bool   `yaml:"isJson"`
	}
	TokenPolicy struct {
		AccessSecret string `yaml:"accessSecret"`
		AccessExpire int64  `yaml:"accessExpire"`
	} `yaml:"tokenPolicy"`
	VerifyCode struct {
		ValidTime int    `yaml:"validTime"`
		UintTime  int    `yaml:"uintTime"`
		MaxCount  int    `yaml:"maxCount"`
		SuperCode string `yaml:"superCode"`
		Len       int    `yaml:"len"`
		Use       string `yaml:"use"`
		Ali       struct {
			Endpoint                     string `yaml:"endpoint"`
			AccessKeyId                  string `yaml:"accessKeyId"`
			AccessKeySecret              string `yaml:"accessKeySecret"`
			SignName                     string `yaml:"signName"`
			VerificationCodeTemplateCode string `yaml:"verificationCodeTemplateCode"`
		}
	}
	ProxyHeader string `yaml:"proxyheader"`
}
