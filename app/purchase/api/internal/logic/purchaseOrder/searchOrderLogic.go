package purchaseOrder

import (
	"context"
	"erp/app/hr/rpc/client/employeedetail"
	"erp/app/product/rpc/client/product"
	"erp/app/purchase/api/internal/svc"
	"erp/app/purchase/api/internal/types"
	"erp/app/purchase/rpc/pb"
	supplierpb "erp/app/supplier/rpc/pb"
	"erp/app/supplier/rpc/supplier"
	"erp/common/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchOrderLogic {
	return &SearchOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchOrderLogic) SearchOrder(req *types.SearchOrderReq) (resp *types.SearchOrderResp, err error) {
	supplierId, err := util.StringToInt64(req.SupplierId)
	if err != nil {
		return nil, err
	}
	ret, err := l.svcCtx.PurchaseRPC.SearchOrder(l.ctx, &pb.SearchOrderReq{
		Page:       req.Page,
		Limit:      req.Limit,
		OrderNo:    req.OrderNo,
		SupplierId: supplierId,
		Status:     req.Status,
	})
	if err != nil {
		return nil, err
	}

	// 批量收集采购员ID和供应商ID
	//purchaserIds := make(map[int64]bool)
	//supplierIds := make(map[int64]bool)
	//for _, od := range ret.OrdersWithDetails {
	//	if od.Orders.PurchaserId > 0 {
	//		purchaserIds[od.Orders.PurchaserId] = true
	//	}
	//	if od.Orders.SupplierId > 0 {
	//		supplierIds[od.Orders.SupplierId] = true
	//	}
	//}
	//
	//// 批量查询采购员信息
	//purchaserMap := make(map[int64]struct{No, Name string})
	//for id := range purchaserIds {
	//	empResp, err := l.svcCtx.HrRPC.GetEmployeeDetailByEmployeeNo(l.ctx, &hrpb.GetEmployeeDetailByEmployeeNoReq{
	//		EmployeeId: id,
	//	})
	//	if err != nil {
	//		logx.Errorw("获取采购员信息失败", logx.Field("purchaser_id", id), logx.Field("error", err))
	//	} else if empResp.EmployeeNonSensitiveDetail != nil {
	//		purchaserMap[id] = struct{No, Name string}{
	//			No:   empResp.EmployeeNonSensitiveDetail.EmployeeNo,
	//			Name: empResp.EmployeeNonSensitiveDetail.Name,
	//		}
	//	}
	//}
	//
	//// 批量查询供应商信息
	//supplierMap := make(map[int64]string)
	//for id := range supplierIds {
	//	supResp, err := l.svcCtx.SupplierRPC.GetSupplierById(l.ctx, &supplierpb.GetSupplierByIdReq{
	//		Id: id,
	//	})
	//	if err != nil {
	//		logx.Errorw("获取供应商信息失败", logx.Field("supplier_id", id), logx.Field("error", err))
	//	} else if supResp.Supplier != nil {
	//		supplierMap[id] = supResp.Supplier.Name
	//	}
	//}

	employeeMap := make(map[int64]*employeedetail.EmployeeNonSensitiveDetail)
	supplierMap := make(map[int64]*supplier.Supplier)
	productMap := make(map[int64]*product.Product)

	// 组装响应列表
	list := make([]*types.PurchaseOrderWithDetails, 0, ret.Total)
	for _, od := range ret.OrdersWithDetails {

		// 供应商信息
		if _, ok := supplierMap[od.Orders.SupplierId]; !ok {
			ret, err := l.svcCtx.SupplierRPC.GetSupplierById(l.ctx, &supplierpb.GetSupplierByIdReq{
				Id: od.Orders.SupplierId,
			})
			if err != nil {
				return nil, err
			}
			supplierMap[od.Orders.SupplierId] = ret.Supplier
		}

		// 采购人
		if _, ok := employeeMap[od.Orders.PurchaserId]; !ok {
			empResp, err := l.svcCtx.HrRPC.GetEmployeeDetailById(l.ctx, &employeedetail.GetEmployeeDetailByIdReq{
				Id: od.Orders.PurchaserId,
			})
			if err != nil {
				return nil, err
			}
			employeeMap[od.Orders.PurchaserId] = empResp.EmployeeNonSensitiveDetail
		}

		order := types.PurchaseOrder{
			Id:            util.Int64ToString(od.Orders.Id),
			OrderNo:       od.Orders.OrderNo,
			SupplierId:    util.Int64ToString(od.Orders.SupplierId),
			SupplierCode:  supplierMap[od.Orders.SupplierId].Code,
			SupplierName:  supplierMap[od.Orders.SupplierId].Name,
			OrderDate:     od.Orders.OrderDate,
			ExpectedDate:  od.Orders.ExpectedDate,
			TotalAmount:   od.Orders.TotalAmount,
			Status:        od.Orders.Status,
			PurchaserId:   util.Int64ToString(od.Orders.PurchaserId),
			PurchaserNo:   employeeMap[od.Orders.PurchaserId].EmployeeNo,
			PurchaserName: employeeMap[od.Orders.PurchaserId].Name,
			ContractURL:   od.Orders.ContractUrl,
			CreatedAt:     od.Orders.CreatedAt,
			UpdatedAt:     od.Orders.UpdatedAt,
		}

		tOd := &types.PurchaseOrderWithDetails{
			Order: order,
			Total: od.Total,
			Details: func() []*types.PurchaseOrderDetail {
				details := make([]*types.PurchaseOrderDetail, 0, od.Total)
				for _, d := range od.Details {

					if _, ok := productMap[d.ProductId]; !ok {
						prod, err := l.svcCtx.ProductRPC.GetProductById(l.ctx, &product.GetProductByIdReq{
							Id: d.ProductId,
						})
						if err != nil {
							return nil
						}
						productMap[d.ProductId] = prod.Product
					}

					details = append(details, &types.PurchaseOrderDetail{
						Id:           util.Int64ToString(d.Id),
						OrderId:      util.Int64ToString(d.OrderId),
						ProductId:    util.Int64ToString(d.ProductId),
						ProductNo:    productMap[d.ProductId].ProductNo,
						ProductName:  d.ProductName,
						CategoryType: d.CategoryType,
						Quantity:     d.Quantity,
						UnitPrice:    d.UnitPrice,
						Amount:       d.Amount,
						ReceivedQty:  d.ReceivedQty,
						Remark:       d.Remark,
					})
				}
				return details
			}(),
		}
		list = append(list, tOd)
	}

	resp = &types.SearchOrderResp{
		Total:             ret.Total,
		OrdersWithDetails: list,
	}
	return
}
