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

type GetSalesDeliveryByIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSalesDeliveryByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSalesDeliveryByIdLogic {
	return &GetSalesDeliveryByIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSalesDeliveryByIdLogic) GetSalesDeliveryById(req *types.GetSalesDeliveryByIdReq) (resp *types.GetSalesDeliveryByIdResp, err error) {
	id, err := util.StringToInt64(req.Id)
	if err != nil {
		return nil, err
	}
	ret, err := l.svcCtx.SaleRPC.GetSalesDeliveryById(l.ctx, &pb.GetSalesDeliveryByIdReq{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	// 获取操作人信息
	var createdByNo, createdByName string
	if ret.DeliveryWithDetails.SalesDelivery.CreatedBy > 0 {
		empResp, err := l.svcCtx.HrRPC.GetEmployeeDetailById(l.ctx, &employeedetail.GetEmployeeDetailByIdReq{
			Id: ret.DeliveryWithDetails.SalesDelivery.CreatedBy,
		})
		if err != nil {
			logx.Errorw("获取操作人信息失败", logx.Field("created_by", ret.DeliveryWithDetails.SalesDelivery.CreatedBy), logx.Field("error", err))
		} else if empResp.EmployeeNonSensitiveDetail != nil {
			createdByNo = empResp.EmployeeNonSensitiveDetail.EmployeeNo
			createdByName = empResp.EmployeeNonSensitiveDetail.Name
		}
	}

	var warehouseNo, warehouseName string
	if ret.DeliveryWithDetails.SalesDelivery.WarehouseId > 0 {
		w, e := l.svcCtx.InventoryRPC.GetWarehouseById(l.ctx, &inventory.GetWarehouseByIdReq{
			Id: ret.DeliveryWithDetails.SalesDelivery.WarehouseId,
		})
		if e != nil {
			return nil, e
		}
		warehouseNo = w.Warehouse.No
		warehouseName = w.Warehouse.Name
	}

	resp = &types.GetSalesDeliveryByIdResp{
		DeliveryWithDetails: types.DeliveryWithDetails{
			SalesDelivery: types.SalesDelivery{
				Id:            util.Int64ToString(ret.DeliveryWithDetails.SalesDelivery.Id),
				DeliveryNo:    ret.DeliveryWithDetails.SalesDelivery.DeliveryNo,
				OrderId:       util.Int64ToString(ret.DeliveryWithDetails.SalesDelivery.OrderId),
				OrderNo:       ret.DeliveryWithDetails.SalesDelivery.OrderNo,
				WarehouseId:   util.Int64ToString(ret.DeliveryWithDetails.SalesDelivery.WarehouseId),
				WarehouseNo:   warehouseNo,
				WarehouseName: warehouseName,
				DeliveryDate:  ret.DeliveryWithDetails.SalesDelivery.DeliveryDate,
				TotalQuantity: ret.DeliveryWithDetails.SalesDelivery.TotalQuantity,
				TotalAmount:   ret.DeliveryWithDetails.SalesDelivery.TotalAmount,
				Status:        ret.DeliveryWithDetails.SalesDelivery.Status,
				CreatedBy:     util.Int64ToString(ret.DeliveryWithDetails.SalesDelivery.CreatedBy),
				CreatedByNo:   createdByNo,
				CreatedByName: createdByName,
				CreatedAt:     ret.DeliveryWithDetails.SalesDelivery.CreatedAt,
			},
			Total: ret.DeliveryWithDetails.Total,
			SalesDeliveryDetail: func() []types.SalesDeliveryDetail {
				list := make([]types.SalesDeliveryDetail, 0, len(ret.DeliveryWithDetails.SalesDeliveryDetail))
				productBatchMap := make(map[int64]*productbatch.ProductBatch)
				productMap := make(map[int64]*product.Product)
				for _, d := range ret.DeliveryWithDetails.SalesDeliveryDetail {

					if _, ok := productMap[d.ProductId]; !ok {
						prod, err := l.svcCtx.ProductRPC.GetProductById(l.ctx, &product.GetProductByIdReq{
							Id: d.ProductId,
						})
						if err != nil {
							return nil
						}
						productMap[d.ProductId] = prod.Product
					}

					detail := types.SalesDeliveryDetail{
						Id:             util.Int64ToString(d.Id),
						DeliveryId:     util.Int64ToString(d.DeliveryId),
						ProductId:      util.Int64ToString(d.ProductId),
						ProductNo:      productMap[d.ProductId].ProductNo,
						ProductName:    d.ProductName,
						Unit:           d.Unit,
						Quantity:       d.Quantity,
						UnitPrice:      d.UnitPrice,
						Amount:         d.Amount,
						BatchId:        "",
						BatchNo:        "",
						ProductionDate: 0,
					}
					if d.BatchId > 0 {
						if _, ok := productBatchMap[d.BatchId]; !ok {
							productBatch, err := l.svcCtx.ProductRPC.GetProductBatchById(l.ctx, &productbatch.GetProductBatchByIdReq{
								Id: d.BatchId,
							})
							if err != nil {
								return nil
							}
							productBatchMap[d.BatchId] = productBatch.ProductBatch
						}
						detail.BatchId = util.Int64ToString(d.BatchId)
						detail.BatchNo = productBatchMap[d.BatchId].BatchNo
						detail.ProductionDate = productBatchMap[d.BatchId].ProductionDate

					}
					list = append(list, detail)
				}
				return list
			}(),
		},
	}
	return
}
