package finance

import (
	"erp/app/finance/rpc/client/fixedasset"
	"erp/app/finance/rpc/client/paymentrecord"
	"erp/app/finance/rpc/client/receiptrecord"
	"erp/app/finance/rpc/client/salarypayment"
	"github.com/zeromicro/go-zero/zrpc"
)

type FinanceZrpcClient struct {
	fixedasset.FixedAssetZrpcClient
	paymentrecord.PaymentRecordZrpcClient
	receiptrecord.ReceiptRecordZrpcClient
	salarypayment.SalaryPaymentZrpcClient
}

func NewFinanceZrpcClient(cli zrpc.Client) FinanceZrpcClient {
	return FinanceZrpcClient{
		FixedAssetZrpcClient:    fixedasset.NewFixedAssetZrpcClient(cli),
		PaymentRecordZrpcClient: paymentrecord.NewPaymentRecordZrpcClient(cli),
		ReceiptRecordZrpcClient: receiptrecord.NewReceiptRecordZrpcClient(cli),
		SalaryPaymentZrpcClient: salarypayment.NewSalaryPaymentZrpcClient(cli),
	}
}
