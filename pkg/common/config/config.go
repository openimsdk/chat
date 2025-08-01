package config

import (
	_ "embed"

	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/db/redisutil"
)

var (
	//go:embed version
	Version string
	//go:embed template.xlsx
	ImportTemplate []byte
)

type Share struct {
	OpenIM struct {
		ApiURL               string `mapstructure:"apiURL"`
		Secret               string `mapstructure:"secret"`
		AdminUserID          string `mapstructure:"adminUserID"`
		TokenRefreshInterval int    `mapstructure:"tokenRefreshInterval"`
	} `mapstructure:"openIM"`
	ChatAdmin   []string `mapstructure:"chatAdmin"`
	ProxyHeader string   `mapstructure:"proxyHeader"`
}

type RpcService struct {
	Chat  string `mapstructure:"chat"`
	Admin string `mapstructure:"admin"`
	Bot   string `mapstructure:"bot"`
}

func (r *RpcService) GetServiceNames() []string {
	return []string{
		r.Chat,
		r.Admin,
	}
}

type API struct {
	Api struct {
		ListenIP string `mapstructure:"listenIP"`
		Ports    []int  `mapstructure:"ports"`
	} `mapstructure:"api"`
}

type APIBot struct {
	Api struct {
		ListenIP string `mapstructure:"listenIP"`
		Ports    []int  `mapstructure:"ports"`
	} `mapstructure:"api"`
}

type Mongo struct {
	URI         string   `mapstructure:"uri"`
	Address     []string `mapstructure:"address"`
	Database    string   `mapstructure:"database"`
	Username    string   `mapstructure:"username"`
	Password    string   `mapstructure:"password"`
	AuthSource  string   `mapstructure:"authSource"`
	MaxPoolSize int      `mapstructure:"maxPoolSize"`
	MaxRetry    int      `mapstructure:"maxRetry"`
}

func (m *Mongo) Build() *mongoutil.Config {
	return &mongoutil.Config{
		Uri:         m.URI,
		Address:     m.Address,
		Database:    m.Database,
		Username:    m.Username,
		Password:    m.Password,
		AuthSource:  m.AuthSource,
		MaxPoolSize: m.MaxPoolSize,
		MaxRetry:    m.MaxRetry,
	}
}

type Redis struct {
	Address        []string `mapstructure:"address"`
	Username       string   `mapstructure:"username"`
	Password       string   `mapstructure:"password"`
	EnablePipeline bool     `mapstructure:"enablePipeline"`
	ClusterMode    bool     `mapstructure:"clusterMode"`
	DB             int      `mapstructure:"db"`
	MaxRetry       int      `mapstructure:"MaxRetry"`
}

func (r *Redis) Build() *redisutil.Config {
	return &redisutil.Config{
		ClusterMode: r.ClusterMode,
		Address:     r.Address,
		Username:    r.Username,
		Password:    r.Password,
		DB:          r.DB,
		MaxRetry:    r.MaxRetry,
	}
}

type Discovery struct {
	Enable     string     `mapstructure:"enable"`
	Etcd       Etcd       `mapstructure:"etcd"`
	Kubernetes Kubernetes `mapstructure:"kubernetes"`
	RpcService RpcService `mapstructure:"rpcService"`
}

type Kubernetes struct {
	Namespace string `mapstructure:"namespace"`
}

type Etcd struct {
	RootDirectory string   `mapstructure:"rootDirectory"`
	Address       []string `mapstructure:"address"`
	Username      string   `mapstructure:"username"`
	Password      string   `mapstructure:"password"`
}

type Chat struct {
	RPC struct {
		RegisterIP string `mapstructure:"registerIP"`
		ListenIP   string `mapstructure:"listenIP"`
		Ports      []int  `mapstructure:"ports"`
	} `mapstructure:"rpc"`
	VerifyCode VerifyCode `mapstructure:"verifyCode"`
	LiveKit    struct {
		URL    string `mapstructure:"url"`
		Key    string `mapstructure:"key"`
		Secret string `mapstructure:"secret"`
	} `mapstructure:"liveKit"`
	AllowRegister bool `mapstructure:"allowRegister"`
}

type Bot struct {
	RPC struct {
		RegisterIP string `mapstructure:"registerIP"`
		ListenIP   string `mapstructure:"listenIP"`
		Ports      []int  `mapstructure:"ports"`
	} `mapstructure:"rpc"`
	Timeout int `mapstructure:"timeout"`
}
type VerifyCode struct {
	ValidTime  int    `mapstructure:"validTime"`
	ValidCount int    `mapstructure:"validCount"`
	UintTime   int    `mapstructure:"uintTime"`
	MaxCount   int    `mapstructure:"maxCount"`
	SuperCode  string `mapstructure:"superCode"`
	Len        int    `mapstructure:"len"`
	Phone      struct {
		Use string `mapstructure:"use"`
		Ali struct {
			Endpoint                     string `mapstructure:"endpoint"`
			AccessKeyID                  string `mapstructure:"accessKeyId"`
			AccessKeySecret              string `mapstructure:"accessKeySecret"`
			SignName                     string `mapstructure:"signName"`
			VerificationCodeTemplateCode string `mapstructure:"verificationCodeTemplateCode"`
		} `mapstructure:"ali"`
	} `mapstructure:"phone"`
	Mail struct {
		Use                     string `mapstructure:"use"`
		Title                   string `mapstructure:"title"`
		SenderMail              string `mapstructure:"senderMail"`
		SenderAuthorizationCode string `mapstructure:"senderAuthorizationCode"`
		SMTPAddr                string `mapstructure:"smtpAddr"`
		SMTPPort                int    `mapstructure:"smtpPort"`
	} `mapstructure:"mail"`
}

type Admin struct {
	RPC struct {
		RegisterIP string `mapstructure:"registerIP"`
		ListenIP   string `mapstructure:"listenIP"`
		Ports      []int  `mapstructure:"ports"`
	} `mapstructure:"rpc"`
	TokenPolicy struct {
		Expire int `mapstructure:"expire"`
	} `mapstructure:"tokenPolicy"`
	Secret string `mapstructure:"secret"`
}

type Log struct {
	StorageLocation     string `mapstructure:"storageLocation"`
	RotationTime        uint   `mapstructure:"rotationTime"`
	RemainRotationCount uint   `mapstructure:"remainRotationCount"`
	RemainLogLevel      int    `mapstructure:"remainLogLevel"`
	IsStdout            bool   `mapstructure:"isStdout"`
	IsJson              bool   `mapstructure:"isJson"`
	IsSimplify          bool   `mapstructure:"isSimplify"`
	WithStack           bool   `mapstructure:"withStack"`
}

type AllConfig struct {
	AdminAPI  API
	ChatAPI   API
	Admin     Admin
	Chat      Chat
	Discovery Discovery
	Log       Log
	Mongo     Mongo
	Redis     Redis
	Share     Share
}

func (a *AllConfig) Name2Config(name string) any {
	switch name {
	case ChatAPIAdminCfgFileName:
		return a.AdminAPI
	case ChatAPIChatCfgFileName:
		return a.ChatAPI
	case ChatRPCAdminCfgFileName:
		return a.Admin
	case ChatRPCChatCfgFileName:
		return a.Chat
	case DiscoveryConfigFileName:
		return a.Discovery
	case LogConfigFileName:
		return a.Log
	case MongodbConfigFileName:
		return a.Mongo
	case RedisConfigFileName:
		return a.Redis
	case ShareFileName:
		return a.Share
	default:
		return nil
	}
}

func (a *AllConfig) GetConfigNames() []string {
	return []string{
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
}
