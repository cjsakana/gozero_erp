package productcategorylogic

import (
	"context"
	types2 "erp/app/product/rpc/internal/types"

	"erp/app/product/rpc/internal/svc"
	"erp/app/product/rpc/pb"

	"erp/app/product/rpc/internal/code"
	"github.com/zeromicro/go-zero/core/logx"
)

type SearchProductCategoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchProductCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchProductCategoryLogic {
	return &SearchProductCategoryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SearchProductCategoryLogic) SearchProductCategory(in *pb.SearchProductCategoryReq) (*pb.SearchProductCategoryResp, error) {
	categories, total, err := l.svcCtx.ProductCategoryModel.Search(l.ctx, &types2.SearchProductCategoryParams{
		SearchCom: types2.SearchCom{
			Page:  in.Page,
			Limit: in.Limit,
		},
		CategoryName: in.CategoryName,
		ParentId:     in.ParentId,
	})
	if err != nil {

		return nil, code.GetProductCategoryFail

	}
	var list []*pb.ProductCategory
	for _, one := range categories {
		list = append(list, &pb.ProductCategory{
			CategoryId:   one.Id,
			CategoryName: one.Name,
			ParentId:     one.ParentId,
		})
	}

	return &pb.SearchProductCategoryResp{
		Total:           total,
		ProductCategory: list,
	}, nil
}
