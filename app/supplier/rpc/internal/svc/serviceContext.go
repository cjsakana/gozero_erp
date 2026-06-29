package svc

import (
	"erp/app/supplier/rpc/internal/config"
	"erp/app/supplier/rpc/internal/model"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config          config.Config
	SupplierModel   model.SupplierModel
	EvaluationModel model.SupplierEvaluationModel
	BizRedis        *redis.Redis
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
		Config:          c,
		SupplierModel:   model.NewSupplierModel(sqlx.NewMysql(c.DataSource), c.CacheRedis),
		EvaluationModel: model.NewSupplierEvaluationModel(sqlx.NewMysql(c.DataSource), c.CacheRedis),
		BizRedis:        rds,
	}
}
