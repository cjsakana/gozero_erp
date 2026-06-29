package inventorylogic

import (
	"context"

	"erp/app/inventory/rpc/internal/svc"
	"erp/app/inventory/rpc/pb"

	"erp/app/inventory/rpc/internal/code"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
)

type GetInventoryTransactionByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetInventoryTransactionByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetInventoryTransactionByIdLogic {
	return &GetInventoryTransactionByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetInventoryTransactionByIdLogic) GetInventoryTransactionById(in *pb.GetInventoryTransactionByIdReq) (*pb.GetInventoryTransactionByIdResp, error) {
	one, err := l.svcCtx.InventoryTransactionModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if err == sqlc.ErrNotFound {
			return nil, code.InventoryTransactionNotFound
		}
		return nil, code.InventoryTransactionNotFound
	}

	return &pb.GetInventoryTransactionByIdResp{
		InventoryTransaction: &pb.InventoryTransaction{
			Id:              one.Id,
			ProductId:       one.ProductId,
			WarehouseId:     one.WarehouseId,
			BatchId:         one.BatchId.Int64,
			TransactionType: one.TransactionType,
			Quantity:        one.Quantity,
			ReferenceType:   one.ReferenceType,
			ReferenceId:     one.ReferenceId.String,
			OperatorId:      one.OperatorId,
			CreatedAt:       one.CreatedAt.Unix(),
		},
	}, nil
}
