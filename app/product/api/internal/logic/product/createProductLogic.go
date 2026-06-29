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

type CreateProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateProductLogic {
	return &CreateProductLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateProductLogic) CreateProduct(req *types.AddProductRequest) (resp *types.AddProductResponse, err error) {
	employeeId, err := util.GetInt64FromCtx(l.ctx, xtypes.EmployeeIdKey)
	if err != nil {
		return nil, err
	}

	categoryId, err := util.StringToInt64(req.CategoryId)
	if err != nil {
		return nil, err
	}

	ret, err := l.svcCtx.ProductRPC.AddProduct(l.ctx, &product.AddProductReq{
		ProductName:    req.ProductName,
		CategoryId:     categoryId,
		Specifications: req.Specifications,
		Unit:           req.Unit,
		PurchasePrice:  req.PurchasePrice,
		SellingPrice:   req.SellingPrice,
		IsActive:       req.IsActive,
		IsMaterial:     req.IsMaterial,
		CreatedBy:      employeeId,
	})
	if err != nil {
		return nil, err
	}
	resp = &types.AddProductResponse{
		Id: util.Int64ToString(ret.Id),
	}
	return
}
