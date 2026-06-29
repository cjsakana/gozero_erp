package productBatch

import (
	"context"
	"erp/app/product/api/internal/svc"
	"erp/app/product/api/internal/types"
	"erp/app/product/rpc/pb"
	"erp/common/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchProductBatchLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchProductBatchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchProductBatchLogic {
	return &SearchProductBatchLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchProductBatchLogic) SearchProductBatch(req *types.SearchProductBatchRequest) (resp *types.SearchProductBatchResponse, err error) {
	productId, err := util.StringToInt64(req.ProductId)
	if err != nil {
		return nil, err
	}
	ret, err := l.svcCtx.ProductRPC.SearchProductBatch(l.ctx, &pb.SearchProductBatchReq{
		Page:      req.Page,
		Limit:     req.Limit,
		ProductId: productId,
		BatchNo:   req.BatchNo,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
	})
	if err != nil {
		return nil, err
	}
	resp = &types.SearchProductBatchResponse{
		Total: ret.Total,
	}
	for _, v := range ret.ProductBatch {
		resp.ProductBatch = append(resp.ProductBatch, &types.ProductBatch{
			Id:             util.Int64ToString(v.Id),
			ProductId:      util.Int64ToString(v.ProductId),
			BatchNo:        v.BatchNo,
			ProductionDate: v.ProductionDate,
			CreatedAt:      v.CreatedAt,
		})
	}

	return
}
