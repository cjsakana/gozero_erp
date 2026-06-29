package svc

import (
	"erp/app/inventory/rpc/internal/config"
	"erp/app/inventory/rpc/internal/model"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config                    config.Config
	InventoryModel            model.InventoryModel
	InventoryTransactionModel model.InventoryTransactionModel
	WarehouseModel            model.WarehouseModel
	BizRedis                  *redis.Redis
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
		Config:                    c,
		InventoryModel:            model.NewInventoryModel(sqlx.NewMysql(c.DataSource), c.CacheRedis),
		InventoryTransactionModel: model.NewInventoryTransactionModel(sqlx.NewMysql(c.DataSource), c.CacheRedis),
		WarehouseModel:            model.NewWarehouseModel(sqlx.NewMysql(c.DataSource), c.CacheRedis),
		BizRedis:                  rds,
	}
}
