package workorder

import (
	"context"
	"erp/app/inventory/rpc/client/inventory"
	"erp/app/product/rpc/client/product"
	"erp/app/production/rpc/production"
	"erp/common/util"

	"erp/app/production/api/internal/svc"
	"erp/app/production/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetWorkOrderListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取生产工单列表
func NewGetWorkOrderListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetWorkOrderListLogic {
	return &GetWorkOrderListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetWorkOrderListLogic) GetWorkOrderList(req *types.WorkOrderListReq) (resp *types.WorkOrderListResp, err error) {
	productId, err := util.StringToInt64(req.ProductId)
	if err != nil {
		return nil, err
	}
	listResp, err := l.svcCtx.ProductionRPC.GetWorkOrderList(l.ctx, &production.WorkOrderListReq{
		Page:      req.Page,
		PageSize:  req.PageSize,
		ProductId: productId,
		Status:    req.Status,
		Priority:  req.Priority,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
	})
	if err != nil {
		return nil, err
	}

	productMap := make(map[int64]*product.Product)
	warehouseMap := make(map[int64]*inventory.WarehouseDetail)

	list := make([]types.WorkOrderInfo, 0, len(listResp.List))
	for _, wo := range listResp.List {
		if _, ok := productMap[wo.ProductId]; !ok {
			ret4, err := l.svcCtx.ProductRPC.GetProductById(l.ctx, &product.GetProductByIdReq{
				Id: wo.ProductId,
			})
			if err != nil {
				return nil, err
			}
			productMap[wo.ProductId] = ret4.Product
		}
		if _, ok := warehouseMap[wo.WarehouseId]; !ok {
			whRet, err := l.svcCtx.InventoryRPC.GetWarehouseById(l.ctx, &inventory.GetWarehouseByIdReq{
				Id: wo.WarehouseId,
			})
			if err != nil {
				return nil, err
			}
			warehouseMap[wo.WarehouseId] = whRet.Warehouse
		}
		list = append(list, types.WorkOrderInfo{
			Id:              util.Int64ToString(wo.Id),
			OrderNo:         wo.OrderNo,
			ProductId:       util.Int64ToString(wo.ProductId),
			ProductNo:       productMap[wo.ProductId].ProductNo,
			ProductName:     wo.ProductName,
			BomId:           util.Int64ToString(wo.BomId),
			BomVersion:      wo.BomVersion,
			Quantity:        wo.Quantity,
			CompletedQty:    wo.CompletedQty,
			QualifiedQty:    wo.QualifiedQty,
			DefectiveQty:    wo.DefectiveQty,
			WarehouseId:     util.Int64ToString(wo.WarehouseId),
			WarehouseNo:     warehouseMap[wo.WarehouseId].No,
			WarehouseName:   warehouseMap[wo.WarehouseId].Name,
			Status:          wo.Status,
			Priority:        wo.Priority,
			PlanStartDate:   wo.PlanStartDate,
			PlanEndDate:     wo.PlanEndDate,
			ActualStartDate: wo.ActualStartDate,
			ActualEndDate:   wo.ActualEndDate,
			Workshop:        wo.Workshop,
			Remark:          wo.Remark,
			CreatedAt:       wo.CreatedAt,
			CreatedBy:       util.Int64ToString(wo.CreatedBy),
			UpdatedAt:       wo.CreatedAt,
			UpdatedBy:       util.Int64ToString(wo.UpdatedBy),
		})
	}

	resp = &types.WorkOrderListResp{
		Total: listResp.Total,
		List:  list,
	}
	return
}
