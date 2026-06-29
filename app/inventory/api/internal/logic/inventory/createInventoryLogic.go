package inventory

import (
	"context"
	"erp/app/inventory/rpc/pb"
	"erp/common/util"
	"erp/common/xtypes"

	"erp/app/inventory/api/internal/svc"
	"erp/app/inventory/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateInventoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateInventoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateInventoryLogic {
	return &CreateInventoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateInventoryLogic) CreateInventory(req *types.CreateInventoryReq) (resp *types.CreateInventoryResp, err error) {
	employeeId, err := util.GetInt64FromCtx(l.ctx, xtypes.EmployeeIdKey)
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
	ret, err := l.svcCtx.InventoryRPC.AddInventory(l.ctx, &pb.AddInventoryReq{
		ProductId:       productId,
		WarehouseId:     warehouseId,
		CurrentStock:    req.CurrentStock,
		SafetyStock:     req.SafetyStock,
		LockedStock:     req.LockedStock,
		TransactionType: req.TransactionType,
		ReferenceType:   req.ReferenceType,
		ReferenceId:     req.ReferenceId,
		OperatorId:      employeeId,
		BatchId:         batchId,
	})
	if err != nil {
		return nil, err
	}
	resp = &types.CreateInventoryResp{
		InventoryId:            util.Int64ToString(ret.InventoryId),
		InventoryTransactionId: util.Int64ToString(ret.InventoryTransactionId),
	}

	return
}
