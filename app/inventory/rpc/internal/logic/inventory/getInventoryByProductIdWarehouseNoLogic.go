package inventorylogic

import (
	"context"

	"erp/app/inventory/rpc/internal/code"
	"erp/app/inventory/rpc/internal/svc"
	"erp/app/inventory/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
)

type GetInventoryByProductIdWarehouseNoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetInventoryByProductIdWarehouseNoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetInventoryByProductIdWarehouseNoLogic {
	return &GetInventoryByProductIdWarehouseNoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetInventoryByProductIdWarehouseNoLogic) GetInventoryByProductIdWarehouseNo(in *pb.GetInventoryByProductIdWarehouseNoReq) (*pb.GetInventoryByProductIdWarehouseNoResp, error) {
	one, err := l.svcCtx.InventoryModel.FindOneByProductIdWarehouseId(l.ctx, in.ProductId, in.WarehouseId)
	if err != nil {
		if err == sqlc.ErrNotFound {
			return nil, code.InventoryNotFound
		}
		return nil, code.GetInventoryFail
	}

	return &pb.GetInventoryByProductIdWarehouseNoResp{
		Inventory: &pb.Inventory{
			InventoryId:  one.InventoryId,
			ProductId:    one.ProductId,
			WarehouseId:  one.WarehouseId,
			CurrentStock: one.CurrentStock,
			SafetyStock:  one.SafetyStock,
			LockedStock:  one.LockedStock,
		},
	}, nil
}
