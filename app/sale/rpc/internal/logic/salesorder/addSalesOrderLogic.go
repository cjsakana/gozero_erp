package salesorderlogic

import (
	"context"
	"erp/common/util"

	"erp/app/sale/rpc/internal/code"
	"erp/app/sale/rpc/internal/svc"
	"erp/app/sale/rpc/internal/types"
	"erp/app/sale/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddSalesOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddSalesOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddSalesOrderLogic {
	return &AddSalesOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// -----------------------salesOrder-----------------------
func (l *AddSalesOrderLogic) AddSalesOrder(in *pb.AddSalesOrderReq) (*pb.AddSalesOrderResp, error) {
	orderId := util.GenerateSnowflake()
	param := &types.AddSalesOrderParam{
		Id:           orderId,
		OrderNo:      in.OrderNo,
		CustomerId:   in.CustomerId,
		OrderDate:    in.OrderDate,
		PromisedDate: in.PromisedDate,
		TotalAmount:  in.TotalAmount,
		Status:       in.Status,
		SalesmanId:   in.SalesmanId,
		ContractUrl:  in.ContractUrl,
	}
	for _, d := range in.Details {
		id := util.GenerateSnowflake()
		param.Details = append(param.Details, types.SalesOrderDetailParam{
			Id:           id,
			OrderId:      orderId,
			ProductId:    d.ProductId,
			ProductName:  d.ProductName,
			Unit:         d.Unit,
			Quantity:     d.Quantity,
			UnitPrice:    d.UnitPrice,
			Amount:       d.Amount,
			DeliveredQty: d.DeliveredQty,
			Remark:       d.Remark,
		})
	}

	err := l.svcCtx.SalesOrderModel.AddWithDetails(l.ctx, param)
	if err != nil {
		return nil, code.CreateSalesOrderFail
	}
	return &pb.AddSalesOrderResp{Id: orderId}, nil
}
