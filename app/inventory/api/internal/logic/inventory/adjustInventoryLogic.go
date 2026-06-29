package inventory

import (
	"context"
	"erp/app/inventory/api/internal/svc"
	"erp/app/inventory/api/internal/types"
	"erp/app/inventory/rpc/pb"
	"erp/common/util"
	"erp/common/xtypes"
	"github.com/zeromicro/go-zero/core/logx"
)

type AdjustInventoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdjustInventoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdjustInventoryLogic {
	return &AdjustInventoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdjustInventoryLogic) AdjustInventory(req *types.AdjustInventoryReq) (resp *types.AdjustInventoryResp, err error) {
	employeeId, err := util.GetInt64FromCtx(l.ctx, xtypes.EmployeeIdKey)
	if err != nil {
		return nil, err
	}
	inventoryId, err := util.StringToInt64(req.InventoryId)
	if err != nil {
		return nil, err
	}
	productId, err := util.StringToInt64(req.ProductId)
	if err != nil {
		return nil, err
	}
	warehouseId, err := util.StringToInt64(req.WarehouseId)
	if err != nil {
		return nil, err
	}
	batchId, err := util.StringToInt64(req.BatchId)
	if err != nil {
		return nil, err
	}

	_, err = l.svcCtx.InventoryRPC.UpdateInventory(l.ctx, &pb.UpdateInventoryReq{
		InventoryId:     inventoryId,
		ProductId:       productId,
		WarehouseId:     warehouseId,
		AdjustType:      req.AdjustType,
		CurrentStock:    req.Quantity,
		TransactionType: req.TransactionType,
		ReferenceType:   req.ReferenceType,
		ReferenceId:     req.ReferenceId,
		OperatorId:      employeeId,
		BatchId:         batchId,
	})
	if err != nil {
		return nil, err
	}
	return
}
