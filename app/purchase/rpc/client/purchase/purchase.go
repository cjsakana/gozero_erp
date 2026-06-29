package purchase

import (
	"erp/app/purchase/rpc/client/purchaseorder"
	"erp/app/purchase/rpc/client/purchasereceipt"
	"erp/app/purchase/rpc/client/purchaserequisition"
	"github.com/zeromicro/go-zero/zrpc"
)

type PurchaseZrpcClient struct {
	purchaseorder.PurchaseOrderZrpcClient
	purchasereceipt.PurchaseReceiptZrpcClient
	purchaserequisition.PurchaseRequisitionZrpcClient
}

func NewPurchaseZrpcClient(cli zrpc.Client) PurchaseZrpcClient {
	return PurchaseZrpcClient{
		PurchaseOrderZrpcClient:       purchaseorder.NewPurchaseOrderZrpcClient(cli),
		PurchaseReceiptZrpcClient:     purchasereceipt.NewPurchaseReceiptZrpcClient(cli),
		PurchaseRequisitionZrpcClient: purchaserequisition.NewPurchaseRequisitionZrpcClient(cli),
	}
}
