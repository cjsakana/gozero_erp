package inventorylogic

import (
	"context"
	"database/sql"
	"erp/app/inventory/rpc/internal/model"

	"erp/app/inventory/rpc/internal/svc"
	"erp/app/inventory/rpc/pb"

	"erp/app/inventory/rpc/internal/code"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateWarehouseLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateWarehouseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateWarehouseLogic {
	return &UpdateWarehouseLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateWarehouseLogic) UpdateWarehouse(in *pb.UpdateWarehouseReq) (*pb.UpdateWarehouseResp, error) {
	err := l.svcCtx.WarehouseModel.XUpdate(l.ctx, &model.Warehouse{
		Id:        in.Id,
		Name:      in.Name,
		Location:  sql.NullString{String: in.Location, Valid: in.Location != ""},
		ManagerId: sql.NullInt64{Int64: in.ManagerId, Valid: in.ManagerId > 0},
		Capacity:  sql.NullFloat64{Float64: in.Capacity, Valid: in.Capacity > 0},
		IsActive:  in.IsActive,
		UpdatedBy: sql.NullInt64{Int64: in.UpdatedBy, Valid: in.UpdatedBy > 0},
	})
	if err != nil {
		return nil, code.UpdateWarehouseFail
	}

	return &pb.UpdateWarehouseResp{}, nil
}
