package svc

import (
	"erp/app/production/rpc/internal/config"
	"erp/app/production/rpc/internal/model"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config                       config.Config
	BomModel                     model.BomModel
	BomItemModel                 model.BomItemModel
	WorkOrderModel               model.WorkOrderModel
	MaterialRequisitionModel     model.MaterialRequisitionModel
	MaterialRequisitionItemModel model.MaterialRequisitionItemModel
	BizRedis                     *redis.Redis
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
		Config:                       c,
		BizRedis:                     rds,
		BomModel:                     model.NewBomModel(sqlx.NewMysql(c.DataSource), c.CacheRedis),
		BomItemModel:                 model.NewBomItemModel(sqlx.NewMysql(c.DataSource), c.CacheRedis),
		WorkOrderModel:               model.NewWorkOrderModel(sqlx.NewMysql(c.DataSource), c.CacheRedis),
		MaterialRequisitionModel:     model.NewMaterialRequisitionModel(sqlx.NewMysql(c.DataSource), c.CacheRedis),
		MaterialRequisitionItemModel: model.NewMaterialRequisitionItemModel(sqlx.NewMysql(c.DataSource), c.CacheRedis),
	}
}
