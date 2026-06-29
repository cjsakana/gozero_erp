package svc

import (
	"erp/app/finance/rpc/internal/config"
	"erp/app/finance/rpc/internal/model"
	"erp/app/hr/rpc/client/hr"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config             config.Config
	FixedAssetModel    model.FixedAssetModel
	PaymentRecordModel model.PaymentRecordModel
	ReceiptRecordModel model.ReceiptRecordModel
	SalaryPaymentModel model.SalaryPaymentModel
	BizRedis           *redis.Redis
	HrRPC              hr.HrZrpcClient // HR RPC 客户端
	DtmServer          string          // DTM 服务器地址
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
	
	// 初始化 HR RPC 客户端
	hrRPC := zrpc.MustNewClient(c.HrRPC)
	
	return &ServiceContext{
		Config:             c,
		FixedAssetModel:    model.NewFixedAssetModel(sqlConn, c.CacheRedis),
		PaymentRecordModel: model.NewPaymentRecordModel(sqlConn, c.CacheRedis),
		ReceiptRecordModel: model.NewReceiptRecordModel(sqlConn, c.CacheRedis),
		SalaryPaymentModel: model.NewSalaryPaymentModel(sqlConn, c.CacheRedis),
		BizRedis:           rds,
		HrRPC:              hr.NewHrZrpcClient(hrRPC),
		DtmServer:          c.DtmServer,
	}
}
