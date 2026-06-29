package svc

import (
	"erp/app/hr/rpc/client/hr"
	"erp/app/inventory/rpc/client/inventory"
	"erp/app/product/rpc/client/prod"
	"erp/app/purchase/api/internal/config"
	"erp/app/purchase/api/internal/middleware"
	"erp/app/purchase/rpc/client/purchase"
	"erp/app/supplier/rpc/supplier"
	"erp/common/interceptors"
	"erp/common/upload"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config         config.Config
	AuthMiddleware rest.Middleware

	PurchaseRPC  purchase.PurchaseZrpcClient
	HrRPC        hr.HrZrpcClient
	InventoryRPC inventory.InventoryZrpcClient
	SupplierRPC  supplier.SupplierZrpcClient
	ProductRPC   prod.ProdZrpcClient
	BizRedis     *redis.Redis

	UploadClient upload.Oss
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 自定义拦截器
	purchaseRPC := zrpc.MustNewClient(c.PurchaseRPC, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))
	hrRPC := zrpc.MustNewClient(c.HrRPC, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))
	inventoryRPC := zrpc.MustNewClient(c.InventoryRPC, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))
	supplierRPC := zrpc.MustNewClient(c.SupplierRPC, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))
	prodRPC := zrpc.MustNewClient(c.ProductRPC, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))

	return &ServiceContext{
		Config:       c,
		BizRedis:     redis.New(c.BizRedis.Host, redis.WithPass(c.BizRedis.Pass)),
		PurchaseRPC:  purchase.NewPurchaseZrpcClient(purchaseRPC),
		HrRPC:        hr.NewHrZrpcClient(hrRPC),
		InventoryRPC: inventory.NewInventoryZrpcClient(inventoryRPC),
		SupplierRPC:  supplier.NewSupplierZrpcClient(supplierRPC),
		ProductRPC:   prod.NewProdZrpcClient(prodRPC),

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
