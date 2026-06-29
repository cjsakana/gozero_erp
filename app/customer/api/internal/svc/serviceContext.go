package svc

import (
	"erp/app/customer/api/internal/config"
	"erp/app/customer/api/internal/middleware"
	"erp/app/customer/rpc/customer"
	"erp/app/hr/rpc/client/hr"
	"erp/common/interceptors"

	"github.com/zeromicro/go-zero/core/collection"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config         config.Config
	AuthMiddleware rest.Middleware

	CustomerRPC customer.CustomerZrpcClient
	HrRPC       hr.HrZrpcClient
	BizRedis    *redis.Redis
	LocCache    *collection.Cache
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 自定义拦截器
	customerRPC := zrpc.MustNewClient(c.CustomerRPC, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))
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
		CustomerRPC: customer.NewCustomerZrpcClient(customerRPC),
		HrRPC:       hr.NewHrZrpcClient(hrRPC),
	}
}
