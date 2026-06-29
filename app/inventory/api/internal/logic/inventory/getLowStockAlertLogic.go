package inventory

import (
	"context"
	"erp/app/inventory/rpc/pb"
	pb2 "erp/app/product/rpc/pb"
	"erp/common/util"

	"erp/app/inventory/api/internal/svc"
	"erp/app/inventory/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLowStockAlertLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetLowStockAlertLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLowStockAlertLogic {
	return &GetLowStockAlertLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetLowStockAlertLogic) GetLowStockAlert(req *types.GetLowStockAlertReq) (resp *types.GetLowStockAlertResp, err error) {
	ret, err := l.svcCtx.InventoryRPC.LowStockAlert(l.ctx, &pb.LowStockAlertReq{
		Page:  req.Page,
		Limit: req.Limit,
	})
	if err != nil {
		return nil, err
	}

	resp = &types.GetLowStockAlertResp{
		Total: ret.Total,
	}

	list := make([]types.LowStockAlert, 0)
	for _, inventory := range ret.Inventory {
		ret3, err := l.svcCtx.ProductRPC.GetProductById(l.ctx, &pb2.GetProductByIdReq{
			Id: inventory.ProductId,
		})
		if err != nil {
			return nil, err
		}
		ret4, err := l.svcCtx.InventoryRPC.GetWarehouseById(l.ctx, &pb.GetWarehouseByIdReq{
			Id: inventory.WarehouseId,
		})
		if err != nil {
			return nil, err
		}
		list = append(list, types.LowStockAlert{
			InventoryId:   util.Int64ToString(inventory.InventoryId),
			ProductId:     util.Int64ToString(inventory.ProductId),
			ProductNo:     ret3.Product.ProductNo,
			ProductName:   ret3.Product.ProductName,
			WarehouseId:   util.Int64ToString(inventory.WarehouseId),
			WarehouseNo:   ret4.Warehouse.No,
			WarehouseName: ret4.Warehouse.Name,
			CurrentStock:  inventory.CurrentStock,
			SafetyStock:   inventory.SafetyStock,
			Difference:    inventory.CurrentStock - inventory.SafetyStock,
		})
	}

	resp.Alerts = list
	return
}
