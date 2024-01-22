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

package chatrpcstart

import (
	"fmt"
	"net"
	"strconv"

	"github.com/OpenIMSDK/chat/pkg/common/config"
	chatMw "github.com/OpenIMSDK/chat/pkg/common/mw"
	"github.com/OpenIMSDK/chat/pkg/discovery_register"
	"github.com/OpenIMSDK/tools/discoveryregistry"
	"github.com/OpenIMSDK/tools/errs"
	"github.com/OpenIMSDK/tools/mw"
	"github.com/OpenIMSDK/tools/network"
	"github.com/OpenIMSDK/tools/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Start(rpcPort int, rpcRegisterName string, prometheusPort int, rpcFn func(client discoveryregistry.SvcDiscoveryRegistry, server *grpc.Server) error, options ...grpc.ServerOption) error {
	fmt.Println("start", rpcRegisterName, "server, port: ", rpcPort, "prometheusPort:", prometheusPort, ", OpenIM version: ", config.Version)

	var zkClient discoveryregistry.SvcDiscoveryRegistry
	zkClient, err := discovery_register.NewDiscoveryRegister(config.Config.Envs.Discovery)
	/*
		zkClient, err := openKeeper.NewClient(config.Config.Zookeeper.ZkAddr, config.Config.Zookeeper.Schema,
			openKeeper.WithFreq(time.Hour), openKeeper.WithUserNameAndPassword(config.Config.Zookeeper.Username,
				config.Config.Zookeeper.Password), openKeeper.WithRoundRobin(), openKeeper.WithTimeout(10), openKeeper.WithLogger(log.NewZkLogger()))*/if err != nil {
		return errs.Wrap(err, fmt.Sprintf(";the addr is:%v", &config.Config.Zookeeper.ZkAddr))
	}
	// defer zkClient.CloseZK()
	defer zkClient.Close()
	zkClient.AddOption(chatMw.AddUserType(), mw.GrpcClient(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	registerIP, err := network.GetRpcRegisterIP(config.Config.Rpc.RegisterIP)
	if err != nil {
		return utils.Wrap1(err)
	}
	srv := grpc.NewServer(append(options, mw.GrpcServer())...)
	defer srv.GracefulStop()
	err = rpcFn(zkClient, srv)
	if err != nil {
		return err
	}
	err = zkClient.Register(rpcRegisterName, registerIP, rpcPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return utils.Wrap1(err)
	}
	listener, err := net.Listen("tcp", net.JoinHostPort(network.GetListenIP(config.Config.Rpc.ListenIP), strconv.Itoa(rpcPort)))
	if err != nil {
		return utils.Wrap1(err)
	}
	defer listener.Close()
	return utils.Wrap1(srv.Serve(listener))
}
