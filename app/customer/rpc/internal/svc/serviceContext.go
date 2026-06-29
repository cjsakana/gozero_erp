package svc

import (
	"erp/app/customer/rpc/internal/config"
	"erp/app/customer/rpc/internal/model"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config                config.Config
	CustomerModel         model.CustomerModel
	CustomerCategoryModel model.CustomerCategoryModel
	SatisfactionModel     model.CustomerSatisfactionSurveyModel
	BizRedis              *redis.Redis
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
		Config:                c,
		CustomerModel:         model.NewCustomerModel(sqlx.NewMysql(c.DataSource), c.CacheRedis),
		CustomerCategoryModel: model.NewCustomerCategoryModel(sqlx.NewMysql(c.DataSource), c.CacheRedis),
		SatisfactionModel:     model.NewCustomerSatisfactionSurveyModel(sqlx.NewMysql(c.DataSource), c.CacheRedis),
		BizRedis:              rds,
	}
}
