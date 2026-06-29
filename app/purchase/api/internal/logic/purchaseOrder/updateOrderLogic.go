package purchaseOrder

import (
	"context"
	"erp/common/util"

	"erp/app/purchase/api/internal/svc"
	"erp/app/purchase/api/internal/types"
	"erp/app/purchase/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateOrderLogic {
	return &UpdateOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateOrderLogic) UpdateOrder(req *types.UpdateOrderReq) (resp *types.UpdateOrderResp, err error) {
	id, err := util.StringToInt64(req.Id)
	if err != nil {
		return nil, err
	}
	supplierId, err := util.StringToInt64(req.SupplierId)
	if err != nil {
		return nil, err
	}
	purchaserId, err := util.StringToInt64(req.PurchaserId)
	if err != nil {
		return nil, err
	}

	// 调用RPC服务
	_, err = l.svcCtx.PurchaseRPC.UpdateOrder(l.ctx, &pb.UpdateOrderReq{
		Id:           id,
		SupplierId:   supplierId,
		OrderDate:    req.OrderDate,
		ExpectedDate: req.ExpectedDate,
		TotalAmount:  req.TotalAmount,
		Status:       req.Status,
		PurchaserId:  purchaserId,
		ContractUrl:  req.ContractURL,
	})
	if err != nil {
		return nil, err
	}

	return &types.UpdateOrderResp{}, nil
}
