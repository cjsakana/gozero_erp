package svc

import (
	"erp/app/hr/rpc/client/hr"
	"erp/app/inventory/rpc/client/inventory"
	"erp/app/product/rpc/client/prod"
	"erp/app/production/api/internal/config"
	"erp/app/production/api/internal/middleware"
	"erp/app/production/rpc/client/production"
	"erp/common/interceptors"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config         config.Config
	AuthMiddleware rest.Middleware

	ProductionRPC production.Production
	HrRPC         hr.HrZrpcClient
	InventoryRPC  inventory.InventoryZrpcClient
	ProductRPC    prod.ProdZrpcClient
	BizRedis      *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 自定义拦截器
	productionRPC := zrpc.MustNewClient(c.ProductionRPC, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))
	hrRPC := zrpc.MustNewClient(c.HrRPC, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))
	iRPC := zrpc.MustNewClient(c.InventoryRPC, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))
	prodRPC := zrpc.MustNewClient(c.ProductRPC, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))

	return &ServiceContext{
		Config:        c,
		BizRedis:      redis.New(c.BizRedis.Host, redis.WithPass(c.BizRedis.Pass)),
		ProductionRPC: production.NewProduction(productionRPC),
		HrRPC:         hr.NewHrZrpcClient(hrRPC),
		InventoryRPC:  inventory.NewInventoryZrpcClient(iRPC),
		ProductRPC:    prod.NewProdZrpcClient(prodRPC),

		AuthMiddleware: middleware.NewAuthMiddleware(middleware.AuthMiddlewareConfig{
			AuthConf: middleware.AuthConf{
				AccessSecret: c.Auth.AccessSecret,
				AccessExpire: c.Auth.AccessExpire,
			},
			AuthRPCConf: c.AuthRPC,
			RedisConf:   c.BizRedis,
		}).Handle,
	}
}
