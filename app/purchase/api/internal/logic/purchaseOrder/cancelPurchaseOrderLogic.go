package purchaseOrder

import (
	"context"
	"erp/app/purchase/api/internal/svc"
	"erp/app/purchase/api/internal/types"
	"erp/app/purchase/rpc/pb"
	"erp/common/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type CancelPurchaseOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCancelPurchaseOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelPurchaseOrderLogic {
	return &CancelPurchaseOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CancelPurchaseOrderLogic) CancelPurchaseOrder(req *types.CancelPurchaseOrderReq) (resp *types.CancelPurchaseOrderResp, err error) {
	orderId, err := util.StringToInt64(req.OrderId)
	if err != nil {
		return nil, err
	}
	_, err = l.svcCtx.PurchaseRPC.CancelPurchaseOrder(l.ctx, &pb.CancelPurchaseOrderReq{
		OrderId:      orderId,
		CancelReason: req.CancelReason,
	})
	if err != nil {
		return nil, err
	}

	resp = &types.CancelPurchaseOrderResp{}
	return
}
