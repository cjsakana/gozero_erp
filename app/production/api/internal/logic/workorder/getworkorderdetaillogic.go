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

type GetWorkOrderDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取生产工单详情
func NewGetWorkOrderDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetWorkOrderDetailLogic {
	return &GetWorkOrderDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetWorkOrderDetailLogic) GetWorkOrderDetail(req *types.IdReq) (resp *types.WorkOrderInfo, err error) {
	id, err := util.StringToInt64(req.Id)
	if err != nil {
		return nil, err
	}
	woInfo, err := l.svcCtx.ProductionRPC.GetWorkOrder(l.ctx, &production.IdReq{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	productRet, err := l.svcCtx.ProductRPC.GetProductById(l.ctx, &product.GetProductByIdReq{
		Id: woInfo.ProductId,
	})
	if err != nil {
		return nil, err
	}

	whRet, err := l.svcCtx.InventoryRPC.GetWarehouseById(l.ctx, &inventory.GetWarehouseByIdReq{
		Id: woInfo.WarehouseId,
	})
	if err != nil {
		return nil, err
	}

	resp = &types.WorkOrderInfo{
		Id:              util.Int64ToString(woInfo.Id),
		OrderNo:         woInfo.OrderNo,
		ProductId:       util.Int64ToString(woInfo.ProductId),
		ProductNo:       productRet.Product.ProductNo,
		ProductName:     woInfo.ProductName,
		BomId:           util.Int64ToString(woInfo.BomId),
		BomVersion:      woInfo.BomVersion,
		Quantity:        woInfo.Quantity,
		CompletedQty:    woInfo.CompletedQty,
		QualifiedQty:    woInfo.QualifiedQty,
		DefectiveQty:    woInfo.DefectiveQty,
		WarehouseId:     util.Int64ToString(woInfo.WarehouseId),
		WarehouseNo:     whRet.Warehouse.No,
		WarehouseName:   whRet.Warehouse.Name,
		Status:          woInfo.Status,
		Priority:        woInfo.Priority,
		PlanStartDate:   woInfo.PlanStartDate,
		PlanEndDate:     woInfo.PlanEndDate,
		ActualStartDate: woInfo.ActualStartDate,
		ActualEndDate:   woInfo.ActualEndDate,
		Workshop:        woInfo.Workshop,
		Remark:          woInfo.Remark,
		CreatedAt:       woInfo.CreatedAt,
		CreatedBy:       util.Int64ToString(woInfo.CreatedBy),
		UpdatedAt:       woInfo.CreatedAt,
		UpdatedBy:       util.Int64ToString(woInfo.UpdatedBy),
	}
	return
}
