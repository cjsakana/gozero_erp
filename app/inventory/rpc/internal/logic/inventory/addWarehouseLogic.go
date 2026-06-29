package inventorylogic

import (
	"context"
	"database/sql"
	"erp/app/inventory/rpc/internal/model"
	"erp/app/inventory/rpc/internal/svc"
	"erp/app/inventory/rpc/pb"
	"erp/common/util"

	"erp/app/inventory/rpc/internal/code"
	"github.com/zeromicro/go-zero/core/logx"
)

type AddWarehouseLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddWarehouseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddWarehouseLogic {
	return &AddWarehouseLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// -----------------------warehouse-----------------------
func (l *AddWarehouseLogic) AddWarehouse(in *pb.AddWarehouseReq) (*pb.AddWarehouseResp, error) {
	// 生成仓库ID
	id := util.GenerateSnowflake()

	_, err := l.svcCtx.WarehouseModel.Insert(l.ctx, &model.Warehouse{
		Id:        id,
		No:        sql.NullString{String: in.No, Valid: in.No != ""},
		Name:      in.Name,
		Location:  sql.NullString{String: in.Location, Valid: in.Location != ""},
		ManagerId: sql.NullInt64{Int64: in.ManagerId, Valid: in.ManagerId > 0},
		Capacity:  sql.NullFloat64{Float64: in.Capacity, Valid: in.Capacity > 0},
		IsActive:  in.IsActive,
		CreatedBy: sql.NullInt64{Int64: in.CreatedBy, Valid: in.CreatedBy > 0},
		UpdatedBy: sql.NullInt64{Int64: in.CreatedBy, Valid: in.CreatedBy > 0},
	})
	if err != nil {
		return nil, code.AddWarehouseFail
	}

	return &pb.AddWarehouseResp{Id: id}, nil
}
