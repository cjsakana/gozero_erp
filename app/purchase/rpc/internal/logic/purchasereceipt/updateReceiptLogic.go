package purchasereceiptlogic

import (
	"context"

	"erp/app/purchase/rpc/internal/svc"
	"erp/app/purchase/rpc/internal/types"
	"erp/app/purchase/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateReceiptLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateReceiptLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateReceiptLogic {
	return &UpdateReceiptLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新采购入库单（动态SQL拼接）
func (l *UpdateReceiptLogic) UpdateReceipt(in *pb.UpdateReceiptReq) (*pb.UpdateReceiptResp, error) {
	// 构建更新参数
	param := &types.UpdateReceiptParam{
		Id: in.Id,
	}

	// 只有非零值才设置指针
	if in.OrderId != 0 {
		param.OrderId = &in.OrderId
	}
	if in.WarehouseId != 0 {
		param.WarehouseId = &in.WarehouseId
	}
	if in.ReceiptDate != 0 {
		param.ReceiptDate = &in.ReceiptDate
	}
	if in.TotalQuantity != 0 {
		param.TotalQuantity = &in.TotalQuantity
	}
	if in.TotalAmount != 0 {
		param.TotalAmount = &in.TotalAmount
	}
	if in.Status != 0 {
		param.Status = &in.Status
	}
	if in.CreatedBy != 0 {
		param.CreatedBy = &in.CreatedBy
	}

	// 调用model层更新方法
	err := l.svcCtx.PurchaseReceiptModel.UpdateReceipt(l.ctx, param)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateReceiptResp{}, nil
}
