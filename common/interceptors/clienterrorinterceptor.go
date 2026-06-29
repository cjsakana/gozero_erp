package interceptors

import (
	"context"
	"erp/common/xcode"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func ClientErrorInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption) error {
		// 调用原始的 gRPC 方法，返回一个错误
		// 注意！！！这个错误可能是一个 gRPC 状态错误（Status），也可能是封装了gRPC status其他类型的错误
		err := invoker(ctx, method, req, reply, cc, opts...)
		if err != nil {
			// 需要从疑似gRPC status 中提取真正的 gRPC status
			grpcStatus, _ := status.FromError(err)
			// 再将 gRPC 状态转换为自定义的 xcode，保证万无一失
			xc := xcode.GrpcStatusToXCode(grpcStatus)
			// 将自定义的 xcode 作为新错误返回
			err = errors.WithMessage(xc, grpcStatus.Message())
		}

		return err
	}
}
