package product

import (
	"context"
	"erp/app/hr/rpc/client/employeedetail"
	"erp/app/inventory/rpc/client/inventory"
	"erp/app/product/api/internal/svc"
	"erp/app/product/api/internal/types"
	"erp/app/product/rpc/pb"
	"erp/common/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetProductListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetProductListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProductListLogic {
	return &GetProductListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetProductListLogic) GetProductList(req *types.SearchProductRequest) (resp *types.SearchProductResponse, err error) {
	categoryId, err := util.StringToInt64(req.CategoryId)
	if err != nil {
		return nil, err
	}

	ret, err := l.svcCtx.ProductRPC.SearchProduct(l.ctx, &pb.SearchProductReq{
		Page:        req.Page,
		Limit:       req.Limit,
		ProductNo:   req.ProductNo,
		ProductName: req.ProductName,
		CategoryId:  categoryId,
		IsActive:    req.IsActive,
		IsMaterial:  req.IsMaterial,
	})
	if err != nil {
		return nil, err
	}

	resp = &types.SearchProductResponse{
		Total: ret.Total,
	}

	employeeMap := make(map[int64]*employeedetail.EmployeeNonSensitiveDetail)

	for _, product := range ret.Product {
		if _, ok := employeeMap[product.CreatedBy]; !ok {
			employeeDetail, err := l.svcCtx.HrRPC.EmployeeDetailZrpcClient.GetEmployeeDetailById(l.ctx, &employeedetail.GetEmployeeDetailByIdReq{
				Id: product.CreatedBy,
			})
			if err != nil {
				logx.Errorf("查询员工信息失败: employeeId=%d, err=%v", product.CreatedBy, err)
			}
			employeeMap[product.CreatedBy] = employeeDetail.EmployeeNonSensitiveDetail
		}
		if _, ok := employeeMap[product.UpdatedBy]; !ok {
			employeeDetail, err := l.svcCtx.HrRPC.EmployeeDetailZrpcClient.GetEmployeeDetailById(l.ctx, &employeedetail.GetEmployeeDetailByIdReq{
				Id: product.UpdatedBy,
			})
			if err != nil {
				logx.Errorf("查询员工信息失败: employeeId=%d, err=%v", product.UpdatedBy, err)
			}
			employeeMap[product.UpdatedBy] = employeeDetail.EmployeeNonSensitiveDetail
		}

		totalStock, totalLocked, err := l.GetSummary(product.Id)
		if err != nil {
			return nil, err
		}

		resp.Products = append(resp.Products, types.Product{
			Id:             util.Int64ToString(product.Id),
			ProductNo:      product.ProductNo,
			ProductName:    product.ProductName,
			CategoryId:     util.Int64ToString(product.CategoryId),
			Specifications: product.Specifications,
			Unit:           product.Unit,
			PurchasePrice:  product.PurchasePrice,
			SellingPrice:   product.SellingPrice,
			TotalStock:     totalStock,
			TotalLocked:    totalLocked,
			IsActive:       product.IsActive,
			IsMaterial:     product.IsMaterial,
			CreatedAt:      product.CreatedAt,
			CreatedBy:      util.Int64ToString(product.CreatedBy),
			CreatedByNo:    employeeMap[product.CreatedBy].EmployeeNo,
			CreatedByName:  employeeMap[product.CreatedBy].Name,
			UpdatedAt:      product.UpdatedAt,
			UpdatedBy:      util.Int64ToString(product.UpdatedBy),
			UpdatedByNo:    employeeMap[product.UpdatedBy].EmployeeNo,
			UpdatedByName:  employeeMap[product.UpdatedBy].Name,
		})
	}

	return
}

func (l *GetProductListLogic) GetSummary(productId int64) (float64, float64, error) {
	ret, err := l.svcCtx.InventoryRPC.SearchInventory(l.ctx, &inventory.SearchInventoryReq{
		Limit:     -1,
		ProductId: productId,
	})
	if err != nil {
		return 0, 0, err
	}

	totalStock := 0.0
	totalLocked := 0.0

	for _, v := range ret.Inventory {
		totalStock += v.CurrentStock
		totalLocked += v.LockedStock
	}
	return totalStock, totalLocked, nil
}
