package warehouse

import (
	"context"
	"erp/app/inventory/api/internal/svc"
	"erp/app/inventory/api/internal/types"
	"erp/app/inventory/rpc/client/inventory"
	"erp/common/util"
	"erp/common/xtypes"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateWarehouseLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateWarehouseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateWarehouseLogic {
	return &CreateWarehouseLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateWarehouseLogic) CreateWarehouse(req *types.CreateWarehouseReq) (resp *types.CreateWarehouseResp, err error) {
	employeeId, err := util.GetInt64FromCtx(l.ctx, xtypes.EmployeeIdKey)
	if err != nil {
		return nil, err
	}
	managerId, err := util.StringToInt64(req.ManagerId)
	if err != nil {
		return nil, err
	}

	no := util.GenerateNo("WH")
	addResp, err := l.svcCtx.InventoryRPC.AddWarehouse(l.ctx, &inventory.AddWarehouseReq{
		No:        no,
		Name:      req.Name,
		Location:  req.Location,
		ManagerId: managerId,
		Capacity:  req.Capacity,
		IsActive:  req.IsActive,
		CreatedBy: employeeId,
	})
	if err != nil {
		return nil, err
	}

	resp = &types.CreateWarehouseResp{
		Id: util.Int64ToString(addResp.Id),
	}
	return
}
