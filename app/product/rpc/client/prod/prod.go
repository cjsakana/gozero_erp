package prod

import (
	"erp/app/product/rpc/client/product"
	"erp/app/product/rpc/client/productbatch"
	"erp/app/product/rpc/client/productcategory"
	"github.com/zeromicro/go-zero/zrpc"
)

type ProdZrpcClient struct {
	product.ProductZrpcClient
	productcategory.ProductCategoryZrpcClient
	productbatch.ProductBatchZrpcClient
}

func NewProdZrpcClient(cli zrpc.Client) ProdZrpcClient {
	return ProdZrpcClient{
		product.NewProductZrpcClient(cli),
		productcategory.NewProductCategoryZrpcClient(cli),
		productbatch.NewProductBatchZrpcClient(cli),
	}
}
