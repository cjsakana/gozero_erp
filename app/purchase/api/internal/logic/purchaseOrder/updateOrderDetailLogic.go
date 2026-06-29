package purchaseOrder

import (
	"context"
	"erp/common/util"

	"erp/app/purchase/api/internal/svc"
	"erp/app/purchase/api/internal/types"
	"erp/app/purchase/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateOrderDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateOrderDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateOrderDetailLogic {
	return &UpdateOrderDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateOrderDetailLogic) UpdateOrderDetail(req *types.UpdateOrderDetailReq) (resp *types.UpdateOrderDetailResp, err error) {
	id, err := util.StringToInt64(req.Id)
	if err != nil {
		return nil, err
	}
	productId, err := util.StringToInt64(req.ProductId)
	if err != nil {
		return nil, err
	}

	// 调用RPC服务
	_, err = l.svcCtx.PurchaseRPC.UpdateOrderDetail(l.ctx, &pb.UpdateOrderDetailReq{
		Id:           id,
		ProductId:    productId,
		ProductName:  req.ProductName,
		CategoryType: req.CategoryType,
		Quantity:     req.Quantity,
		UnitPrice:    req.UnitPrice,
		Amount:       req.Amount,
		ReceivedQty:  req.ReceivedQty,
		Remark:       req.Remark,
	})
	if err != nil {
		return nil, err
	}

	return &types.UpdateOrderDetailResp{}, nil
}
