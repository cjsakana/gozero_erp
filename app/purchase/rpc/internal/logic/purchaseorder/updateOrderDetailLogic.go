package purchaseorderlogic

import (
	"context"

	"erp/app/purchase/rpc/internal/svc"
	"erp/app/purchase/rpc/internal/types"
	"erp/app/purchase/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateOrderDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateOrderDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateOrderDetailLogic {
	return &UpdateOrderDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新采购订单明细（动态SQL拼接）
func (l *UpdateOrderDetailLogic) UpdateOrderDetail(in *pb.UpdateOrderDetailReq) (*pb.UpdateOrderDetailResp, error) {
	// 构建更新参数
	param := &types.UpdateOrderDetailParam{
		Id: in.Id,
	}

	// 只有非零值才设置指针
	if in.ProductId != 0 {
		param.ProductId = &in.ProductId
	}
	if in.ProductName != "" {
		param.ProductName = &in.ProductName
	}
	if in.CategoryType != 0 {
		param.CategoryType = &in.CategoryType
	}
	if in.Quantity != 0 {
		param.Quantity = &in.Quantity
	}
	if in.UnitPrice != 0 {
		param.UnitPrice = &in.UnitPrice
	}
	if in.Amount != 0 {
		param.Amount = &in.Amount
	}
	if in.ReceivedQty != 0 {
		param.ReceivedQty = &in.ReceivedQty
	}
	if in.Remark != "" {
		param.Remark = &in.Remark
	}

	// 调用model层更新方法
	err := l.svcCtx.PurchaseOrderDetailModel.UpdateOrderDetail(l.ctx, param)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateOrderDetailResp{}, nil
}
