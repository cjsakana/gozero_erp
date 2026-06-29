package inventorylogic

import (
	"context"
	"erp/app/inventory/rpc/internal/code"

	"erp/app/inventory/rpc/internal/svc"
	"erp/app/inventory/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetWarehouseByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetWarehouseByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetWarehouseByIdLogic {
	return &GetWarehouseByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetWarehouseByIdLogic) GetWarehouseById(in *pb.GetWarehouseByIdReq) (*pb.GetWarehouseByIdResp, error) {
	// 根据ID查询仓库
	one, err := l.svcCtx.WarehouseModel.FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, code.WarehouseNotFound
	}

	// 获取仓库已使用容量
	usedCapacity, err := l.svcCtx.InventoryModel.GetUsedCapacity(l.ctx, one.Id)
	if err != nil {
		l.Logger.Errorf("获取仓库已使用容量失败: %v", err)
		usedCapacity = 0
	}

	// 构建响应
	return &pb.GetWarehouseByIdResp{
		Warehouse: &pb.WarehouseDetail{
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
		},
	}, nil
}
