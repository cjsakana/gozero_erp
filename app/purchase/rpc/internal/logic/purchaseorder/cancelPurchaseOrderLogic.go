package purchaseorderlogic

import (
	"context"

	"erp/app/purchase/rpc/internal/code"
	"erp/app/purchase/rpc/internal/svc"
	"erp/app/purchase/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CancelPurchaseOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCancelPurchaseOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelPurchaseOrderLogic {
	return &CancelPurchaseOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 取消采购订单
func (l *CancelPurchaseOrderLogic) CancelPurchaseOrder(in *pb.CancelPurchaseOrderReq) (*pb.CancelPurchaseOrderResp, error) {
	err := l.svcCtx.PurchaseOrderModel.CancelOrder(l.ctx, in.OrderId)
	if err != nil {
		return nil, code.CancelOrderFail
	}
	return &pb.CancelPurchaseOrderResp{}, nil
}
