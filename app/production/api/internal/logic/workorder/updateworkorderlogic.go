package workorder

import (
	"context"
	"erp/app/production/rpc/production"
	"erp/common/util"
	"erp/common/xtypes"

	"erp/app/production/api/internal/svc"
	"erp/app/production/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateWorkOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新生产工单
func NewUpdateWorkOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateWorkOrderLogic {
	return &UpdateWorkOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateWorkOrderLogic) UpdateWorkOrder(req *types.UpdateWorkOrderReq) (resp *types.UpdateWorkOrderResp, err error) {
	updatedBy, err := util.GetInt64FromCtx(l.ctx, xtypes.EmployeeIdKey)
	if err != nil {
		return nil, err
	}
	id, err := util.StringToInt64(req.Id)
	if err != nil {
		return nil, err
	}

	_, err = l.svcCtx.ProductionRPC.UpdateWorkOrder(l.ctx, &production.UpdateWorkOrderReq{
		Id:              id,
		Quantity:        req.Quantity,
		CompletedQty:    req.CompletedQty,
		QualifiedQty:    req.QualifiedQty,
		DefectiveQty:    req.DefectiveQty,
		Status:          req.Status,
		Priority:        req.Priority,
		PlanStartDate:   req.PlanStartDate,
		PlanEndDate:     req.PlanEndDate,
		ActualStartDate: req.ActualStartDate,
		ActualEndDate:   req.ActualEndDate,
		Workshop:        req.Workshop,
		Remark:          req.Remark,
		UpdatedBy:       updatedBy,
	})
	if err != nil {
		return nil, err
	}

	resp = &types.UpdateWorkOrderResp{}
	return
}
