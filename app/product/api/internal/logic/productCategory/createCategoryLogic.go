package productCategory

import (
	"context"
	"erp/app/product/rpc/pb"
	"erp/common/util"

	"erp/app/product/api/internal/svc"
	"erp/app/product/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateCategoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCategoryLogic {
	return &CreateCategoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateCategoryLogic) CreateCategory(req *types.AddProductCategoryRequest) (resp *types.AddProductCategoryResponse, err error) {
	parentId, err := util.StringToInt64(req.ParentId)
	if err != nil {
		return nil, err
	}
	ret, err := l.svcCtx.ProductRPC.AddProductCategory(l.ctx, &pb.AddProductCategoryReq{
		CategoryName: req.CategoryName,
		ParentId:     parentId,
	})
	if err != nil {
		return nil, err
	}
	resp = &types.AddProductCategoryResponse{
		Id: util.Int64ToString(ret.CategoryId),
	}

	return
}
