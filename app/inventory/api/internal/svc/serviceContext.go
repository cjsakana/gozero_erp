package svc

import (
	"erp/app/hr/rpc/client/hr"
	"erp/app/inventory/api/internal/config"
	"erp/app/inventory/api/internal/middleware"
	"erp/app/inventory/rpc/client/inventory"
	"erp/app/product/rpc/client/prod"
	"erp/app/user/rpc/user"
	"erp/common/interceptors"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config         config.Config
	AuthMiddleware rest.Middleware

	InventoryRPC inventory.InventoryZrpcClient
	ProductRPC   prod.ProdZrpcClient
	HrRPC        hr.HrZrpcClient
	UserRPC      user.UserZrpcClient
	BizRedis     *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 自定义拦截器
	iRPC := zrpc.MustNewClient(c.InventoryRPC, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))
	productRPC := zrpc.MustNewClient(c.ProductRPC, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))
	hrRPC := zrpc.MustNewClient(c.HrRPC, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))
	userRPC := zrpc.MustNewClient(c.UserRPC, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))

	return &ServiceContext{
		Config:   c,
		BizRedis: redis.New(c.BizRedis.Host, redis.WithPass(c.BizRedis.Pass)),
		AuthMiddleware: middleware.NewAuthMiddleware(middleware.AuthMiddlewareConfig{
			AuthConf: middleware.AuthConf{
				AccessSecret: c.Auth.AccessSecret,
				AccessExpire: c.Auth.AccessExpire,
			},
			AuthRPCConf: c.AuthRPC,
			RedisConf:   c.BizRedis,
		}).Handle,
		InventoryRPC: inventory.NewInventoryZrpcClient(iRPC),
		ProductRPC:   prod.NewProdZrpcClient(productRPC),
		HrRPC:        hr.NewHrZrpcClient(hrRPC),
		UserRPC:      user.NewUserZrpcClient(userRPC),
	}
}
