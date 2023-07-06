package mw

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/chat/pkg/common/constant"
	"google.golang.org/grpc"
	"strconv"
)

func AddUserType() grpc.DialOption {
	return grpc.WithChainUnaryInterceptor(func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		log.ZInfo(ctx, "add user type", "method", method)
		if arr, _ := ctx.Value(constant.RpcOpUserType).([]string); len(arr) > 0 {
			userType, err := strconv.Atoi(arr[0])
			if err != nil {
				return errs.ErrInternalServer.Wrap("user type is not int")
			}
			log.ZInfo(ctx, "add user type", "method", method, "userType", userType)
			headers, _ := ctx.Value(constant.RpcCustomHeader).([]string)
			ctx = context.WithValue(ctx, constant.RpcCustomHeader, append(headers, constant.RpcOpUserType))
			ctx = context.WithValue(ctx, constant.RpcOpUserType, userType)
		}
		return invoker(ctx, method, req, reply, cc, opts...)
	})
}
