package svc

import (
	"erp/app/sale/rpc/internal/config"
	"erp/app/sale/rpc/internal/model"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config                   config.Config
	SalesDeliveryModel       model.SalesDeliveryModel
	SalesDeliveryDetailModel model.SalesDeliveryDetailModel
	SalesOrderModel          model.SalesOrderModel
	SalesOrderDetailModel    model.SalesOrderDetailModel
	BizRedis                 *redis.Redis
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
		Config:                   c,
		BizRedis:                 rds,
		SalesDeliveryModel:       model.NewSalesDeliveryModel(sqlx.NewMysql(c.DataSource), c.CacheRedis),
		SalesDeliveryDetailModel: model.NewSalesDeliveryDetailModel(sqlx.NewMysql(c.DataSource), c.CacheRedis),
		SalesOrderModel:          model.NewSalesOrderModel(sqlx.NewMysql(c.DataSource), c.CacheRedis),
		SalesOrderDetailModel:    model.NewSalesOrderDetailModel(sqlx.NewMysql(c.DataSource), c.CacheRedis),
	}
}
