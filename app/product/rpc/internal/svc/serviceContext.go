package svc

import (
	"erp/app/product/rpc/internal/config"
	"erp/app/product/rpc/internal/model"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config               config.Config
	ProductModel         model.ProductModel
	ProductBatchModel    model.ProductBatchModel
	ProductCategoryModel model.ProductCategoryModel
	BizRedis             *redis.Redis
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
		Config:               c,
		ProductModel:         model.NewProductModel(sqlx.NewMysql(c.DataSource), c.CacheRedis),
		ProductBatchModel:    model.NewProductBatchModel(sqlx.NewMysql(c.DataSource), c.CacheRedis),
		ProductCategoryModel: model.NewProductCategoryModel(sqlx.NewMysql(c.DataSource), c.CacheRedis),
		BizRedis:             rds,
	}
}
