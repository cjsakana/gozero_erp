package product

import (
	"context"
	"erp/app/hr/rpc/client/employeedetail"
	"erp/app/inventory/rpc/client/inventory"
	"erp/app/product/api/internal/svc"
	"erp/app/product/api/internal/types"
	"erp/app/product/rpc/client/product"
	"erp/common/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetProductDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetProductDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProductDetailLogic {
	return &GetProductDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetProductDetailLogic) GetProductDetail(req *types.GetProductByIdRequest) (resp *types.GetProductByIdResponse, err error) {
	id, err := util.StringToInt64(req.Id)
	if err != nil {
		return nil, err
	}

	one, err := l.svcCtx.ProductRPC.GetProductById(l.ctx, &product.GetProductByIdReq{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	// 查询创建人和更新人信息
	var createdByNo, createdByName, updatedByNo, updatedByName string
	if one.Product.CreatedBy > 0 {
		eRet, err := l.svcCtx.HrRPC.EmployeeDetailZrpcClient.GetEmployeeDetailById(l.ctx, &employeedetail.GetEmployeeDetailByIdReq{
			Id: one.Product.CreatedBy,
		})
		if err == nil {
			createdByNo = eRet.EmployeeNonSensitiveDetail.EmployeeNo
			createdByName = eRet.EmployeeNonSensitiveDetail.Name
		}
	}
	if one.Product.UpdatedBy > 0 {
		eRet, err := l.svcCtx.HrRPC.EmployeeDetailZrpcClient.GetEmployeeDetailById(l.ctx, &employeedetail.GetEmployeeDetailByIdReq{
			Id: one.Product.CreatedBy,
		})
		if err == nil {
			updatedByNo = eRet.EmployeeNonSensitiveDetail.EmployeeNo
			updatedByName = eRet.EmployeeNonSensitiveDetail.Name
		}
	}
	IRet, err := l.svcCtx.InventoryRPC.SearchInventory(l.ctx, &inventory.SearchInventoryReq{
		Limit:     -1,
		ProductId: id,
	})
	if err != nil {
		return nil, err
	}

	totalStock := 0.0
	totalLocked := 0.0

	for _, v := range IRet.Inventory {
		totalStock += v.CurrentStock
		totalLocked += v.LockedStock
	}

	resp = &types.GetProductByIdResponse{
		Product: types.Product{
			Id:             util.Int64ToString(one.Product.Id),
			ProductNo:      one.Product.ProductNo,
			ProductName:    one.Product.ProductName,
			CategoryId:     util.Int64ToString(one.Product.CategoryId),
			Specifications: one.Product.Specifications,
			Unit:           one.Product.Unit,
			PurchasePrice:  one.Product.PurchasePrice,
			SellingPrice:   one.Product.SellingPrice,
			TotalStock:     totalStock,
			TotalLocked:    totalLocked,
			IsActive:       one.Product.IsActive,
			IsMaterial:     one.Product.IsMaterial,
			CreatedAt:      one.Product.CreatedAt,
			CreatedBy:      util.Int64ToString(one.Product.CreatedBy),
			CreatedByNo:    createdByNo,
			CreatedByName:  createdByName,
			UpdatedAt:      one.Product.UpdatedAt,
			UpdatedBy:      util.Int64ToString(one.Product.UpdatedBy),
			UpdatedByNo:    updatedByNo,
			UpdatedByName:  updatedByName,
		},
	}
	return
}
