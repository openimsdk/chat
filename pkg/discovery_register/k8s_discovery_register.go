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

package discovery_register

import (
	"context"
	"errors"
	"fmt"
	"github.com/OpenIMSDK/tools/errs"
	"strings"
	"time"

	"github.com/OpenIMSDK/chat/pkg/common/config"
	"github.com/OpenIMSDK/tools/discoveryregistry"
	openkeeper "github.com/OpenIMSDK/tools/discoveryregistry/zookeeper"
	"github.com/OpenIMSDK/tools/log"
	"google.golang.org/grpc"
)

func NewDiscoveryRegister(envType string) (discoveryregistry.SvcDiscoveryRegistry, error) {
	var client discoveryregistry.SvcDiscoveryRegistry
	var err error
	switch envType {
	case "zookeeper":
		client, err = openkeeper.NewClient(config.Config.Zookeeper.ZkAddr, config.Config.Zookeeper.Schema,
			openkeeper.WithFreq(time.Hour), openkeeper.WithUserNameAndPassword(
				config.Config.Zookeeper.Username,
				config.Config.Zookeeper.Password,
			), openkeeper.WithRoundRobin(), openkeeper.WithTimeout(10), openkeeper.WithLogger(log.NewZkLogger()))
		err = errs.Wrap(err,
			"Zookeeper ZkAddr: "+strings.Join(config.Config.Zookeeper.ZkAddr, ",")+
				", Zookeeper Schema: "+config.Config.Zookeeper.Schema+
				", Zookeeper Username: "+config.Config.Zookeeper.Username+
				", Zookeeper Password: "+config.Config.Zookeeper.Password)
	case "k8s":
		client, err = NewK8sDiscoveryRegister()
		err = errs.Wrap(err,
			"envType: "+"k8s")
	default:
		client = nil
		err = errs.Wrap(errors.New("envType not correct"))
	}
	return client, err
}

type K8sDR struct {
	options         []grpc.DialOption
	rpcRegisterAddr string
}

func (cli *K8sDR) GetUserIdHashGatewayHost(ctx context.Context, userId string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func NewK8sDiscoveryRegister() (discoveryregistry.SvcDiscoveryRegistry, error) {
	return &K8sDR{}, nil
}

func (cli *K8sDR) Register(serviceName, host string, port int, opts ...grpc.DialOption) error {
	cli.rpcRegisterAddr = serviceName
	return nil
}

func (cli *K8sDR) UnRegister() error {
	return nil
}

func (cli *K8sDR) CreateRpcRootNodes(serviceNames []string) error {
	return nil
}

func (cli *K8sDR) RegisterConf2Registry(key string, conf []byte) error {
	return nil
}

func (cli *K8sDR) GetConfFromRegistry(key string) ([]byte, error) {
	return nil, nil
}

func (cli *K8sDR) GetConns(ctx context.Context, serviceName string, opts ...grpc.DialOption) ([]*grpc.ClientConn, error) {
	conn, err := grpc.DialContext(ctx, serviceName, append(cli.options, opts...)...)
	return []*grpc.ClientConn{conn}, err
}

func (cli *K8sDR) GetConn(ctx context.Context, serviceName string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	return grpc.DialContext(ctx, serviceName, append(cli.options, opts...)...)
}

func (cli *K8sDR) GetSelfConnTarget() string {
	return cli.rpcRegisterAddr
}

func (cli *K8sDR) AddOption(opts ...grpc.DialOption) {
	cli.options = append(cli.options, opts...)
}

func (cli *K8sDR) CloseConn(conn *grpc.ClientConn) {
	conn.Close()
}

// do not use this method for call rpc.
func (cli *K8sDR) GetClientLocalConns() map[string][]*grpc.ClientConn {
	fmt.Println("should not call this function!!!!!!!!!!!!!!!!!!!!!!!!!")
	return nil
}

func (cli *K8sDR) Close() {
	return
}
