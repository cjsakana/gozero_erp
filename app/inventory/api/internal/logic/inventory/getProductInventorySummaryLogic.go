package inventory

import (
	"context"
	"erp/app/inventory/api/internal/code"
	"erp/app/inventory/api/internal/svc"
	"erp/app/inventory/api/internal/types"
	"erp/app/inventory/rpc/client/inventory"
	"erp/app/inventory/rpc/pb"
	pb2 "erp/app/product/rpc/pb"
	"erp/common/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetProductInventorySummaryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetProductInventorySummaryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProductInventorySummaryLogic {
	return &GetProductInventorySummaryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetProductInventorySummaryLogic) GetProductInventorySummary(req *types.GetProductInventorySummaryReq) (resp *types.GetProductInventorySummaryResp, err error) {
	productId, err := util.StringToInt64(req.ProductId)
	if err != nil {
		return nil, err
	}

	// 如果没有商品ID，返回错误
	if productId == 0 {
		return nil, code.ParamsInvalid
	}

	ret, err := l.svcCtx.InventoryRPC.SearchInventory(l.ctx, &inventory.SearchInventoryReq{
		Limit:     -1,
		ProductId: productId,
	})
	if err != nil {
		return nil, err
	}
	ret2, err := l.svcCtx.ProductRPC.GetProductById(l.ctx, &pb2.GetProductByIdReq{Id: productId})
	if err != nil {
		return nil, err
	}

	totalStock := 0.0
	totalLocked := 0.0
	list := make([]*types.Inventory, 0)
	for _, v := range ret.Inventory {
		ret4, err := l.svcCtx.InventoryRPC.GetWarehouseById(l.ctx, &pb.GetWarehouseByIdReq{
			Id: v.WarehouseId,
		})
		if err != nil {
			return nil, err
		}
		totalStock += v.CurrentStock
		totalLocked += v.LockedStock
		list = append(list, &types.Inventory{
			InventoryId:   util.Int64ToString(v.InventoryId),
			ProductId:     util.Int64ToString(v.ProductId),
			ProductNo:     ret2.Product.ProductNo,
			ProductName:   ret2.Product.ProductName,
			WarehouseId:   util.Int64ToString(v.WarehouseId),
			WarehouseNo:   ret4.Warehouse.No,
			WarehouseName: ret4.Warehouse.Name,
			CurrentStock:  v.CurrentStock,
			SafetyStock:   v.SafetyStock,
			LockedStock:   v.LockedStock,
		})
	}
	resp = &types.GetProductInventorySummaryResp{
		ProductId:        req.ProductId,
		ProductName:      ret2.Product.ProductName,
		TotalStock:       totalStock,
		TotalLocked:      totalLocked,
		TotalAvailable:   totalStock - totalLocked,
		InventoryDetails: list,
	}
	return
}
