package admin

import (
	"encoding/json"
	"reflect"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/openimsdk/chat/pkg/common/apistruct"
	"github.com/openimsdk/chat/pkg/common/config"
	"github.com/openimsdk/chat/pkg/common/kdisc"
	"github.com/openimsdk/chat/pkg/common/kdisc/etcd"
	"github.com/openimsdk/chat/version"
	"github.com/openimsdk/tools/apiresp"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/utils/runtimeenv"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type ConfigManager struct {
	config     *config.AllConfig
	client     *clientv3.Client
	configPath string
	runtimeEnv string
}

func NewConfigManager(cfg *config.AllConfig, client *clientv3.Client, configPath string, runtimeEnv string) *ConfigManager {
	return &ConfigManager{
		config:     cfg,
		client:     client,
		configPath: configPath,
		runtimeEnv: runtimeEnv,
	}
}

func (cm *ConfigManager) GetConfig(c *gin.Context) {
	var req apistruct.GetConfigReq
	if err := c.BindJSON(&req); err != nil {
		apiresp.GinError(c, errs.ErrArgs.WithDetail(err.Error()).Wrap())
		return
	}
	conf := cm.config.Name2Config(req.ConfigName)
	if conf == nil {
		apiresp.GinError(c, errs.ErrArgs.WithDetail("config name not found").Wrap())
		return
	}
	b, err := json.Marshal(conf)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiresp.GinSuccess(c, string(b))
}

func (cm *ConfigManager) GetConfigList(c *gin.Context) {
	var resp apistruct.GetConfigListResp
	resp.ConfigNames = cm.config.GetConfigNames()
	resp.Environment = runtimeenv.PrintRuntimeEnvironment()
	resp.Version = version.Version

	apiresp.GinSuccess(c, resp)
}

func (cm *ConfigManager) SetConfig(c *gin.Context) {
	if cm.config.Discovery.Enable != kdisc.ETCDCONST {
		apiresp.GinError(c, errs.New("only etcd support set config").Wrap())
		return
	}
	var req apistruct.SetConfigReq
	if err := c.BindJSON(&req); err != nil {
		apiresp.GinError(c, errs.ErrArgs.WithDetail(err.Error()).Wrap())
		return
	}
	var err error
	switch req.ConfigName {
	case config.DiscoveryConfigFileName:
		err = compareAndSave[config.Discovery](c, cm.config.Name2Config(req.ConfigName), &req, cm.client)
	case config.LogConfigFileName:
		err = compareAndSave[config.Log](c, cm.config.Name2Config(req.ConfigName), &req, cm.client)
	case config.MongodbConfigFileName:
		err = compareAndSave[config.Mongo](c, cm.config.Name2Config(req.ConfigName), &req, cm.client)
	case config.ChatAPIAdminCfgFileName:
		err = compareAndSave[config.API](c, cm.config.Name2Config(req.ConfigName), &req, cm.client)
	case config.ChatAPIChatCfgFileName:
		err = compareAndSave[config.API](c, cm.config.Name2Config(req.ConfigName), &req, cm.client)
	case config.ChatRPCAdminCfgFileName:
		err = compareAndSave[config.Admin](c, cm.config.Name2Config(req.ConfigName), &req, cm.client)
	case config.ChatRPCChatCfgFileName:
		err = compareAndSave[config.Chat](c, cm.config.Name2Config(req.ConfigName), &req, cm.client)
	case config.ShareFileName:
		err = compareAndSave[config.Share](c, cm.config.Name2Config(req.ConfigName), &req, cm.client)
	case config.RedisConfigFileName:
		err = compareAndSave[config.Redis](c, cm.config.Name2Config(req.ConfigName), &req, cm.client)
	default:
		apiresp.GinError(c, errs.ErrArgs.Wrap())
		return
	}
	if err != nil {
		apiresp.GinError(c, errs.ErrArgs.WithDetail(err.Error()).Wrap())
		return
	}
	apiresp.GinSuccess(c, nil)
}

func compareAndSave[T any](c *gin.Context, old any, req *apistruct.SetConfigReq, client *clientv3.Client) error {
	conf := new(T)
	err := json.Unmarshal([]byte(req.Data), &conf)
	if err != nil {
		return errs.ErrArgs.WithDetail(err.Error()).Wrap()
	}
	eq := reflect.DeepEqual(old, conf)
	if eq {
		return nil
	}
	data, err := json.Marshal(conf)
	if err != nil {
		return errs.ErrArgs.WithDetail(err.Error()).Wrap()
	}
	_, err = client.Put(c, etcd.BuildKey(req.ConfigName), string(data))
	if err != nil {
		return errs.WrapMsg(err, "save to etcd failed")
	}
	return nil
}

func (cm *ConfigManager) ResetConfig(c *gin.Context) {
	go cm.resetConfig(c)
	apiresp.GinSuccess(c, nil)
}

func (cm *ConfigManager) resetConfig(c *gin.Context) {
	txn := cm.client.Txn(c)
	type initConf struct {
		old       any
		new       any
		isChanged bool
	}
	configMap := map[string]*initConf{
		config.DiscoveryConfigFileName: {old: &cm.config.Discovery, new: new(config.Discovery)},
		config.LogConfigFileName:       {old: &cm.config.Log, new: new(config.Log)},
		config.MongodbConfigFileName:   {old: &cm.config.Mongo, new: new(config.Mongo)},
		config.ChatAPIAdminCfgFileName: {old: &cm.config.AdminAPI, new: new(config.API)},
		config.ChatAPIChatCfgFileName:  {old: &cm.config.ChatAPI, new: new(config.API)},
		config.ChatRPCAdminCfgFileName: {old: &cm.config.Admin, new: new(config.Admin)},
		config.ChatRPCChatCfgFileName:  {old: &cm.config.Chat, new: new(config.Chat)},
		config.RedisConfigFileName:     {old: &cm.config.Redis, new: new(config.Redis)},
		config.ShareFileName:           {old: &cm.config.Share, new: new(config.Share)},
	}

	changedKeys := make([]string, 0, len(configMap))
	for k, v := range configMap {
		err := config.Load(
			cm.configPath,
			k,
			config.EnvPrefixMap[k],
			cm.runtimeEnv,
			v.new,
		)
		if err != nil {
			log.ZError(c, "load config failed", err)
			continue
		}
		v.isChanged = reflect.DeepEqual(v.old, v.new)
		if !v.isChanged {
			changedKeys = append(changedKeys, k)
		}
	}

	ops := make([]clientv3.Op, 0)
	for _, k := range changedKeys {
		data, err := json.Marshal(configMap[k].new)
		if err != nil {
			log.ZError(c, "marshal config failed", err)
			continue
		}
		ops = append(ops, clientv3.OpPut(etcd.BuildKey(k), string(data)))
	}
	if len(ops) > 0 {
		txn.Then(ops...)
		_, err := txn.Commit()
		if err != nil {
			log.ZError(c, "commit etcd txn failed", err)
			return
		}
	}
}

func (cm *ConfigManager) Restart(c *gin.Context) {
	go cm.restart(c)
	apiresp.GinSuccess(c, nil)
}

func (cm *ConfigManager) restart(c *gin.Context) {
	time.Sleep(time.Millisecond * 200) // wait for Restart http call return
	t := time.Now().Unix()
	_, err := cm.client.Put(c, etcd.BuildKey(etcd.RestartKey), strconv.Itoa(int(t)))
	if err != nil {
		log.ZError(c, "restart etcd put key failed", err)
	}
}
