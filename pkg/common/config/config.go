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
	} `yaml:"chatApi"`
	AdminApi struct {
		GinPort  []int  `yaml:"openImAdminApiPort"`
		ListenIP string `yaml:"listenIP"`
	} `yaml:"adminApi"`
	RpcPort struct {
		OpenImAdminPort        []int `yaml:"openImAdminPort"`
		OpenImChatPort         []int `yaml:"openImChatPort"`
		OpenImOrganizationPort []int `yaml:"openImOrganizationPort"`
	} `yaml:"rpcPort"`
	RpcRegisterName struct {
		OpenImAdminName string `yaml:"openImAdminName"`
		OpenImChatName  string `yaml:"openImChatName"`
		OpenImOrganizationName string `yaml:"openImOrganizationName"
	} `yaml:"rpcRegisterName"`
	Mysql struct {
		Address       *[]string `yaml:"address"`
		Username      *string   `yaml:"username"`
		Password      *string   `yaml:"password"`
		Database      *string   `yaml:"database"`
		MaxOpenConn   *int      `yaml:"maxOpenConn"`
		MaxIdleConn   *int      `yaml:"maxIdleConn"`
		MaxLifeTime   *int      `yaml:"maxLifeTime"`
		LogLevel      *int      `yaml:"logLevel"`
		SlowThreshold *int      `yaml:"slowThreshold"`
	} `yaml:"mysql"`
	Log struct {
		StorageLocation     *string `yaml:"storageLocation"`
		RotationTime        *int    `yaml:"rotationTime"`
		RemainRotationCount *uint   `yaml:"remainRotationCount"`
		RemainLogLevel      *int    `yaml:"remainLogLevel"`
		IsStdout            *bool   `yaml:"isStdout"`
		IsJson              *bool   `yaml:"isJson"`
		WithStack           *bool   `yaml:"withStack"`
	} `yaml:"log"`
	Secret      *string `yaml:"secret"`
	TokenPolicy struct {
		Expire *int64 `yaml:"expire"`
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
		} `yaml:"ali"`
	} `yaml:"verifyCode"`
	ProxyHeader string `yaml:"proxyHeader"`
}
