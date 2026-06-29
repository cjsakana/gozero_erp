package productbatchlogic

import (
	"context"
	"erp/app/product/rpc/internal/model"
	"erp/common/util"
	"github.com/zeromicro/go-zero/core/logx"
	"time"

	"erp/app/product/rpc/internal/svc"
	"erp/app/product/rpc/pb"

	"erp/app/product/rpc/internal/code"
)

type AddProductBatchLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddProductBatchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddProductBatchLogic {
	return &AddProductBatchLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddProductBatchLogic) AddProductBatch(in *pb.AddProductBatchReq) (*pb.AddProductBatchResp, error) {
	id := util.GenerateSnowflake()
	_, err := l.svcCtx.ProductBatchModel.Insert(l.ctx, &model.ProductBatch{
		Id:             id,
		ProductId:      in.ProductId,
		BatchNo:        in.BatchNo,
		ProductionDate: time.Unix(in.ProductionDate, 0),
	})
	if err != nil {

		return nil, code.AddBatchFail

	}

	return &pb.AddProductBatchResp{
		Id: id,
	}, nil
}
