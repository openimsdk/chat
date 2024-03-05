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
	"fmt"
	"github.com/OpenIMSDK/tools/errs"
	"strings"
	"time"

	"github.com/OpenIMSDK/chat/pkg/common/config"
	component "github.com/OpenIMSDK/tools/component"
)

var (
	MaxConnectTimes = 100
)

type checkFunc struct {
	name     string
	function func() error
	flag     bool
}

func ComponentCheck() error {
	if config.Config.Envs.Discovery != "k8s" {
		checks := []checkFunc{
			{name: "Zookeeper", function: checkZookeeper},
			{name: "Redis", function: checkRedis},
			//	{name: "Mongo", function: checkMongo, config: conf},
		}

		for i := 0; i < MaxConnectTimes; i++ {
			if i != 0 {
				time.Sleep(1 * time.Second)
			}
			fmt.Printf("Checking components Round %v...\n", i+1)

			var err error
			allSuccess := true
			for index, check := range checks {
				if !check.flag {
					err = check.function()
					if err != nil {
						allSuccess = false
						component.ErrorPrint(fmt.Sprintf("Starting %s failed:%v.", check.name, errs.Unwrap(err).Error()))
						if !strings.Contains(errs.Unwrap(err).Error(), "connection refused") &&
							!strings.Contains(errs.Unwrap(err).Error(), "timeout") {
							return err
						}
					} else {
						checks[index].flag = true
						component.SuccessPrint(fmt.Sprintf("%s connected successfully", check.name))
					}
				}
			}

			if allSuccess {
				component.SuccessPrint("All components started successfully!")
				return nil
			}
		}
	}
	return errs.Wrap(fmt.Errorf("components started failed"))
}

// checkZookeeper checks the Zookeeper connection
func checkZookeeper() error {
	zkStu := &component.Zookeeper{
		Schema:   config.Config.Zookeeper.Schema,
		ZkAddr:   config.Config.Zookeeper.ZkAddr,
		Username: config.Config.Zookeeper.Username,
		Password: config.Config.Zookeeper.Password,
	}
	err := component.CheckZookeeper(zkStu)
	return err
}

// checkRedis checks the Redis connection
func checkRedis() error {
	redisStu := &component.Redis{
		Address:  *config.Config.Redis.Address,
		Username: config.Config.Redis.Username,
		Password: config.Config.Redis.Password,
	}
	err := component.CheckRedis(redisStu)
	return err
}

// checkMongo checks the MongoDB connection without retries
//func checkMongo(config *config.GlobalConfig) error {
//	mongoStu := &component.Mongo{
//		URL:         config.Mongo.Uri,
//		Address:     config.Mongo.Address,
//		Database:    config.Mongo.Database,
//		Username:    config.Mongo.Username,
//		Password:    config.Mongo.Password,
//		MaxPoolSize: config.Mongo.MaxPoolSize,
//	}
//	err := component.CheckMongo(mongoStu)
//
//	return err
//}
