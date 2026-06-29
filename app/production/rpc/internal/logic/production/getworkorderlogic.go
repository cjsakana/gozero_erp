package productionlogic

import (
	"context"

	"erp/app/production/rpc/internal/svc"
	"erp/app/production/rpc/production"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetWorkOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetWorkOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetWorkOrderLogic {
	return &GetWorkOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetWorkOrderLogic) GetWorkOrder(in *production.IdReq) (*production.WorkOrderInfo, error) {
	wo, err := l.svcCtx.WorkOrderModel.FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	bom, err := l.svcCtx.BomModel.FindOne(l.ctx, wo.BomId.Int64)
	if err != nil {
		return nil, err
	}

	return &production.WorkOrderInfo{
		Id:              wo.Id,
		OrderNo:         wo.OrderNo,
		ProductId:       wo.ProductId,
		ProductName:     wo.ProductName.String,
		BomId:           wo.BomId.Int64,
		Quantity:        wo.Quantity,
		CompletedQty:    wo.CompletedQty,
		QualifiedQty:    wo.QualifiedQty,
		DefectiveQty:    wo.DefectiveQty,
		WarehouseId:     wo.WarehouseId.Int64,
		Status:          wo.Status,
		Priority:        wo.Priority,
		PlanStartDate:   wo.PlanStartDate.Time.Unix(),
		PlanEndDate:     wo.PlanEndDate.Time.Unix(),
		ActualStartDate: wo.ActualStartDate.Time.Unix(),
		ActualEndDate:   wo.ActualEndDate.Time.Unix(),
		Workshop:        wo.Workshop.String,
		Remark:          wo.Remark.String,
		CreatedAt:       wo.CreatedAt.Unix(),
		CreatedBy:       wo.CreatedBy.Int64,
		UpdatedAt:       wo.UpdatedAt.Unix(),
		UpdatedBy:       wo.UpdatedBy.Int64,
		BomVersion:      bom.Version,
	}, nil
}
