package productbatchlogic

import (
	"context"

	"erp/app/product/rpc/internal/svc"
	"erp/app/product/rpc/pb"

	"erp/app/product/rpc/internal/code"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
)

type GetProductBatchByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetProductBatchByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProductBatchByIdLogic {
	return &GetProductBatchByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetProductBatchByIdLogic) GetProductBatchById(in *pb.GetProductBatchByIdReq) (*pb.GetProductBatchByIdResp, error) {
	one, err := l.svcCtx.ProductBatchModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if err == sqlc.ErrNotFound {
			return nil, code.ProductBatchNotFound
		}
		return nil, code.GetProductBatchFail

	}

	return &pb.GetProductBatchByIdResp{
		ProductBatch: &pb.ProductBatch{
			Id:             one.Id,
			ProductId:      one.ProductId,
			BatchNo:        one.BatchNo,
			ProductionDate: one.ProductionDate.Unix(),
			CreatedAt:      one.CreatedAt.Unix(),
		},
	}, nil
}
