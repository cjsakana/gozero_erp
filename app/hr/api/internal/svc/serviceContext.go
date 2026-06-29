package svc

import (
	"erp/app/auth/rpc/client/auth"
	"erp/app/hr/api/internal/config"
	"erp/app/hr/api/internal/middleware"
	"erp/app/hr/rpc/client/hr"
	"erp/app/user/rpc/user"
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

	AuthRPC  auth.AuthZrpcClient
	HrRPC    hr.HrZrpcClient
	UserRPC  user.UserZrpcClient
	BizRedis *redis.Redis
	LocCache *collection.Cache
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 自定义拦截器
	authRPC := zrpc.MustNewClient(c.AuthRPC, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))
	hrRPC := zrpc.MustNewClient(c.HrRPC, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))
	userRPC := zrpc.MustNewClient(c.UserRPC, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))

	// 本地缓存
	locCache, _ := collection.NewCache(7 * time.Duration(c.Auth.AccessExpire) * time.Second)

	return &ServiceContext{
		Config:   c,
		LocCache: locCache,
		BizRedis: redis.New(c.BizRedis.Host, redis.WithPass(c.BizRedis.Pass)),
		HrRPC:    hr.NewHrZrpcClient(hrRPC),
		UserRPC:  user.NewUserZrpcClient(userRPC),
		AuthRPC:  auth.NewAuthZrpcClient(authRPC),

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
