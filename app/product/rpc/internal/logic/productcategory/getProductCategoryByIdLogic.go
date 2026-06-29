package productcategorylogic

import (
	"context"

	"erp/app/product/rpc/internal/code"
	"erp/app/product/rpc/internal/svc"
	"erp/app/product/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
)

type GetProductCategoryByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetProductCategoryByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProductCategoryByIdLogic {
	return &GetProductCategoryByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetProductCategoryByIdLogic) GetProductCategoryById(in *pb.GetProductCategoryByIdReq) (*pb.GetProductCategoryByIdResp, error) {
	one, err := l.svcCtx.ProductCategoryModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if err == sqlc.ErrNotFound {
			return nil, code.ProductCategoryNotFound
		}
		return nil, code.GetProductCategoryFail
	}

	return &pb.GetProductCategoryByIdResp{
		ProductCategory: &pb.ProductCategory{
			CategoryId:   one.Id,
			CategoryName: one.Name,
			ParentId:     one.ParentId,
		},
	}, nil
}
