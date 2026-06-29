package svc

import (
	"erp/app/customer/rpc/customer"
	"erp/app/finance/api/internal/config"
	"erp/app/finance/api/internal/middleware"
	"erp/app/finance/rpc/client/finance"
	"erp/app/hr/rpc/client/hr"
	"erp/app/supplier/rpc/supplier"
	"erp/common/interceptors"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config         config.Config
	AuthMiddleware rest.Middleware

	FinanceRPC  finance.FinanceZrpcClient
	HrRPC       hr.HrZrpcClient
	SupplierRPC supplier.SupplierZrpcClient
	CustomerRPC customer.CustomerZrpcClient
	BizRedis    *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 自定义拦截器
	financeRPC := zrpc.MustNewClient(c.FinanceRPC, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))
	hrRPC := zrpc.MustNewClient(c.HrRPC, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))
	supplierRPC := zrpc.MustNewClient(c.SupplierRPC, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))
	customerRPC := zrpc.MustNewClient(c.CustomerRPC, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))

	return &ServiceContext{
		Config:      c,
		BizRedis:    redis.New(c.BizRedis.Host, redis.WithPass(c.BizRedis.Pass)),
		FinanceRPC:  finance.NewFinanceZrpcClient(financeRPC),
		HrRPC:       hr.NewHrZrpcClient(hrRPC),
		SupplierRPC: supplier.NewSupplierZrpcClient(supplierRPC),
		CustomerRPC: customer.NewCustomerZrpcClient(customerRPC),
		AuthMiddleware: middleware.NewAuthMiddleware(middleware.AuthMiddlewareConfig{
			AuthConf: middleware.AuthConf{
				AccessSecret: c.Auth.AccessSecret,
				AccessExpire: c.Auth.AccessExpire,
			},
			AuthRPCConf: c.AuthRPC,
			RedisConf:   c.BizRedis,
			UserRPCConf: c.UserRPC,
		}).Handle,
	}
}
