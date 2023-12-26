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

package zookeeper

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/OpenIMSDK/tools/errs"
	"github.com/OpenIMSDK/tools/log"

	"github.com/go-zookeeper/zk"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
)

var (
	ErrConnIsNil               = errors.New("conn is nil")
	ErrConnIsNilButLocalNotNil = errors.New("conn is nil, but local is not nil")
)

func (s *ZkClient) watch(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			s.logger.Printf("zk watch ctx done")
			return
		case event := <-s.eventChan:
			s.logger.Printf("zk eventChan recv new event: %+v", event)
			switch event.Type {
			case zk.EventSession:
				switch event.State {
				case zk.StateHasSession:
					if s.isRegistered && !s.isStateDisconnected {
						s.logger.Printf("zk session event stateHasSession: %+v, client prepare to create new temp node", event)
						node, err := s.CreateTempNode(s.rpcRegisterName, s.rpcRegisterAddr)
						if err != nil {
							s.logger.Printf("zk session event stateHasSession: %+v, create temp node error: %v", event, err)
						} else {
							s.node = node
						}
					}
				case zk.StateDisconnected:
					s.isStateDisconnected = true
				case zk.StateConnected:
					s.isStateDisconnected = false
				default:
					s.logger.Printf("zk session event: %+v", event)
				}
			case zk.EventNodeChildrenChanged:
				s.logger.Printf("zk event: %s", event.Path)
				l := strings.Split(event.Path, "/")
				if len(l) > 1 {
					serviceName := l[len(l)-1]
					s.lock.Lock()
					s.flushResolverAndDeleteLocal(serviceName)
					s.lock.Unlock()
				}
				s.logger.Printf("zk event handle success: %s", event.Path)
			case zk.EventNodeDataChanged:
			case zk.EventNodeCreated:
			case zk.EventNodeDeleted:
			case zk.EventNotWatching:
			}
		}
	}
}

func (s *ZkClient) GetConnsRemote(serviceName string) (conns []resolver.Address, err error) {
	path := s.getPath(serviceName)
	_, _, _, err = s.conn.ChildrenW(path)
	if err != nil {
		return nil, errors.Wrap(err, "children watch error")
	}
	childNodes, _, err := s.conn.Children(path)
	if err != nil {
		return nil, errors.Wrap(err, "get children error")
	} else {
		for _, child := range childNodes {
			fullPath := path + "/" + child
			data, _, err := s.conn.Get(fullPath)
			if err != nil {
				if err == zk.ErrNoNode {
					return nil, errors.Wrap(err, "this is zk ErrNoNode")
				}
				return nil, errors.Wrap(err, "get children error")
			}
			log.ZDebug(context.Background(), "get addrs from remote", "conn", string(data))
			conns = append(conns, resolver.Address{Addr: string(data), ServerName: serviceName})
		}
	}
	return conns, nil
}
func (s *ZkClient) GetUserIdHashGatewayHost(ctx context.Context, userId string) (string, error) {
	log.ZWarn(ctx, "not impliment", errors.New("zkclinet not impliment GetUserIdHashGatewayHost method"))
	return "", nil
}
func (s *ZkClient) GetConns(ctx context.Context, serviceName string, opts ...grpc.DialOption) ([]*grpc.ClientConn, error) {
	s.logger.Printf("get conns from client, serviceName: %s", serviceName)
	s.lock.Lock()
	defer s.lock.Unlock()
	conns := s.localConns[serviceName]
	if len(conns) == 0 {
		var err error
		s.logger.Printf("get conns from zk remote, serviceName: %s", serviceName)
		addrs, err := s.GetConnsRemote(serviceName)
		if err != nil {
			return nil, err
		}
		if len(addrs) == 0 {
			return nil, fmt.Errorf("no conn for service %s, grpc server may not exist, local conn is %v, please check zookeeper server %v, path: %s", serviceName, s.localConns, s.zkServers, s.zkRoot)
		}
		for _, addr := range addrs {
			cc, err := grpc.DialContext(ctx, addr.Addr, append(s.options, opts...)...)
			if err != nil {
				log.ZError(context.Background(), "dialContext failed", err, "addr", addr.Addr, "opts", append(s.options, opts...))
				return nil, errs.Wrap(err)
			}
			conns = append(conns, cc)
		}
		s.localConns[serviceName] = conns
	}
	return conns, nil
}

func (s *ZkClient) GetConn(ctx context.Context, serviceName string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	newOpts := append(s.options, grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, s.balancerName)))
	s.logger.Printf("get conn from client, serviceName: %s", serviceName)
	return grpc.DialContext(ctx, fmt.Sprintf("%s:///%s", s.scheme, serviceName), append(newOpts, opts...)...)
}

func (s *ZkClient) GetSelfConnTarget() string {
	return s.rpcRegisterAddr
}

func (s *ZkClient) CloseConn(conn *grpc.ClientConn) {
	conn.Close()
}
