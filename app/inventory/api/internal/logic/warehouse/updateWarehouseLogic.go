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

type UpdateWarehouseLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateWarehouseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateWarehouseLogic {
	return &UpdateWarehouseLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateWarehouseLogic) UpdateWarehouse(req *types.UpdateWarehouseReq) (resp *types.EmptyResponse, err error) {
	employeeId, err := util.GetInt64FromCtx(l.ctx, xtypes.EmployeeIdKey)
	if err != nil {
		return nil, err
	}

	id, err := util.StringToInt64(req.Id)
	if err != nil {
		return nil, err
	}
	managerId, err := util.StringToInt64(req.ManagerId)
	if err != nil {
		return nil, err
	}
	_, err = l.svcCtx.InventoryRPC.UpdateWarehouse(l.ctx, &inventory.UpdateWarehouseReq{
		Id:        id,
		Name:      req.Name,
		Location:  req.Location,
		ManagerId: managerId,
		Capacity:  req.Capacity,
		IsActive:  req.IsActive,
		UpdatedBy: employeeId,
	})
	if err != nil {
		return nil, err
	}
	return
}
