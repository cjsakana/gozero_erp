package svc

import (
	"erp/app/auth/rpc/client/userrole"
	"erp/app/user/api/internal/config"
	"erp/app/user/api/internal/middleware"
	"erp/app/user/rpc/user"
	"erp/common/interceptors"
	"erp/common/upload"
	"time"

	"github.com/zeromicro/go-zero/core/collection"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config

	UserRPC     user.UserZrpcClient
	UserRoleRPC userrole.UserRoleZrpcClient

	LocCache       *collection.Cache
	AuthMiddleware rest.Middleware
	BizRedis       *redis.Redis

	UploadClient upload.Oss
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 自定义拦截器
	userRPC := zrpc.MustNewClient(c.UserRPC, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))
	authRPC := zrpc.MustNewClient(c.AuthRPC, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))

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
		LocCache:    locCache,
		BizRedis:    redis.New(c.BizRedis.Host, redis.WithPass(c.BizRedis.Pass)),
		UserRPC:     user.NewUserZrpcClient(userRPC),
		UserRoleRPC: userrole.NewUserRoleZrpcClient(authRPC),

		UploadClient: upload.NewR2Client(&c.R2Conf),
	}
}
