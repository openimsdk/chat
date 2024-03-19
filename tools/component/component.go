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

package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/OpenIMSDK/tools/errs"

	"github.com/OpenIMSDK/tools/errs"

	"github.com/OpenIMSDK/chat/pkg/common/config"
	component "github.com/OpenIMSDK/tools/component"
)

const (
	// defaultCfgPath is the default path of the configuration file.
	defaultCfgPath  = "../../../../../config/config.yaml"
	MaxConnectTimes = 100
)

var cfgPath = flag.String("config_folder_path", defaultCfgPath, "Path to the configuration file")

type checkFunc struct {
	name     string
	function func() error
	flag     bool
}

func main() {
	flag.Parse()
	if err := config.InitConfig(*cfgPath); err != nil {
		fmt.Printf("Read config failed: %v\n", err)
		return
	}

	if config.Config.Envs.Discovery != "k8s" {
		checks := []checkFunc{
			{name: "Zookeeper", function: checkZookeeper},
			{name: "Redis", function: checkRedis},
			{name: "Mongo", function: checkMongo},
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
							os.Exit(-1)
						}
					} else {
						checks[index].flag = true
						component.SuccessPrint(fmt.Sprintf("%s connected successfully", check.name))
					}
				}
			}

			if allSuccess {
				component.SuccessPrint("All components started successfully!")
				return
			}
		}
	}
	component.ErrorPrint("Some components started failed!")
	os.Exit(-1)
}

// checkZookeeper checks the Zookeeper connection.
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

// checkRedis checks the Redis connection.
func checkRedis() error {
	redisStu := &component.Redis{
		Address:  *config.Config.Redis.Address,
		Username: config.Config.Redis.Username,
		Password: config.Config.Redis.Password,
	}
	err := component.CheckRedis(redisStu)
	return err
}

// checkMongo checks the MongoDB connection without retries.
func checkMongo() error {
	mongoStu := &component.Mongo{
		URL:         config.Config.Mongo.Uri,
		Address:     config.Config.Mongo.Address,
		Database:    config.Config.Mongo.Database,
		Username:    config.Config.Mongo.Username,
		Password:    config.Config.Mongo.Password,
		MaxPoolSize: config.Config.Mongo.MaxPoolSize,
	}
	err := component.CheckMongo(mongoStu)

	return err
}
