package inventory

import (
	"context"
	"erp/app/inventory/rpc/pb"
	"erp/common/util"

	"erp/app/inventory/api/internal/svc"
	"erp/app/inventory/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SetInventoryLockLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSetInventoryLockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetInventoryLockLogic {
	return &SetInventoryLockLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SetInventoryLockLogic) SetInventoryLock(req *types.SetInventoryLockReq) (resp *types.EmptyResponse, err error) {
	inventoryId, err := util.StringToInt64(req.InventoryId)
	if err != nil {
		return nil, err
	}
	_, err = l.svcCtx.InventoryRPC.UpdateInventory(l.ctx, &pb.UpdateInventoryReq{
		InventoryId: inventoryId,
		SafetyStock: req.SafetyStock,
		LockedStock: req.LockedStock,
	})
	if err != nil {
		return nil, err
	}

	return
}
