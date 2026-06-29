package productcategorylogic

import (
	"context"
	"erp/app/product/rpc/internal/code"
	"erp/app/product/rpc/internal/model"

	"erp/app/product/rpc/internal/svc"
	"erp/app/product/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateProductCategoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateProductCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateProductCategoryLogic {
	return &UpdateProductCategoryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateProductCategoryLogic) UpdateProductCategory(in *pb.UpdateProductCategoryReq) (*pb.UpdateProductCategoryResp, error) {
	err := l.svcCtx.ProductCategoryModel.XUpdate(l.ctx, &model.ProductCategory{
		Id:       in.CategoryId,
		Name:     in.CategoryName,
		ParentId: in.ParentId,
	})
	if err != nil {
		return nil, code.UpdateCategoryFail
	}

	return &pb.UpdateProductCategoryResp{}, nil
}
