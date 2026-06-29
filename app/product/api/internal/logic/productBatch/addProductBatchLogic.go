package productBatch

import (
	"context"
	"erp/app/product/rpc/client/product"
	"erp/common/util"

	"erp/app/product/api/internal/svc"
	"erp/app/product/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddProductBatchLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddProductBatchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddProductBatchLogic {
	return &AddProductBatchLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddProductBatchLogic) AddProductBatch(req *types.AddProductBatchRequest) (resp *types.AddProductBatchResponse, err error) {
	productId, err := util.StringToInt64(req.ProductId)
	if err != nil {
		return nil, err
	}
	ret, err := l.svcCtx.ProductRPC.AddProductBatch(l.ctx, &product.AddProductBatchReq{
		ProductId:      productId,
		BatchNo:        req.BatchNo,
		ProductionDate: req.ProductionDate,
	})
	if err != nil {
		return nil, err
	}
	resp = &types.AddProductBatchResponse{
		Id: util.Int64ToString(ret.Id),
	}
	return
}
