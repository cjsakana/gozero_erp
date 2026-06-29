package salesDelivery

import (
	"context"
	"erp/app/sale/api/internal/svc"
	"erp/app/sale/api/internal/types"
	"erp/app/sale/rpc/pb"
	"erp/common/util"
	"erp/common/xtypes"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddSalesDeliveryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddSalesDeliveryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddSalesDeliveryLogic {
	return &AddSalesDeliveryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddSalesDeliveryLogic) AddSalesDelivery(req *types.AddSalesDeliveryReq) (resp *types.AddSalesDeliveryResp, err error) {
	createdBy, err := util.GetInt64FromCtx(l.ctx, xtypes.EmployeeIdKey)
	if err != nil {
		return nil, err
	}

	var pbDetails []*pb.DeliveryDetailItem
	for _, d := range req.Details {

		productId, err := util.StringToInt64(d.ProductId)
		if err != nil {
			return nil, err
		}
		batchId, err := util.StringToInt64(d.BatchId)
		if err != nil {
			return nil, err
		}

		pbDetails = append(pbDetails, &pb.DeliveryDetailItem{
			ProductId:   productId,
			ProductName: d.ProductName,
			Unit:        d.Unit,
			Quantity:    d.Quantity,
			UnitPrice:   d.UnitPrice,
			Amount:      d.Amount,
			BatchId:     batchId,
		})
	}

	orderId, err := util.StringToInt64(req.OrderId)
	if err != nil {
		return nil, err
	}
	warehouseId, err := util.StringToInt64(req.WarehouseId)
	if err != nil {
		return nil, err
	}
	deliveryNo := util.GenerateNo("DEL")
	ret, err := l.svcCtx.SaleRPC.AddSalesDelivery(l.ctx, &pb.AddSalesDeliveryReq{
		DeliveryNo:    deliveryNo,
		OrderId:       orderId,
		WarehouseId:   warehouseId,
		DeliveryDate:  req.DeliveryDate,
		TotalQuantity: req.TotalQuantity,
		TotalAmount:   req.TotalAmount,
		CreatedBy:     createdBy,
		Details:       pbDetails,
	})
	if err != nil {
		return nil, err
	}

	resp = &types.AddSalesDeliveryResp{
		Id: util.Int64ToString(ret.Id),
	}

	return
}
