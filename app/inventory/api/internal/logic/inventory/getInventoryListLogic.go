package inventory

import (
	"context"
	"erp/app/inventory/api/internal/svc"
	"erp/app/inventory/api/internal/types"
	"erp/app/inventory/rpc/client/inventory"
	"erp/app/inventory/rpc/pb"
	"erp/app/product/rpc/client/product"
	pb2 "erp/app/product/rpc/pb"
	"erp/common/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetInventoryListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetInventoryListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetInventoryListLogic {
	return &GetInventoryListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetInventoryListLogic) GetInventoryList(req *types.GetInventoryListReq) (resp *types.GetInventoryListResp, err error) {
	// 返回对象
	resp = &types.GetInventoryListResp{}

	warehouseId, err := util.StringToInt64(req.WarehouseId)
	if err != nil {
		return nil, err
	}

	warehouseMap := make(map[int64]*inventory.WarehouseDetail) //
	productMap := make(map[int64]*product.Product)             // productId -> product

	// 获取按 productName 搜索到的产品（如果有提供）
	if req.ProductName != "" {
		ret2, err := l.svcCtx.ProductRPC.SearchProduct(l.ctx, &pb2.SearchProductReq{
			Limit:       -1,
			ProductName: req.ProductName,
		})
		if err != nil {
			return nil, err
		}
		for _, p := range ret2.Product {
			productMap[p.Id] = p
		}
	}

	// 1. 只有 productName
	if req.ProductName != "" {
		for productId, p := range productMap {
			invResp, err := l.svcCtx.InventoryRPC.SearchInventory(l.ctx, &pb.SearchInventoryReq{
				Limit:       -1,
				ProductId:   productId,
				WarehouseId: warehouseId,
			})
			if err != nil {
				return nil, err
			}

			for _, inv := range invResp.Inventory {
				// 低库存过滤：只保留 CurrentStock < SafetyStock
				if req.LowStock && !(inv.CurrentStock < inv.SafetyStock) {
					continue
				}

				if _, ok := warehouseMap[inv.WarehouseId]; !ok {
					w, e := l.svcCtx.InventoryRPC.GetWarehouseById(l.ctx, &inventory.GetWarehouseByIdReq{
						Id: inv.WarehouseId,
					})
					if e != nil {
						return nil, nil
					}
					warehouseMap[inv.WarehouseId] = w.Warehouse
				}

				resp.Inventories = append(resp.Inventories, types.Inventory{
					InventoryId:   util.Int64ToString(inv.InventoryId),
					ProductId:     util.Int64ToString(productId),
					ProductNo:     p.ProductNo,
					ProductName:   p.ProductName,
					WarehouseId:   util.Int64ToString(inv.WarehouseId),
					WarehouseNo:   warehouseMap[inv.WarehouseId].No,
					WarehouseName: warehouseMap[inv.WarehouseId].Name,
					CurrentStock:  inv.CurrentStock,
					SafetyStock:   inv.SafetyStock,
					LockedStock:   inv.LockedStock,
				})
			}
		}
	} else {
		invResp, err := l.svcCtx.InventoryRPC.SearchInventory(l.ctx, &pb.SearchInventoryReq{
			Limit:       -1,
			WarehouseId: warehouseId,
		})
		if err != nil {
			return nil, err
		}

		for _, inv := range invResp.Inventory {
			// 低库存过滤：只保留 CurrentStock < SafetyStock
			if req.LowStock && !(inv.CurrentStock < inv.SafetyStock) {
				continue
			}

			if _, ok := warehouseMap[inv.WarehouseId]; !ok {
				w, e := l.svcCtx.InventoryRPC.GetWarehouseById(l.ctx, &inventory.GetWarehouseByIdReq{
					Id: inv.WarehouseId,
				})
				if e != nil {
					return nil, nil
				}
				warehouseMap[inv.WarehouseId] = w.Warehouse
			}

			if _, ok := productMap[inv.ProductId]; !ok {
				ret2, err := l.svcCtx.ProductRPC.GetProductById(l.ctx, &pb2.GetProductByIdReq{
					Id: inv.ProductId,
				})
				if err != nil {
					return nil, nil
				}
				productMap[inv.ProductId] = ret2.Product
			}

			resp.Inventories = append(resp.Inventories, types.Inventory{
				InventoryId:   util.Int64ToString(inv.InventoryId),
				ProductId:     util.Int64ToString(inv.ProductId),
				ProductNo:     productMap[inv.ProductId].ProductNo,
				ProductName:   productMap[inv.ProductId].ProductName,
				WarehouseId:   util.Int64ToString(inv.WarehouseId),
				WarehouseNo:   warehouseMap[inv.WarehouseId].No,
				WarehouseName: warehouseMap[inv.WarehouseId].Name,
				CurrentStock:  inv.CurrentStock,
				SafetyStock:   inv.SafetyStock,
				LockedStock:   inv.LockedStock,
			})
		}
	}

	// 使用过滤后的结果数量作为 total
	resp.Total = int64(len(resp.Inventories))
	return resp, nil
}
