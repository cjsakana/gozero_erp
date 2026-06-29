package productCategory

import (
	"context"
	"erp/app/product/rpc/pb"
	"erp/common/util"
	"fmt"

	"erp/app/product/api/internal/svc"
	"erp/app/product/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCategoryListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCategoryListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCategoryListLogic {
	return &GetCategoryListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCategoryListLogic) GetCategoryList(req *types.SearchProductCategoryRequest) (resp *types.SearchProductCategoryResponse, err error) {
	categoryId, err := util.StringToInt64(req.CategoryId)
	if err != nil {
		return nil, err
	}
	parentId, err := util.StringToInt64(req.ParentId)
	if err != nil {
		return nil, err
	}
	ret, err := l.svcCtx.ProductRPC.SearchProductCategory(l.ctx, &pb.SearchProductCategoryReq{
		Page:         req.Page,
		Limit:        req.Limit,
		CategoryId:   categoryId,
		CategoryName: req.CategoryName,
		ParentId:     parentId,
	})
	if err != nil {
		return nil, err
	}

	resp = &types.SearchProductCategoryResponse{
		Total: ret.Total,
	}

	for i, category := range ret.ProductCategory {
		fmt.Println(i, category)
		resp.Categories = append(resp.Categories, &types.ProductCategory{
			CategoryId:   util.Int64ToString(category.CategoryId),
			CategoryName: category.CategoryName,
			ParentId:     util.Int64ToString(category.ParentId),
		})
	}

	return
}
