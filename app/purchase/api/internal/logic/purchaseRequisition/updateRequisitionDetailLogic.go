package purchaseRequisition

import (
	"context"
	"erp/common/util"

	"erp/app/purchase/api/internal/svc"
	"erp/app/purchase/api/internal/types"
	"erp/app/purchase/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateRequisitionDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateRequisitionDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateRequisitionDetailLogic {
	return &UpdateRequisitionDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateRequisitionDetailLogic) UpdateRequisitionDetail(req *types.UpdateRequisitionDetailReq) (resp *types.UpdateRequisitionDetailResp, err error) {
	id, err := util.StringToInt64(req.Id)
	if err != nil {
		return nil, err
	}
	productId, err := util.StringToInt64(req.ProductId)
	if err != nil {
		return nil, err
	}

	// 调用RPC服务
	_, err = l.svcCtx.PurchaseRPC.UpdateRequisitionDetail(l.ctx, &pb.UpdateRequisitionDetailReq{
		Id:           id,
		ProductId:    productId,
		ProductName:  req.ProductName,
		CategoryType: req.CategoryType,
		Quantity:     req.Quantity,
		UnitPrice:    req.UnitPrice,
		Amount:       req.Amount,
		Remark:       req.Remark,
	})
	if err != nil {
		return nil, err
	}

	return &types.UpdateRequisitionDetailResp{}, nil
}
