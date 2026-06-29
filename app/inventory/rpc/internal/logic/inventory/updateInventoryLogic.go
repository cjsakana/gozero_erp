package inventorylogic

import (
	"context"
	"database/sql"
	"erp/app/inventory/rpc/internal/code"
	"erp/app/inventory/rpc/internal/model"
	"erp/app/inventory/rpc/internal/svc"
	"erp/app/inventory/rpc/pb"
	"erp/common/util"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateInventoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateInventoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateInventoryLogic {
	return &UpdateInventoryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateInventoryLogic) UpdateInventory(in *pb.UpdateInventoryReq) (*pb.UpdateInventoryResp, error) {
	// 只有在需要更新安全库存或锁定库存时才调用XUpdate
	if in.SafetyStock != 0 || in.LockedStock != 0 {
		err := l.svcCtx.InventoryModel.XUpdate(l.ctx, &model.Inventory{
			InventoryId: in.InventoryId,
			SafetyStock: in.SafetyStock,
			LockedStock: in.LockedStock,
		})
		if err != nil {
			return nil, code.UpdateInventoryFail
		}
		return &pb.UpdateInventoryResp{}, nil
	}

	transactionId := util.GenerateSnowflake()

	// 只有在需要更新库存数量时才调用UpdateTransactCtx
	if in.CurrentStock != 0 || in.ProductId != 0 || in.WarehouseId != 0 {
		err := l.svcCtx.InventoryModel.UpdateTransactCtx(l.ctx, &model.Inventory{
			InventoryId:  in.InventoryId,
			ProductId:    in.ProductId,
			WarehouseId:  in.WarehouseId,
			CurrentStock: in.CurrentStock,
		}, in.AdjustType, &model.InventoryTransaction{
			Id:              transactionId,
			ProductId:       in.ProductId,
			WarehouseId:     in.WarehouseId,
			BatchId:         sql.NullInt64{Int64: in.BatchId, Valid: in.BatchId > 0},
			TransactionType: in.TransactionType,
			Quantity:        in.CurrentStock,
			ReferenceType:   in.ReferenceType,
			ReferenceId:     sql.NullString{String: in.ReferenceId, Valid: in.ReferenceId != ""},
			OperatorId:      in.OperatorId,
		})
		if err != nil {
			return nil, code.UpdateInventoryFail
		}
	}

	return &pb.UpdateInventoryResp{}, nil
}
