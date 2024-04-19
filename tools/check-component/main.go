// Copyright © 2023 OpenIM. All rights reserved.
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
	"context"
	"flag"
	"fmt"
	"github.com/openimsdk/chat/pkg/common/cmd"
	"github.com/openimsdk/chat/pkg/common/config"
	"github.com/openimsdk/chat/pkg/common/imapi"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/db/redisutil"
	"github.com/openimsdk/tools/discovery/zookeeper"
	"github.com/openimsdk/tools/mcontext"
	"github.com/openimsdk/tools/system/program"
	"github.com/openimsdk/tools/utils/idutil"
	"path/filepath"
	"time"
)

const maxRetry = 180

func CheckZookeeper(ctx context.Context, config *config.ZooKeeper) error {
	return zookeeper.Check(ctx, config.Address, config.Schema, zookeeper.WithUserNameAndPassword(config.Username, config.Password))
}

func CheckMongo(ctx context.Context, config *config.Mongo) error {
	return mongoutil.Check(ctx, config.Build())
}

func CheckRedis(ctx context.Context, config *config.Redis) error {
	return redisutil.Check(ctx, config.Build())
}

func CheckOpenIM(ctx context.Context, apiURL, secret, adminUserID string) error {
	api2 := imapi.New(apiURL, secret, adminUserID)
	_, err := api2.UserToken(mcontext.SetOperationID(ctx, "CheckOpenIM"+idutil.OperationIDGenerator()), adminUserID, constant.AdminPlatformID)
	return err
}

func initConfig(configDir string) (*config.Mongo, *config.Redis, *config.ZooKeeper, *config.Share, error) {
	var (
		mongoConfig     = &config.Mongo{}
		redisConfig     = &config.Redis{}
		zookeeperConfig = &config.ZooKeeper{}
		shareConfig     = &config.Share{}
	)
	err := config.LoadConfig(filepath.Join(configDir, cmd.MongodbConfigFileName), cmd.ConfigEnvPrefixMap[cmd.MongodbConfigFileName], mongoConfig)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	err = config.LoadConfig(filepath.Join(configDir, cmd.RedisConfigFileName), cmd.ConfigEnvPrefixMap[cmd.RedisConfigFileName], redisConfig)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	err = config.LoadConfig(filepath.Join(configDir, cmd.ZookeeperConfigFileName), cmd.ConfigEnvPrefixMap[cmd.ZookeeperConfigFileName], zookeeperConfig)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	err = config.LoadConfig(filepath.Join(configDir, cmd.ShareFileName), cmd.ConfigEnvPrefixMap[cmd.ShareFileName], shareConfig)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	return mongoConfig, redisConfig, zookeeperConfig, shareConfig, nil
}

func main() {
	var index int
	var configDir string
	flag.IntVar(&index, "i", 0, "Index number")
	defaultConfigDir := filepath.Join("..", "..", "..", "..", "..", "config")
	flag.StringVar(&configDir, "c", defaultConfigDir, "Configuration dir")
	flag.Parse()

	fmt.Printf("Index: %d, Config Path: %s\n", index, configDir)

	mongoConfig, redisConfig, zookeeperConfig, shareConfig, err := initConfig(configDir)
	if err != nil {
		program.ExitWithError(err)
	}

	ctx := context.Background()
	err = performChecks(ctx, mongoConfig, redisConfig, zookeeperConfig, shareConfig, maxRetry)
	if err != nil {
		// Assume program.ExitWithError logs the error and exits.
		// Replace with your error handling logic as necessary.
		program.ExitWithError(err)
	}
}

func performChecks(ctx context.Context, mongoConfig *config.Mongo, redisConfig *config.Redis, zookeeperConfig *config.ZooKeeper, shareConfig *config.Share, maxRetry int) error {
	checksDone := make(map[string]bool)

	checks := map[string]func(ctx context.Context) error{
		"Zookeeper": func(ctx context.Context) error {
			return CheckZookeeper(ctx, zookeeperConfig)
		},
		"Mongo": func(ctx context.Context) error {
			return CheckMongo(ctx, mongoConfig)
		},
		"Redis": func(ctx context.Context) error {
			return CheckRedis(ctx, redisConfig)
		},
		"OpenIM": func(ctx context.Context) error {
			return CheckOpenIM(ctx, shareConfig.OpenIM.ApiURL, shareConfig.OpenIM.Secret, shareConfig.OpenIM.AdminUserID)
		},
	}

	for i := 0; i < maxRetry; i++ {
		allSuccess := true
		for name, check := range checks {
			ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			if !checksDone[name] {
				if err := check(ctx); err != nil {
					fmt.Printf("%s check failed: %v\n", name, err)
					allSuccess = false
				} else {
					fmt.Printf("%s check succeeded.\n", name)
					checksDone[name] = true
				}
			}
			cancel()
		}

		if allSuccess {
			fmt.Println("All components checks passed successfully.")
			return nil
		}

		time.Sleep(1 * time.Second)
	}

	return fmt.Errorf("not all components checks passed successfully after %d attempts", maxRetry)
}
