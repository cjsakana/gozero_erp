package svc

import (
	"erp/app/auth/api/internal/config"
	"erp/app/auth/api/internal/middleware"
	"erp/app/auth/rpc/client/permission"
	"erp/app/auth/rpc/client/role"
	"erp/app/auth/rpc/client/rolepermission"
	"erp/app/auth/rpc/client/userrole"
	"erp/common/interceptors"
	"time"

	"github.com/zeromicro/go-zero/core/collection"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config            config.Config
	PermissionRPC     permission.PermissionZrpcClient
	RoleRPC           role.RoleZrpcClient
	RolePermissionRPC rolepermission.RolePermissionZrpcClient
	UserRoleRPC       userrole.UserRoleZrpcClient
	LocCache          *collection.Cache
	AuthMiddleware    rest.Middleware
	BizRedis          *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 自定义拦截器
	authRPC := zrpc.MustNewClient(c.AuthRPC, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))

	// 本地缓存
	locCache, _ := collection.NewCache(7 * time.Duration(c.Auth.AccessExpire) * time.Second)

	return &ServiceContext{
		Config: c,

		PermissionRPC:     permission.NewPermissionZrpcClient(authRPC),
		RoleRPC:           role.NewRoleZrpcClient(authRPC),
		RolePermissionRPC: rolepermission.NewRolePermissionZrpcClient(authRPC),
		UserRoleRPC:       userrole.NewUserRoleZrpcClient(authRPC),

		LocCache: locCache,
		BizRedis: redis.New(c.BizRedis.Host, redis.WithPass(c.BizRedis.Pass)),

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
