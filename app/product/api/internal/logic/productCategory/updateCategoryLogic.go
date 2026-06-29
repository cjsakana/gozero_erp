package productCategory

import (
	"context"
	"erp/app/product/rpc/pb"
	"erp/common/util"

	"erp/app/product/api/internal/svc"
	"erp/app/product/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateCategoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCategoryLogic {
	return &UpdateCategoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateCategoryLogic) UpdateCategory(req *types.UpdateProductCategoryRequest) (resp *types.EmptyResponse, err error) {
	categoryId, err := util.StringToInt64(req.CategoryId)
	if err != nil {
		return nil, err
	}
	parentId, err := util.StringToInt64(req.ParentId)
	if err != nil {
		return nil, err
	}
	_, err = l.svcCtx.ProductRPC.UpdateProductCategory(l.ctx, &pb.UpdateProductCategoryReq{
		CategoryId:   categoryId,
		CategoryName: req.CategoryName,
		ParentId:     parentId,
	})
	if err != nil {
		return nil, err
	}

	return
}
