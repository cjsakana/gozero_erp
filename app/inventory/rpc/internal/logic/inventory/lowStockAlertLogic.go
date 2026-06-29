package inventorylogic

import (
	"context"

	"erp/app/inventory/rpc/internal/code"
	"erp/app/inventory/rpc/internal/svc"
	"erp/app/inventory/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type LowStockAlertLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLowStockAlertLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LowStockAlertLogic {
	return &LowStockAlertLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LowStockAlertLogic) LowStockAlert(in *pb.LowStockAlertReq) (*pb.LowStockAlertResp, error) {
	inventories, err := l.svcCtx.InventoryModel.LowStockAlert(l.ctx)
	if err != nil {
		return nil, code.GetInventoryFail
	}
	list := []*pb.Inventory{}
	for _, item := range inventories {
		list = append(list, &pb.Inventory{
			InventoryId:  item.InventoryId,
			ProductId:    item.ProductId,
			WarehouseId:  item.WarehouseId,
			CurrentStock: item.CurrentStock,
			SafetyStock:  item.SafetyStock,
			LockedStock:  item.LockedStock,
		})
	}

	return &pb.LowStockAlertResp{
		Total:     int64(len(list)),
		Inventory: list,
	}, nil
}
