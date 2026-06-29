package svc

import (
	"erp/app/hr/rpc/client/hr"
	"erp/app/supplier/api/internal/config"
	"erp/app/supplier/api/internal/middleware"
	"erp/app/supplier/rpc/supplier"
	"erp/common/interceptors"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config         config.Config
	AuthMiddleware rest.Middleware

	SupplierRPC supplier.SupplierZrpcClient
	HrRPC       hr.HrZrpcClient
	BizRedis    *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 自定义拦截器
	supplierRPC := zrpc.MustNewClient(c.SupplierRPC, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))
	hrRPC := zrpc.MustNewClient(c.HrRPC, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))

	return &ServiceContext{
		Config: c,
		AuthMiddleware: middleware.NewAuthMiddleware(middleware.AuthMiddlewareConfig{
			AuthConf: middleware.AuthConf{
				AccessSecret: c.Auth.AccessSecret,
				AccessExpire: c.Auth.AccessExpire,
			},
			AuthRPCConf: c.AuthRPC,
			RedisConf:   c.BizRedis,
		}).Handle,
		BizRedis:    redis.New(c.BizRedis.Host, redis.WithPass(c.BizRedis.Pass)),
		SupplierRPC: supplier.NewSupplierZrpcClient(supplierRPC),
		HrRPC:       hr.NewHrZrpcClient(hrRPC),
	}
}
