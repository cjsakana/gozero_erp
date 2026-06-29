package productBatch

import (
	"context"
	"erp/app/product/rpc/client/product"
	"erp/common/util"

	"erp/app/product/api/internal/svc"
	"erp/app/product/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetProductBatchByIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetProductBatchByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProductBatchByIdLogic {
	return &GetProductBatchByIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetProductBatchByIdLogic) GetProductBatchById(req *types.GetProductBatchByIdRequest) (resp *types.GetProductBatchByIdResponse, err error) {
	id, err := util.StringToInt64(req.Id)
	if err != nil {
		return nil, err
	}

	ret, err := l.svcCtx.ProductRPC.GetProductBatchById(l.ctx, &product.GetProductBatchByIdReq{
		Id: id,
	})
	if err != nil {
		return nil, err
	}
	resp = &types.GetProductBatchByIdResponse{
		ProductBatch: types.ProductBatch{
			Id:             util.Int64ToString(ret.ProductBatch.Id),
			ProductId:      util.Int64ToString(ret.ProductBatch.ProductId),
			BatchNo:        ret.ProductBatch.BatchNo,
			ProductionDate: ret.ProductBatch.ProductionDate,
			CreatedAt:      ret.ProductBatch.CreatedAt,
		},
	}

	return
}
