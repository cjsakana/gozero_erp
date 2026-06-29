package productcategorylogic

import (
	"context"

	"erp/app/product/rpc/internal/code"
	"erp/app/product/rpc/internal/svc"
	"erp/app/product/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DelProductCategoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDelProductCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelProductCategoryLogic {
	return &DelProductCategoryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DelProductCategoryLogic) DelProductCategory(in *pb.DelProductCategoryReq) (*pb.DelProductCategoryResp, error) {
	err := l.svcCtx.ProductCategoryModel.XDelete(l.ctx, in.Id)
	if err != nil {
		return nil, code.DeleteCategoryFail
	}

	return &pb.DelProductCategoryResp{}, nil
}
