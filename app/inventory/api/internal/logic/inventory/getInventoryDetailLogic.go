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

type GetInventoryDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetInventoryDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetInventoryDetailLogic {
	return &GetInventoryDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetInventoryDetailLogic) GetInventoryDetail(req *types.GetInventoryDetailReq) (resp *types.GetInventoryDetailResp, err error) {
	id, err := util.StringToInt64(req.Id)
	if err != nil {
		return nil, err
	}
	ret, err := l.svcCtx.InventoryRPC.GetInventoryById(l.ctx, &pb.GetInventoryByIdReq{Id: id})
	if err != nil {
		return nil, err
	}
	ret2, err := l.svcCtx.ProductRPC.GetProductById(l.ctx, &pb2.GetProductByIdReq{
		Id: ret.Inventory.ProductId,
	})
	if err != nil {
		return nil, err
	}
	ret3, err := l.svcCtx.InventoryRPC.GetWarehouseById(l.ctx, &pb.GetWarehouseByIdReq{
		Id: ret.Inventory.WarehouseId,
	})
	if err != nil {
		return nil, err
	}
	resp = &types.GetInventoryDetailResp{
		Inventory: types.Inventory{
			InventoryId:   util.Int64ToString(ret.Inventory.InventoryId),
			ProductId:     util.Int64ToString(ret.Inventory.ProductId),
			ProductNo:     ret2.Product.ProductNo,
			ProductName:   ret2.Product.ProductName,
			WarehouseId:   util.Int64ToString(ret.Inventory.WarehouseId),
			WarehouseNo:   ret3.Warehouse.No,
			WarehouseName: ret3.Warehouse.Name,
			CurrentStock:  ret.Inventory.CurrentStock,
			SafetyStock:   ret.Inventory.SafetyStock,
			LockedStock:   ret.Inventory.LockedStock,
		},
	}

	return
}
