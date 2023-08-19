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

package config

import _ "embed"

//go:embed version
var Version string

var Config struct {
	Zookeeper struct {
		Schema   string   `yaml:"schema"`
		ZkAddr   []string `yaml:"zkAddr"`
		Username string   `yaml:"username"`
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
	Rpc struct {
		RegisterIP string `yaml:"registerIP"`
		ListenIP   string `yaml:"listenIP"`
	} `yaml:"rpc"`
	Redis struct {
		Address  *[]string `yaml:"address"`
		Username *string   `yaml:"username"`
		Password *string   `yaml:"password"`
	} `yaml:"redis"`
	RpcPort struct {
		OpenImAdminPort []int `yaml:"openImAdminPort"`
		OpenImChatPort  []int `yaml:"openImChatPort"`
	} `yaml:"rpcPort"`
	RpcRegisterName struct {
		OpenImAdminName string `yaml:"openImAdminName"`
		OpenImChatName  string `yaml:"openImChatName"`
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
		RotationTime        *uint   `yaml:"rotationTime"`
		RemainRotationCount *uint   `yaml:"remainRotationCount"`
		RemainLogLevel      *int    `yaml:"remainLogLevel"`
		IsStdout            *bool   `yaml:"isStdout"`
		IsJson              *bool   `yaml:"isJson"`
		WithStack           *bool   `yaml:"withStack"`
	} `yaml:"log"`
	Secret      *string `yaml:"secret"`
	OpenIMUrl   string  `yaml:"openIMUrl"`
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
	ProxyHeader string  `yaml:"proxyHeader"`
	AdminList   []Admin `yaml:"adminList"`
}

type Admin struct {
	AdminID   string `yaml:"adminID"`
	NickName  string `yaml:"nickname"`
	ImAdminID string `yaml:"imAdmin"`
}
