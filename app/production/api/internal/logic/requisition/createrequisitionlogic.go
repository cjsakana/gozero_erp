package requisition

import (
	"context"
	"erp/app/production/rpc/production"
	"erp/common/util"
	"erp/common/xtypes"

	"erp/app/production/api/internal/svc"
	"erp/app/production/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateRequisitionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建领料单
func NewCreateRequisitionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateRequisitionLogic {
	return &CreateRequisitionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateRequisitionLogic) CreateRequisition(req *types.CreateRequisitionReq) (resp *types.CreateRequisitionResp, err error) {
	createdBy, err := util.GetInt64FromCtx(l.ctx, xtypes.EmployeeIdKey)
	if err != nil {
		return nil, err
	}
	workOrderId, err := util.StringToInt64(req.WorkOrderId)
	if err != nil {
		return nil, err
	}
	warehouseId, err := util.StringToInt64(req.WarehouseId)
	if err != nil {
		return nil, err
	}
	approvedById, err := util.StringToInt64(req.ApprovedById)
	if err != nil {
		return nil, err
	}

	// 调用 RPC 创建领料单主表（items 现在通过独立接口创建）
	ret, err := l.svcCtx.ProductionRPC.CreateRequisition(l.ctx, &production.CreateRequisitionReq{
		WorkOrderId:     workOrderId,
		WarehouseId:     warehouseId,
		RequisitionDate: req.RequisitionDate,
		ApprovedBy:      approvedById,
		CreatedBy:       createdBy,
	})
	if err != nil {
		return nil, err
	}

	resp = &types.CreateRequisitionResp{
		Id: util.Int64ToString(ret.Id),
	}
	return
}
