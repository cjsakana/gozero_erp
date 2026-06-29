package svc

import (
	"erp/app/purchase/rpc/internal/config"
	"erp/app/purchase/rpc/internal/model"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config                         config.Config
	PurchaseOrderModel             model.PurchaseOrderModel
	PurchaseOrderDetailModel       model.PurchaseOrderDetailModel
	PurchaseReceiptModel           model.PurchaseReceiptModel
	PurchaseReceiptDetailModel     model.PurchaseReceiptDetailModel
	PurchaseRequisitionModel       model.PurchaseRequisitionModel
	PurchaseRequisitionDetailModel model.PurchaseRequisitionDetailModel
	BizRedis                       *redis.Redis
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
	sqlConn := sqlx.NewMysql(c.DataSource)
	return &ServiceContext{
		Config:                         c,
		PurchaseOrderModel:             model.NewPurchaseOrderModel(sqlConn, c.CacheRedis),
		PurchaseOrderDetailModel:       model.NewPurchaseOrderDetailModel(sqlConn, c.CacheRedis),
		PurchaseReceiptModel:           model.NewPurchaseReceiptModel(sqlConn, c.CacheRedis),
		PurchaseReceiptDetailModel:     model.NewPurchaseReceiptDetailModel(sqlConn, c.CacheRedis),
		PurchaseRequisitionModel:       model.NewPurchaseRequisitionModel(sqlConn, c.CacheRedis),
		PurchaseRequisitionDetailModel: model.NewPurchaseRequisitionDetailModel(sqlConn, c.CacheRedis),
		BizRedis:                       rds,
	}
}
