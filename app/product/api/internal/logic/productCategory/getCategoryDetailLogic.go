package productCategory

import (
	"context"
	"erp/app/product/rpc/pb"
	"erp/common/util"

	"erp/app/product/api/internal/svc"
	"erp/app/product/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCategoryDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCategoryDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCategoryDetailLogic {
	return &GetCategoryDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCategoryDetailLogic) GetCategoryDetail(req *types.GetProductCategoryByIdRequest) (resp *types.GetProductCategoryByIdResponse, err error) {
	id, err := util.StringToInt64(req.Id)
	if err != nil {
		return nil, err
	}
	one, err := l.svcCtx.ProductRPC.GetProductCategoryById(l.ctx, &pb.GetProductCategoryByIdReq{
		Id: id,
	})
	if err != nil {
		return nil, err
	}
	resp = &types.GetProductCategoryByIdResponse{
		ProductCategory: types.ProductCategory{
			CategoryId:   util.Int64ToString(one.ProductCategory.CategoryId),
			CategoryName: one.ProductCategory.CategoryName,
			ParentId:     util.Int64ToString(one.ProductCategory.ParentId),
			Children:     nil,
		},
	}

	return
}
