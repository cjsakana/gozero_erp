package product

import (
	"context"
	"erp/app/product/rpc/client/product"
	"erp/common/util"
	"erp/common/xtypes"

	"erp/app/product/api/internal/svc"
	"erp/app/product/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateProductLogic {
	return &UpdateProductLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateProductLogic) UpdateProduct(req *types.UpdateProductRequest) (resp *types.EmptyResponse, err error) {
	employeeId, err := util.GetInt64FromCtx(l.ctx, xtypes.EmployeeIdKey)
	if err != nil {
		return nil, err
	}
	id, err := util.StringToInt64(req.Id)
	if err != nil {
		return nil, err
	}
	categoryId, err := util.StringToInt64(req.CategoryId)
	if err != nil {
		return nil, err
	}
	_, err = l.svcCtx.ProductRPC.UpdateProduct(l.ctx, &product.UpdateProductReq{
		Id:             id,
		ProductName:    req.ProductName,
		CategoryId:     categoryId,
		Specifications: req.Specifications,
		Unit:           req.Unit,
		PurchasePrice:  req.PurchasePrice,
		SellingPrice:   req.SellingPrice,
		IsActive:       req.IsActive,
		IsMaterial:     req.IsMaterial,
		UpdatedBy:      employeeId,
	})
	if err != nil {
		return nil, err
	}

	return
}
