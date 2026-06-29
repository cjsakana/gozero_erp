package purchaseorderlogic

import (
	"context"

	"erp/app/purchase/rpc/internal/svc"
	"erp/app/purchase/rpc/internal/types"
	"erp/app/purchase/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateOrderLogic {
	return &UpdateOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新采购订单（动态SQL拼接）
func (l *UpdateOrderLogic) UpdateOrder(in *pb.UpdateOrderReq) (*pb.UpdateOrderResp, error) {
	// 构建更新参数
	param := &types.UpdateOrderParam{
		Id: in.Id,
	}

	// 只有非零值才设置指针
	if in.SupplierId != 0 {
		param.SupplierId = &in.SupplierId
	}
	if in.OrderDate != 0 {
		param.OrderDate = &in.OrderDate
	}
	if in.ExpectedDate != 0 {
		param.ExpectedDate = &in.ExpectedDate
	}
	if in.TotalAmount != 0 {
		param.TotalAmount = &in.TotalAmount
	}
	if in.Status != 0 {
		param.Status = &in.Status
	}
	if in.PurchaserId != 0 {
		param.PurchaserId = &in.PurchaserId
	}
	if in.ContractUrl != "" {
		param.ContractUrl = &in.ContractUrl
	}

	// 调用model层更新方法
	err := l.svcCtx.PurchaseOrderModel.UpdateOrder(l.ctx, param)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateOrderResp{}, nil
}
