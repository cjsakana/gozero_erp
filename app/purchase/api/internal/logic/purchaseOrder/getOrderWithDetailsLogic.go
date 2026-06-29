package purchaseOrder

import (
	"context"
	pb2 "erp/app/hr/rpc/pb"
	"erp/app/product/rpc/client/product"
	"erp/app/purchase/api/internal/svc"
	"erp/app/purchase/api/internal/types"
	"erp/app/purchase/rpc/pb"
	pb3 "erp/app/supplier/rpc/pb"
	"erp/common/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOrderWithDetailsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetOrderWithDetailsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrderWithDetailsLogic {
	return &GetOrderWithDetailsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetOrderWithDetailsLogic) GetOrderWithDetails(req *types.GetOrderWithDetailsReq) (resp *types.GetOrderWithDetailsResp, err error) {
	orderId, err := util.StringToInt64(req.OrderId)
	if err != nil {
		return nil, err
	}
	ret, err := l.svcCtx.PurchaseRPC.GetOrderWithDetails(l.ctx, &pb.GetOrderWithDetailsReq{
		OrderId: orderId,
	})
	if err != nil {
		return nil, err
	}

	var purchaserNo, purchaserName string
	empResp, err := l.svcCtx.HrRPC.GetEmployeeDetailById(l.ctx, &pb2.GetEmployeeDetailByIdReq{
		Id: ret.Order.PurchaserId,
	})
	if err != nil {
		logx.Errorw("获取采购员信息失败", logx.Field("purchaser_id", ret.Order.PurchaserId), logx.Field("error", err))
		return nil, err
	}
	purchaserNo = empResp.EmployeeNonSensitiveDetail.EmployeeNo
	purchaserName = empResp.EmployeeNonSensitiveDetail.Name

	var supplierCode, supplierName string
	supplierResp, _ := l.svcCtx.SupplierRPC.GetSupplierById(l.ctx, &pb3.GetSupplierByIdReq{
		Id: ret.Order.SupplierId,
	})
	supplierCode = supplierResp.Supplier.Code
	supplierName = supplierResp.Supplier.Name

	resp = &types.GetOrderWithDetailsResp{
		Order: types.PurchaseOrder{
			Id:            util.Int64ToString(ret.Order.Id),
			OrderNo:       ret.Order.OrderNo,
			SupplierId:    util.Int64ToString(ret.Order.SupplierId),
			SupplierCode:  supplierCode,
			SupplierName:  supplierName,
			OrderDate:     ret.Order.OrderDate,
			ExpectedDate:  ret.Order.ExpectedDate,
			TotalAmount:   ret.Order.TotalAmount,
			Status:        ret.Order.Status,
			PurchaserId:   util.Int64ToString(ret.Order.PurchaserId),
			PurchaserNo:   purchaserNo,
			PurchaserName: purchaserName,
			ContractURL:   ret.Order.ContractUrl,
			CreatedAt:     ret.Order.CreatedAt,
			UpdatedAt:     ret.Order.UpdatedAt,
		},
		Details: func() []*types.PurchaseOrderDetail {
			list := make([]*types.PurchaseOrderDetail, 0, len(ret.Details))
			productMap := make(map[int64]*product.Product)

			for _, d := range ret.Details {

				detail := &types.PurchaseOrderDetail{
					Id:           util.Int64ToString(d.Id),
					OrderId:      util.Int64ToString(d.OrderId),
					ProductId:    util.Int64ToString(d.ProductId),
					ProductName:  d.ProductName,
					CategoryType: d.CategoryType,
					Quantity:     d.Quantity,
					UnitPrice:    d.UnitPrice,
					Amount:       d.Amount,
					ReceivedQty:  d.ReceivedQty,
					Remark:       d.Remark,
				}
				if d.ProductId > 0 {
					if _, ok := productMap[d.ProductId]; !ok {
						prod, err := l.svcCtx.ProductRPC.GetProductById(l.ctx, &product.GetProductByIdReq{
							Id: d.ProductId,
						})
						if err != nil {
							return nil
						}
						productMap[d.ProductId] = prod.Product
					}
					detail.ProductNo = productMap[d.ProductId].ProductNo
				}
				list = append(list, detail)
			}
			return list
		}(),
	}
	return
}
