package productbatchlogic

import (
	"context"
	types2 "erp/app/product/rpc/internal/types"
	"time"

	"erp/app/product/rpc/internal/svc"
	"erp/app/product/rpc/pb"

	"erp/app/product/rpc/internal/code"
	"github.com/zeromicro/go-zero/core/logx"
)

type SearchProductBatchLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchProductBatchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchProductBatchLogic {
	return &SearchProductBatchLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SearchProductBatchLogic) SearchProductBatch(in *pb.SearchProductBatchReq) (*pb.SearchProductBatchResp, error) {
	productBatchs, total, err := l.svcCtx.ProductBatchModel.Search(l.ctx, &types2.SearchProductBatchParams{
		SearchCom: types2.SearchCom{
			Page:  in.Page,
			Limit: in.Limit,
		},
		ProductId: in.ProductId,
		BatchNo:   in.BatchNo,
		StartDate: time.Unix(in.StartDate, 0),
		EndDate:   time.Unix(in.EndDate, 0),
	})
	if err != nil {

		return nil, code.SearchBatchFail

	}

	resp := &pb.SearchProductBatchResp{
		Total: total,
	}
	for _, v := range productBatchs {
		resp.ProductBatch = append(resp.ProductBatch, &pb.ProductBatch{
			Id:             v.Id,
			ProductId:      v.ProductId,
			BatchNo:        v.BatchNo,
			ProductionDate: v.ProductionDate.Unix(),
			CreatedAt:      v.CreatedAt.Unix(),
		})
	}
	return resp, nil
}
