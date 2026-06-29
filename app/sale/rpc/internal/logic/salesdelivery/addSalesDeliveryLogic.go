package salesdeliverylogic

import (
	"context"
	"erp/app/sale/rpc/internal/code"
	"erp/app/sale/rpc/internal/svc"
	"erp/app/sale/rpc/internal/types"
	"erp/app/sale/rpc/pb"
	"erp/common/util"
	"github.com/zeromicro/go-zero/core/logx"
)

type AddSalesDeliveryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddSalesDeliveryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddSalesDeliveryLogic {
	return &AddSalesDeliveryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// -----------------------salesDelivery-----------------------
func (l *AddSalesDeliveryLogic) AddSalesDelivery(in *pb.AddSalesDeliveryReq) (*pb.AddSalesDeliveryResp, error) {
	// 转换 proto → model param
	deliveryId := util.GenerateSnowflake()
	param := &types.AddSalesDeliveryParam{
		Id:            deliveryId,
		DeliveryNo:    in.DeliveryNo,
		OrderId:       in.OrderId,
		WarehouseId:   in.WarehouseId,
		DeliveryDate:  in.DeliveryDate,
		TotalQuantity: in.TotalQuantity,
		TotalAmount:   in.TotalAmount,
		CreatedBy:     in.CreatedBy,
	}

	for _, d := range in.Details {
		id := util.GenerateSnowflake()
		param.Details = append(param.Details, types.DeliveryDetailParam{
			Id:          id,
			DeliveryId:  deliveryId,
			ProductId:   d.ProductId,
			ProductName: d.ProductName,
			Unit:        d.Unit,
			Quantity:    d.Quantity,
			UnitPrice:   d.UnitPrice,
			Amount:      d.Amount,
			BatchId:     d.BatchId,
		})
	}

	err := l.svcCtx.SalesDeliveryModel.AddWithDetails(l.ctx, param)
	if err != nil {

		return nil, code.CreateDeliveryFail

	}

	return &pb.AddSalesDeliveryResp{
		Id: deliveryId,
	}, nil
}
