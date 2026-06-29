package svc

import (
	"erp/app/auth/rpc/internal/config"
	"erp/app/auth/rpc/internal/model"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config              config.Config
	PermissionModel     model.PermissionModel
	RolePermissionModel model.RolePermissionModel
	RoleModel           model.RoleModel
	UserRoleModel       model.UserRoleModel
	BizRedis            *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	rds, err := redis.NewRedis(redis.RedisConf{
		Host: c.BizRedis.Host,
		Pass: c.BizRedis.Pass,
		Type: c.BizRedis.Type,
	})
	if err != nil {
		panic(err)
	}

	return &ServiceContext{
		Config:              c,
		PermissionModel:     model.NewPermissionModel(sqlx.NewMysql(c.DataSource), c.CacheRedis),
		RolePermissionModel: model.NewRolePermissionModel(sqlx.NewMysql(c.DataSource), c.CacheRedis),
		RoleModel:           model.NewRoleModel(sqlx.NewMysql(c.DataSource), c.CacheRedis),
		UserRoleModel:       model.NewUserRoleModel(sqlx.NewMysql(c.DataSource), c.CacheRedis),
		BizRedis:            rds,
	}
}
