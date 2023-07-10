package mw

import (
	"context"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"google.golang.org/grpc"

	"github.com/OpenIMSDK/chat/pkg/common/constant"
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
