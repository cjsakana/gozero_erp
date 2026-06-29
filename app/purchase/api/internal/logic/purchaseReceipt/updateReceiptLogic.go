package purchaseReceipt

import (
	"context"
	"erp/common/util"
	"erp/common/xtypes"

	"erp/app/purchase/api/internal/svc"
	"erp/app/purchase/api/internal/types"
	"erp/app/purchase/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateReceiptLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateReceiptLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateReceiptLogic {
	return &UpdateReceiptLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateReceiptLogic) UpdateReceipt(req *types.UpdateReceiptReq) (resp *types.UpdateReceiptResp, err error) {
	id, err := util.StringToInt64(req.Id)
	if err != nil {
		return nil, err
	}
	orderId, err := util.StringToInt64(req.OrderId)
	if err != nil {
		return nil, err
	}
	warehouseId, err := util.StringToInt64(req.WarehouseId)
	if err != nil {
		return nil, err
	}
	employeeId := int64(l.ctx.Value(xtypes.EmployeeIdKey).(float64))

	// 调用RPC服务
	_, err = l.svcCtx.PurchaseRPC.UpdateReceipt(l.ctx, &pb.UpdateReceiptReq{
		Id:            id,
		OrderId:       orderId,
		WarehouseId:   warehouseId,
		ReceiptDate:   req.ReceiptDate,
		TotalQuantity: req.TotalQuantity,
		TotalAmount:   req.TotalAmount,
		Status:        req.Status,
		CreatedBy:     employeeId,
	})
	if err != nil {
		return nil, err
	}

	return &types.UpdateReceiptResp{}, nil
}
