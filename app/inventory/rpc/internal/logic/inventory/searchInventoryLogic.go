package inventorylogic

import (
	"context"
	"erp/app/inventory/rpc/internal/svc"
	types2 "erp/app/inventory/rpc/internal/types"
	"erp/app/inventory/rpc/pb"

	"erp/app/inventory/rpc/internal/code"
	"github.com/zeromicro/go-zero/core/logx"
)

type SearchInventoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchInventoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchInventoryLogic {
	return &SearchInventoryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SearchInventoryLogic) SearchInventory(in *pb.SearchInventoryReq) (*pb.SearchInventoryResp, error) {
	search, total, err := l.svcCtx.InventoryModel.Search(l.ctx, &types2.SearchInventoryParams{
		SearchCom: types2.SearchCom{
			Page:  in.Page,
			Limit: in.Limit,
		},
		ProductId:    in.ProductId,
		WarehouseId:  in.WarehouseId,
		CurrentStock: in.CurrentStock,
		SafetyStock:  in.SafetyStock,
		LockedStock:  in.LockedStock,
	})
	if err != nil {

		return nil, code.GetInventoryFail

	}
	var list []*pb.Inventory
	for _, one := range search {
		list = append(list, &pb.Inventory{
			InventoryId:  one.InventoryId,
			ProductId:    one.ProductId,
			WarehouseId:  one.WarehouseId,
			CurrentStock: one.CurrentStock,
			SafetyStock:  one.SafetyStock,
			LockedStock:  one.LockedStock,
		})
	}

	return &pb.SearchInventoryResp{
		Inventory: list,
		Total:     total,
	}, nil
}
