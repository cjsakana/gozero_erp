package inventorylogic

import (
	"context"
	"erp/app/inventory/rpc/internal/svc"
	types2 "erp/app/inventory/rpc/internal/types"
	"erp/app/inventory/rpc/pb"
	"time"

	"erp/app/inventory/rpc/internal/code"
	"github.com/zeromicro/go-zero/core/logx"
)

type SearchInventoryTransactionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchInventoryTransactionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchInventoryTransactionLogic {
	return &SearchInventoryTransactionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SearchInventoryTransactionLogic) SearchInventoryTransaction(in *pb.SearchInventoryTransactionReq) (*pb.SearchInventoryTransactionResp, error) {

	search, total, err := l.svcCtx.InventoryTransactionModel.Search(l.ctx, &types2.SearchInventoryTransactionParams{
		SearchCom: types2.SearchCom{
			Page:  in.Page,
			Limit: in.Limit,
		},
		ProductId:       in.ProductId,
		WarehouseId:     in.WarehouseId,
		BatchId:         in.BatchId,
		TransactionType: in.TransactionType,
		ReferenceType:   in.ReferenceType,
		StartTime:       time.Unix(in.StartTime, 0),
		EndTime:         time.Unix(in.EndTime, 0),
	})
	if err != nil {
		return nil, code.GetTransactionFail
	}
	list := make([]*pb.InventoryTransaction, len(search))
	for i, one := range search {
		list[i] = &pb.InventoryTransaction{
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
		}
	}

	return &pb.SearchInventoryTransactionResp{
		Total:                total,
		InventoryTransaction: list,
	}, nil
}
