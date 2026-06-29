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

type AddInventoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddInventoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddInventoryLogic {
	return &AddInventoryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// -----------------------inventory-----------------------
func (l *AddInventoryLogic) AddInventory(in *pb.AddInventoryReq) (*pb.AddInventoryResp, error) {
	inventoryId := util.GenerateSnowflake()
	transactionId := util.GenerateSnowflake()
	err := l.svcCtx.InventoryModel.InsertTransactCtx(l.ctx, &model.Inventory{
		InventoryId:  inventoryId,
		ProductId:    in.ProductId,
		WarehouseId:  in.WarehouseId,
		CurrentStock: in.CurrentStock,
		SafetyStock:  in.SafetyStock,
		LockedStock:  in.LockedStock,
	}, &model.InventoryTransaction{
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
		return nil, code.AddInventoryFail
	}
	return &pb.AddInventoryResp{
		InventoryId:            inventoryId,
		InventoryTransactionId: transactionId,
	}, nil
}
