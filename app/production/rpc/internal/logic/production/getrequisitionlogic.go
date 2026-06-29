package productionlogic

import (
	"context"

	"erp/app/production/rpc/internal/svc"
	"erp/app/production/rpc/production"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRequisitionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetRequisitionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRequisitionLogic {
	return &GetRequisitionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetRequisitionLogic) GetRequisition(in *production.IdReq) (*production.RequisitionInfo, error) {
	req, err := l.svcCtx.MaterialRequisitionModel.FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	// 获取明细
	items, err := l.svcCtx.MaterialRequisitionItemModel.FindByRequisitionId(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	var itemInfos []*production.RequisitionItemInfo
	for _, item := range items {
		itemInfos = append(itemInfos, &production.RequisitionItemInfo{
			Id:             item.Id,
			RequisitionId:  item.RequisitionId,
			MaterialId:     item.MaterialId,
			MaterialName:   item.MaterialName.String,
			PlanQuantity:   item.PlanQuantity,
			ActualQuantity: item.ActualQuantity,
			Unit:           item.Unit.String,
			BatchNo:        item.BatchNo.String,
			Remark:         item.Remark.String,
		})
	}

	var approvedAt int64
	if req.ApprovedAt.Valid {
		approvedAt = req.ApprovedAt.Time.Unix()
	}

	return &production.RequisitionInfo{
		Id:              req.Id,
		RequisitionNo:   req.RequisitionNo,
		WorkOrderId:     req.WorkOrderId,
		WorkOrderNo:     req.WorkOrderNo.String,
		WarehouseId:     req.WarehouseId,
		RequisitionDate: req.RequisitionDate.Unix(),
		Status:          req.Status,
		ApprovedBy:      req.ApprovedBy.Int64,
		ApprovedAt:      approvedAt,
		CreatedAt:       req.CreatedAt.Unix(),
		CreatedBy:       req.CreatedBy.Int64,
		UpdatedAt:       req.UpdatedAt.Unix(),
		UpdatedBy:       req.UpdatedBy.Int64,
		Items:           itemInfos,
	}, nil
}
