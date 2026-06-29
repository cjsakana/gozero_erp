package productcategorylogic

import (
	"context"
	"erp/app/product/rpc/internal/model"
	"erp/common/util"
	"github.com/zeromicro/go-zero/core/logx"

	"erp/app/product/rpc/internal/svc"
	"erp/app/product/rpc/pb"

	"erp/app/product/rpc/internal/code"
)

type AddProductCategoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddProductCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddProductCategoryLogic {
	return &AddProductCategoryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// -----------------------productCategory-----------------------
func (l *AddProductCategoryLogic) AddProductCategory(in *pb.AddProductCategoryReq) (*pb.AddProductCategoryResp, error) {
	id := util.GenerateSnowflake()
	_, err := l.svcCtx.ProductCategoryModel.Insert(l.ctx, &model.ProductCategory{
		Id:       id,
		Name:     in.CategoryName,
		ParentId: in.ParentId,
	})
	if err != nil {

		return nil, code.AddCategoryFail

	}

	return &pb.AddProductCategoryResp{
		CategoryId: id,
	}, nil
}
