package productionlogic

import (
	"context"
	"database/sql"
	"time"

	"erp/app/production/rpc/internal/svc"
	"erp/app/production/rpc/production"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateWorkOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateWorkOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateWorkOrderLogic {
	return &UpdateWorkOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateWorkOrderLogic) UpdateWorkOrder(in *production.UpdateWorkOrderReq) (*production.UpdateWorkOrderResp, error) {
	wo, err := l.svcCtx.WorkOrderModel.FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	// 更新字段
	if in.Quantity > 0 {
		wo.Quantity = in.Quantity
	}
	if in.CompletedQty > 0 {
		wo.CompletedQty = in.CompletedQty
	}
	if in.QualifiedQty > 0 {
		wo.QualifiedQty = in.QualifiedQty
	}
	if in.DefectiveQty > 0 {
		wo.DefectiveQty = in.DefectiveQty
	}
	if in.Status > 0 {
		wo.Status = in.Status
	}
	if in.Priority > 0 {
		wo.Priority = in.Priority
	}
	if in.PlanStartDate != 0 {
		wo.PlanStartDate = sql.NullTime{Time: time.Unix(in.PlanStartDate, 0), Valid: true}
	}
	if in.PlanEndDate != 0 {
		wo.PlanEndDate = sql.NullTime{Time: time.Unix(in.PlanEndDate, 0), Valid: true}
	}
	if in.ActualStartDate != 0 {
		wo.ActualStartDate = sql.NullTime{Time: time.Unix(in.ActualStartDate, 0), Valid: true}
	}
	if in.ActualEndDate != 0 {
		wo.ActualEndDate = sql.NullTime{Time: time.Unix(in.ActualEndDate, 0), Valid: true}

	}
	if in.Workshop != "" {
		wo.Workshop = sql.NullString{String: in.Workshop, Valid: true}
	}
	if in.Remark != "" {
		wo.Remark = sql.NullString{String: in.Remark, Valid: true}
	}
	wo.UpdatedBy = sql.NullInt64{Int64: in.UpdatedBy, Valid: in.UpdatedBy > 0}

	err = l.svcCtx.WorkOrderModel.Update(l.ctx, wo)
	if err != nil {
		return nil, err
	}

	return &production.UpdateWorkOrderResp{}, nil
}
