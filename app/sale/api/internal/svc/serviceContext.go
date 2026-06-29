package svc

import (
	"erp/app/customer/rpc/customer"
	"erp/app/hr/rpc/client/hr"
	"erp/app/inventory/rpc/client/inventory"
	"erp/app/product/rpc/client/prod"
	"erp/app/sale/api/internal/config"
	"erp/app/sale/api/internal/middleware"
	"erp/app/sale/rpc/client/sale"
	"erp/common/interceptors"
	"erp/common/upload"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config         config.Config
	AuthMiddleware rest.Middleware

	SaleRPC     sale.SaleZrpcClient
	HrRPC       hr.HrZrpcClient
	CustomerRPC customer.CustomerZrpcClient

	InventoryRPC inventory.InventoryZrpcClient
	ProductRPC   prod.ProdZrpcClient
	BizRedis     *redis.Redis

	UploadClient upload.Oss
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 自定义拦截器
	saleRPC := zrpc.MustNewClient(c.SaleRPC, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))
	hrRPC := zrpc.MustNewClient(c.HrRPC, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))
	iRPC := zrpc.MustNewClient(c.InventoryRPC, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))
	prodRPC := zrpc.MustNewClient(c.ProductRPC, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))
	customerRPC := zrpc.MustNewClient(c.CustomerRPC, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))

	return &ServiceContext{
		Config:       c,
		BizRedis:     redis.New(c.BizRedis.Host, redis.WithPass(c.BizRedis.Pass)),
		SaleRPC:      sale.NewSaleZrpcClient(saleRPC),
		HrRPC:        hr.NewHrZrpcClient(hrRPC),
		InventoryRPC: inventory.NewInventoryZrpcClient(iRPC),
		ProductRPC:   prod.NewProdZrpcClient(prodRPC),
		CustomerRPC:  customer.NewCustomerZrpcClient(customerRPC),

		AuthMiddleware: middleware.NewAuthMiddleware(middleware.AuthMiddlewareConfig{
			AuthConf: middleware.AuthConf{
				AccessSecret: c.Auth.AccessSecret,
				AccessExpire: c.Auth.AccessExpire,
			},
			AuthRPCConf: c.AuthRPC,
			RedisConf:   c.BizRedis,
		}).Handle,
		UploadClient: upload.NewR2Client(&c.R2Conf),
	}
}
