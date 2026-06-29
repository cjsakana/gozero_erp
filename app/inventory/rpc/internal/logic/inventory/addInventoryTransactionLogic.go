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

type AddInventoryTransactionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddInventoryTransactionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddInventoryTransactionLogic {
	return &AddInventoryTransactionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// -----------------------inventoryTransaction-----------------------
func (l *AddInventoryTransactionLogic) AddInventoryTransaction(in *pb.AddInventoryTransactionReq) (*pb.AddInventoryTransactionResp, error) {
	transactionId := util.GenerateSnowflake()

	_, err := l.svcCtx.InventoryTransactionModel.Insert(l.ctx, &model.InventoryTransaction{
		Id:              transactionId,
		ProductId:       in.ProductId,
		WarehouseId:     in.WarehouseId,
		BatchId:         sql.NullInt64{Int64: in.BatchId, Valid: in.BatchId > 0},
		TransactionType: in.TransactionType,
		Quantity:        in.Quantity,
		ReferenceType:   in.ReferenceType,
		ReferenceId:     sql.NullString{String: in.ReferenceId, Valid: in.ReferenceId != ""},
		OperatorId:      in.OperatorId,
	})
	if err != nil {

		return nil, code.CreateTransactionFail

	}

	return &pb.AddInventoryTransactionResp{
		Id: transactionId,
	}, nil
}
