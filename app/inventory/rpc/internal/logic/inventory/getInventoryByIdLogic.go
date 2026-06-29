package inventorylogic

import (
	"context"

	"erp/app/inventory/rpc/internal/svc"
	"erp/app/inventory/rpc/pb"

	"erp/app/inventory/rpc/internal/code"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
)

type GetInventoryByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetInventoryByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetInventoryByIdLogic {
	return &GetInventoryByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetInventoryByIdLogic) GetInventoryById(in *pb.GetInventoryByIdReq) (*pb.GetInventoryByIdResp, error) {
	one, err := l.svcCtx.InventoryModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if err == sqlc.ErrNotFound {
			return nil, code.InventoryNotFound
		}
		return nil, code.InventoryNotFound
	}

	return &pb.GetInventoryByIdResp{
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
