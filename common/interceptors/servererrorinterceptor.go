package interceptors

import (
	"context"
	"erp/common/xcode"

	"google.golang.org/grpc"
)

func ServerErrorInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{},
		info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		// handler就是正常业务处理对应的rpc方法
		resp, err = handler(ctx, req)
		// 将rpc方法【如Register】返回的error，转成对应gRPC能识别的error
		return resp, xcode.FromError(err).Err()
	}
}
