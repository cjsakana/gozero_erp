package svc

import (
	"erp/app/hr/rpc/client/hr"
	"erp/app/inventory/rpc/client/inventory"
	"erp/app/product/api/internal/config"
	"erp/app/product/api/internal/middleware"
	"erp/app/product/rpc/client/prod"
	"erp/common/interceptors"
	"time"

	"github.com/zeromicro/go-zero/core/collection"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config         config.Config
	AuthMiddleware rest.Middleware

	ProductRPC   prod.ProdZrpcClient
	BizRedis     *redis.Redis
	LocCache     *collection.Cache
	HrRPC        hr.HrZrpcClient
	InventoryRPC inventory.InventoryZrpcClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 自定义拦截器
	prodRPC := zrpc.MustNewClient(c.ProductRPC, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))
	hrRPC := zrpc.MustNewClient(c.HrRPC, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))
	iRPC := zrpc.MustNewClient(c.InventoryRPC, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))

	// 本地缓存
	locCache, _ := collection.NewCache(7 * time.Duration(c.Auth.AccessExpire) * time.Second)

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
		LocCache:     locCache,
		BizRedis:     redis.New(c.BizRedis.Host, redis.WithPass(c.BizRedis.Pass)),
		ProductRPC:   prod.NewProdZrpcClient(prodRPC),
		HrRPC:        hr.NewHrZrpcClient(hrRPC),
		InventoryRPC: inventory.NewInventoryZrpcClient(iRPC),
	}
}
