package salesDelivery

import (
	"context"
	"erp/app/hr/rpc/client/employeedetail"
	"erp/app/inventory/rpc/client/inventory"
	"erp/app/product/rpc/client/product"
	"erp/app/product/rpc/client/productbatch"
	"erp/app/sale/api/internal/svc"
	"erp/app/sale/api/internal/types"
	"erp/app/sale/rpc/pb"
	"erp/common/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchSalesDeliveryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchSalesDeliveryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchSalesDeliveryLogic {
	return &SearchSalesDeliveryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchSalesDeliveryLogic) SearchSalesDelivery(req *types.SearchSalesDeliveryReq) (resp *types.SearchSalesDeliveryResp, err error) {
	orderId, err := util.StringToInt64(req.OrderId)
	if err != nil {
		return nil, err
	}
	warehouseId, err := util.StringToInt64(req.WarehouseId)
	if err != nil {
		return nil, err
	}

	ret, err := l.svcCtx.SaleRPC.SearchSalesDelivery(l.ctx, &pb.SearchSalesDeliveryReq{
		Page:         req.Page,
		Limit:        req.Limit,
		DeliveryNo:   req.DeliveryNo,
		OrderId:      orderId,
		WarehouseId:  warehouseId,
		DeliveryDate: req.DeliveryDate,
		Status:       req.Status,
	})
	if err != nil {
		return nil, err
	}

	warehouseMap := make(map[int64]*inventory.WarehouseDetail)
	productBatchMap := make(map[int64]*productbatch.ProductBatch)
	productMap := make(map[int64]*product.Product)
	employeeMap := make(map[int64]*employeedetail.EmployeeNonSensitiveDetail)

	resp = &types.SearchSalesDeliveryResp{
		Total: ret.Total,
		DeliveryWithDetails: func() []types.DeliveryWithDetails {
			list := make([]types.DeliveryWithDetails, 0)
			for _, d := range ret.DeliveryWithDetails {

				if _, ok := warehouseMap[d.SalesDelivery.WarehouseId]; !ok {
					w, e := l.svcCtx.InventoryRPC.GetWarehouseById(l.ctx, &inventory.GetWarehouseByIdReq{
						Id: d.SalesDelivery.WarehouseId,
					})
					if e != nil {
						return nil
					}
					warehouseMap[d.SalesDelivery.WarehouseId] = w.Warehouse
				}
				if _, ok := employeeMap[d.SalesDelivery.CreatedBy]; !ok {
					empResp, err := l.svcCtx.HrRPC.GetEmployeeDetailById(l.ctx, &employeedetail.GetEmployeeDetailByIdReq{
						Id: d.SalesDelivery.CreatedBy,
					})
					if err != nil {
						return nil
					}
					employeeMap[d.SalesDelivery.CreatedBy] = empResp.EmployeeNonSensitiveDetail
				}

				list = append(list, types.DeliveryWithDetails{
					SalesDelivery: types.SalesDelivery{
						Id:            util.Int64ToString(d.SalesDelivery.Id),
						DeliveryNo:    d.SalesDelivery.DeliveryNo,
						OrderId:       util.Int64ToString(d.SalesDelivery.OrderId),
						OrderNo:       d.SalesDelivery.OrderNo,
						WarehouseId:   util.Int64ToString(d.SalesDelivery.WarehouseId),
						WarehouseNo:   warehouseMap[d.SalesDelivery.WarehouseId].No,
						WarehouseName: warehouseMap[d.SalesDelivery.WarehouseId].Name,
						DeliveryDate:  d.SalesDelivery.DeliveryDate,
						TotalQuantity: d.SalesDelivery.TotalQuantity,
						TotalAmount:   d.SalesDelivery.TotalAmount,
						Status:        d.SalesDelivery.Status,
						CreatedBy:     util.Int64ToString(d.SalesDelivery.CreatedBy),
						CreatedByNo:   employeeMap[d.SalesDelivery.CreatedBy].EmployeeNo,
						CreatedByName: employeeMap[d.SalesDelivery.CreatedBy].Name,
						CreatedAt:     d.SalesDelivery.CreatedAt,
					},
					Total: d.Total,
					SalesDeliveryDetail: func() []types.SalesDeliveryDetail {
						detailList := make([]types.SalesDeliveryDetail, 0)
						for _, s := range d.SalesDeliveryDetail {

							if _, ok := productBatchMap[s.BatchId]; !ok {
								productBatch, err := l.svcCtx.ProductRPC.GetProductBatchById(l.ctx, &productbatch.GetProductBatchByIdReq{
									Id: s.BatchId,
								})
								if err != nil {
									return nil
								}
								productBatchMap[s.BatchId] = productBatch.ProductBatch
							}
							if _, ok := productMap[s.ProductId]; !ok {
								prod, err := l.svcCtx.ProductRPC.GetProductById(l.ctx, &product.GetProductByIdReq{
									Id: s.ProductId,
								})
								if err != nil {
									return nil
								}
								productMap[s.ProductId] = prod.Product
							}

							detailList = append(detailList, types.SalesDeliveryDetail{
								Id:             util.Int64ToString(s.Id),
								DeliveryId:     util.Int64ToString(s.DeliveryId),
								ProductId:      util.Int64ToString(s.ProductId),
								ProductNo:      productMap[s.ProductId].ProductNo,
								ProductName:    s.ProductName,
								Unit:           s.Unit,
								Quantity:       s.Quantity,
								UnitPrice:      s.UnitPrice,
								Amount:         s.Amount,
								BatchId:        util.Int64ToString(s.BatchId),
								BatchNo:        productBatchMap[s.BatchId].BatchNo,
								ProductionDate: productBatchMap[s.BatchId].ProductionDate,
							})
						}
						return detailList
					}(),
				})
			}
			return list
		}(),
	}
	return
}
