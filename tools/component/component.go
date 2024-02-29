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
	"time"

	"github.com/OpenIMSDK/chat/pkg/common/config"
	component "github.com/OpenIMSDK/tools/component"
)

var (
	MaxConnectTimes = 200
)

type checkFunc struct {
	name     string
	function func() (string, error)
}

func ComponentCheck() error {
	var err error
	var strInfo string
	if config.Config.Envs.Discovery != "k8s" {
		checks := []checkFunc{
			{name: "Zookeeper", function: checkZookeeper},
			{name: "Redis", function: checkRedis},
			{name: "MySQL", function: checkMySQL},
		}

		for i := 0; i < component.MaxRetry; i++ {
			if i != 0 {
				time.Sleep(1 * time.Second)
			}
			fmt.Printf("Checking components Round %v...\n", i+1)

			allSuccess := true
			for _, check := range checks {
				strInfo, err = check.function()
				if err != nil {
					component.ErrorPrint(fmt.Sprintf("Starting %s failed, %v", check.name, err))
					allSuccess = false
					break
				} else {
					component.SuccessPrint(fmt.Sprintf("%s connected successfully, %s", check.name, strInfo))
				}
			}

			if allSuccess {
				component.SuccessPrint("All components started successfully!")
				return nil
			}
		}
	}
	return err
}

// checkZookeeper checks the Zookeeper connection
func checkZookeeper() (string, error) {
	// Prioritize environment variables
	zk := &component.Zookeeper{
		Schema:   config.Config.Zookeeper.Schema,
		ZkAddr:   config.Config.Zookeeper.ZkAddr,
		Username: config.Config.Zookeeper.Username,
		Password: config.Config.Zookeeper.Password,
	}

	err := component.CheckZookeeper(zk)
	if err != nil {
		return "", err
	}
	return "", nil
}

// checkRedis checks the Redis connection
func checkRedis() (string, error) {
	redis := &component.Redis{
		Address:  *config.Config.Redis.Address,
		Username: config.Config.Redis.Username,
		Password: config.Config.Redis.Password,
	}

	err := component.CheckRedis(redis)
	if err != nil {
		return "", err
	}
	return "", nil
}

func checkMySQL() (string, error) {

	mysql := &component.MySQL{
		Address:  *config.Config.Mysql.Address,
		Username: *config.Config.Mysql.Username,
		Password: *config.Config.Mysql.Password,
		Database: *config.Config.Mysql.Database,
	}
	err := component.CheckMySQL(mysql)
	if err != nil {
		return "", err
	}
	return "", nil
}
