package productionlogic

import (
	"context"
	"erp/app/production/rpc/internal/svc"
	"erp/app/production/rpc/internal/types"
	"erp/app/production/rpc/production"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRequisitionListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetRequisitionListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRequisitionListLogic {
	return &GetRequisitionListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetRequisitionListLogic) GetRequisitionList(in *production.RequisitionListReq) (*production.RequisitionListResp, error) {
	requisitionList, total, err := l.svcCtx.MaterialRequisitionModel.GetRequisitionList(l.ctx, &types.GetRequisitionListParams{
		Page:        in.Page,
		PageSize:    in.PageSize,
		WorkOrderId: in.WorkOrderId,
		WarehouseId: in.WarehouseId,
		Status:      in.Status,
	})
	if err != nil {
		return nil, err
	}
	var pbRequisitionList []*production.RequisitionInfo

	for _, r := range requisitionList {
		var approvedAt int64
		if r.ApprovedAt.Valid {
			approvedAt = r.ApprovedAt.Time.Unix()
		}
		
		pbRequisitionList = append(pbRequisitionList, &production.RequisitionInfo{
			Id:              r.Id,
			RequisitionNo:   r.RequisitionNo,
			WorkOrderId:     r.WorkOrderId,
			WorkOrderNo:     r.WorkOrderNo.String,
			WarehouseId:     r.WarehouseId,
			RequisitionDate: r.RequisitionDate.Unix(),
			Status:          r.Status,
			ApprovedBy:      r.ApprovedBy.Int64,
			ApprovedAt:      approvedAt,
			CreatedAt:       r.CreatedAt.Unix(),
			CreatedBy:       r.CreatedBy.Int64,
			UpdatedAt:       r.UpdatedAt.Unix(),
			UpdatedBy:       r.UpdatedBy.Int64,
			Items:           nil,
		})
	}

	return &production.RequisitionListResp{
		Total: total,
		List:  pbRequisitionList,
	}, nil
}
