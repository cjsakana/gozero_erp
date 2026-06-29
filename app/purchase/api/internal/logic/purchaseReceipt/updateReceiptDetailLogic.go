package purchaseReceipt

import (
	"context"
	"erp/common/util"

	"erp/app/purchase/api/internal/svc"
	"erp/app/purchase/api/internal/types"
	"erp/app/purchase/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateReceiptDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateReceiptDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateReceiptDetailLogic {
	return &UpdateReceiptDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateReceiptDetailLogic) UpdateReceiptDetail(req *types.UpdateReceiptDetailReq) (resp *types.UpdateReceiptDetailResp, err error) {
	id, err := util.StringToInt64(req.Id)
	if err != nil {
		return nil, err
	}
	productId, err := util.StringToInt64(req.ProductId)
	if err != nil {
		return nil, err
	}
	batchId, err := util.StringToInt64(req.BatchId)
	if err != nil {
		return nil, err
	}

	// 调用RPC服务
	_, err = l.svcCtx.PurchaseRPC.UpdateReceiptDetail(l.ctx, &pb.UpdateReceiptDetailReq{
		Id:           id,
		ProductId:    productId,
		ProductName:  req.ProductName,
		CategoryType: req.CategoryType,
		Quantity:     req.Quantity,
		UnitPrice:    req.UnitPrice,
		Amount:       req.Amount,
		BatchId:      batchId,
	})
	if err != nil {
		return nil, err
	}

	return &types.UpdateReceiptDetailResp{}, nil
}
