package sale

import (
	"erp/app/sale/rpc/client/salesdelivery"
	"erp/app/sale/rpc/client/salesorder"
	"github.com/zeromicro/go-zero/zrpc"
)

type SaleZrpcClient struct {
	salesdelivery.SalesDeliveryZrpcClient
	salesorder.SalesOrderZrpcClient
}

func NewSaleZrpcClient(cli zrpc.Client) SaleZrpcClient {
	return SaleZrpcClient{
		SalesDeliveryZrpcClient: salesdelivery.NewSalesDeliveryZrpcClient(cli),
		SalesOrderZrpcClient:    salesorder.NewSalesOrderZrpcClient(cli),
	}
}
