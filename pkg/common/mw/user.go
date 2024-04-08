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

package mw

import (
	"context"

	"github.com/openimsdk/chat/pkg/common/constant"
	"github.com/openimsdk/tools/log"
	"google.golang.org/grpc"
)

func AddUserType() grpc.DialOption {
	return grpc.WithChainUnaryInterceptor(func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		log.ZInfo(ctx, "add user type", "method", method)
		if arr, _ := ctx.Value(constant.RpcOpUserType).([]string); len(arr) > 0 {
			log.ZInfo(ctx, "add user type", "method", method, "userType", arr)
			headers, _ := ctx.Value(constant.RpcCustomHeader).([]string)
			ctx = context.WithValue(ctx, constant.RpcCustomHeader, append(headers, constant.RpcOpUserType))
			ctx = context.WithValue(ctx, constant.RpcOpUserType, arr)
		}
		return invoker(ctx, method, req, reply, cc, opts...)
	})
}
