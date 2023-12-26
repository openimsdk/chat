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

package component

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"os"
	"strings"
	"time"

	"github.com/OpenIMSDK/protocol/constant"

	"github.com/OpenIMSDK/chat/pkg/common/config"
	"github.com/OpenIMSDK/tools/errs"
	"github.com/OpenIMSDK/tools/log"
	"github.com/go-zookeeper/zk"
	"github.com/pkg/errors"
)

var (
	MaxConnectTimes = 200
)

func ComponentCheck(cfgPath string, hide bool) error {
	if config.Config.Envs.Discovery != "k8s" {
		if _, err := checkNewZkClient(hide); err != nil {
			errorPrint(fmt.Sprintf("%v.Please check if your openIM server has started", err.Error()), hide)
			return err
		}
		// if err = checkGetCfg(zkConn, hide); err != nil {
		// 	errorPrint(fmt.Sprintf("%v.Please check if your openIM server has started", err.Error()), hide)
		// 	return err
		// }
	}
	//_, err := checkRedis()

	return nil
}

func errorPrint(s string, hide bool) {
	if !hide {
		fmt.Printf("\x1b[%dm%v\x1b[0m\n", 31, s)
	}
}

func successPrint(s string, hide bool) {
	if !hide {
		fmt.Printf("\x1b[%dm%v\x1b[0m\n", 32, s)
	}
}

func newZkClient() (*zk.Conn, error) {
	var c *zk.Conn
	var err error
	c, eventChan, err := zk.Connect(config.Config.Zookeeper.ZkAddr, time.Second*5, zk.WithLogger(log.NewZkLogger()))
	if err != nil {
		fmt.Println("zookeeper connect error:", err)
		return nil, errs.Wrap(err, "Zookeeper Addr: "+strings.Join(config.Config.Zookeeper.ZkAddr, " "))
	}

	// wait for successfully connect
	timeout := time.After(5 * time.Second)
	for {
		select {
		case event := <-eventChan:
			if event.State == zk.StateConnected {
				fmt.Println("Connected to Zookeeper")
				goto Connected
			}
		case <-timeout:
			return nil, errs.Wrap(errors.New("timeout waiting for Zookeeper connection"), "Zookeeper Addr: "+strings.Join(config.Config.Zookeeper.ZkAddr, " "))
		}
	}
Connected:

	if config.Config.Zookeeper.Username != "" && config.Config.Zookeeper.Password != "" {
		if err := c.AddAuth("digest", []byte(config.Config.Zookeeper.Username+":"+config.Config.Zookeeper.Password)); err != nil {
			return nil, errs.Wrap(err, "Zookeeper Username: "+config.Config.Zookeeper.Username+
				", Zookeeper Password: "+config.Config.Zookeeper.Password+
				", Zookeeper Addr: "+strings.Join(config.Config.Zookeeper.ZkAddr, " "))
		}
	}

	return c, nil
}

func checkNewZkClient(hide bool) (*zk.Conn, error) {
	for i := 0; i < MaxConnectTimes; i++ {
		if i != 0 {
			time.Sleep(3 * time.Second)
		}
		zkConn, err := newZkClient()
		if err != nil {
			if zkConn != nil {
				zkConn.Close()
			}
			errorPrint(fmt.Sprintf("Starting Zookeeper failed: %v.Please make sure your Zookeeper service has started", err.Error()), hide)
			continue
		}
		successPrint(fmt.Sprintf("zk starts successfully after: %v times ", i+1), hide)
		return zkConn, nil
	}
	return nil, errs.Wrap(errors.New("Connecting to zk fails"))
}

func checkGetCfg(conn *zk.Conn, hide bool) error {
	for i := 0; i < MaxConnectTimes; i++ {
		if i != 0 {
			time.Sleep(3 * time.Second)
		}
		path := "/" + config.Config.Zookeeper.Schema + "/" + constant.OpenIMCommonConfigKey

		zkConfig, _, err := conn.Get(path)
		if err != nil {
			fmt.Println("path =", path, "zkConfig is:", zkConfig)
			errorPrint(fmt.Sprintf("! get zk config [%d] error: %v\n", i, err), hide)
			continue
		} else if len(zkConfig) == 0 {
			errorPrint(fmt.Sprintf("! get zk config [%d] data is empty\n", i), hide)
			continue
		}
		successPrint(fmt.Sprint("Chat get config successfully"), hide)
		return nil
	}
	return errors.New("Getting config from zk failed")
}

// checkRedis checks the Redis connection
func checkRedis() (string, error) {
	// Prioritize environment variables
	address := getEnv("REDIS_ADDRESS", strings.Join(*config.Config.Redis.Address, ","))
	username := getEnv("REDIS_USERNAME", config.Config.Redis.Username)
	password := getEnv("REDIS_PASSWORD", config.Config.Redis.Password)

	// Split address to handle multiple addresses for cluster setup
	redisAddresses := strings.Split(address, ",")

	var redisClient redis.UniversalClient
	if len(redisAddresses) > 1 {
		// Use cluster client for multiple addresses
		redisClient = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    redisAddresses,
			Username: username,
			Password: password,
		})
	} else {
		// Use regular client for single address
		redisClient = redis.NewClient(&redis.Options{
			Addr:     redisAddresses[0],
			Username: username,
			Password: password,
		})
	}
	defer redisClient.Close()

	// Ping Redis to check connectivity
	_, err := redisClient.Ping(context.Background()).Result()
	str := "the addr is:" + strings.Join(redisAddresses, ",")
	if err != nil {
		return "", errs.Wrap(err, str)
	}

	return str, nil
}
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
