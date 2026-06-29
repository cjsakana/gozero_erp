package svc

import (
	"erp/app/image/rpc/internal/config"
	"erp/app/image/rpc/internal/model"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config     config.Config
	ImageModel model.ImageModel
	BizRedis   *redis.Redis
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
		Config:     c,
		ImageModel: model.NewImageModel(sqlx.NewMysql(c.DataSource), c.CacheRedis),
		BizRedis:   rds,
	}
}
