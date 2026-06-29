package inventorylogic

import (
	"context"
	types2 "erp/app/inventory/rpc/internal/types"

	"erp/app/inventory/rpc/internal/svc"
	"erp/app/inventory/rpc/pb"

	"erp/app/inventory/rpc/internal/code"
	"github.com/zeromicro/go-zero/core/logx"
)

type SearchWarehouseLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchWarehouseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchWarehouseLogic {
	return &SearchWarehouseLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SearchWarehouseLogic) SearchWarehouse(in *pb.SearchWarehouseReq) (*pb.SearchWarehouseResp, error) {
	search, total, err := l.svcCtx.WarehouseModel.Search(l.ctx, &types2.SearchWarehouseParams{
		SearchCom: types2.SearchCom{
			Page:  in.Page,
			Limit: in.Limit,
		},
		Name:     in.Name,
		Location: in.Location,
		IsActive: in.IsActive,
	})
	if err != nil {

		return nil, code.GetWarehouseFail

	}
	list := make([]*pb.WarehouseDetail, 0, total)
	for _, one := range search {
		// 获取仓库已使用容量
		usedCapacity, err := l.svcCtx.InventoryModel.GetUsedCapacity(l.ctx, one.Id)
		if err != nil {
			l.Logger.Errorf("获取仓库已使用容量失败: %v", err)
			usedCapacity = 0
		}

		list = append(list, &pb.WarehouseDetail{
			Id:           one.Id,
			No:           one.No.String,
			Name:         one.Name,
			Location:     one.Location.String,
			ManagerId:    one.ManagerId.Int64,
			Capacity:     one.Capacity.Float64,
			UsedCapacity: usedCapacity,
			IsActive:     one.IsActive,
			CreatedAt:    one.CreatedAt.Unix(),
			CreatedBy:    one.CreatedBy.Int64,
			UpdatedAt:    one.UpdatedAt.Unix(),
			UpdatedBy:    one.UpdatedBy.Int64,
		})
	}
	return &pb.SearchWarehouseResp{
		Total:     total,
		Warehouse: list,
	}, nil
}
