package productionlogic

import (
	"context"
	"erp/app/production/rpc/internal/svc"
	"erp/app/production/rpc/internal/types"
	"erp/app/production/rpc/production"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetWorkOrderListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetWorkOrderListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetWorkOrderListLogic {
	return &GetWorkOrderListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetWorkOrderListLogic) GetWorkOrderList(in *production.WorkOrderListReq) (*production.WorkOrderListResp, error) {
	workOrderList, total, err := l.svcCtx.WorkOrderModel.GetWorkOrderList(l.ctx, &types.GetWorkOrderListParams{
		Page:      in.Page,
		PageSize:  in.PageSize,
		ProductId: in.ProductId,
		Status:    in.Status,
		Priority:  in.Priority,
		StartDate: func() time.Time {
			if in.StartDate == 0 {
				return time.Time{}
			}
			return time.Unix(in.StartDate, 0)
		}(),
		EndDate: func() time.Time {
			if in.EndDate == 0 {
				return time.Time{}
			}
			return time.Unix(in.EndDate, 0)
		}(),
	})
	if err != nil {
		return nil, err
	}
	var pbWorkOrderList []*production.WorkOrderInfo

	for _, wo := range workOrderList {
		pbWorkOrderList = append(pbWorkOrderList, &production.WorkOrderInfo{
			Id:           wo.Id,
			OrderNo:      wo.OrderNo,
			ProductId:    wo.ProductId,
			ProductName:  wo.ProductName.String,
			BomId:        wo.BomId.Int64,
			Quantity:     wo.Quantity,
			CompletedQty: wo.CompletedQty,
			QualifiedQty: wo.QualifiedQty,
			DefectiveQty: wo.DefectiveQty,
			WarehouseId:  wo.WarehouseId.Int64,
			Status:       wo.Status,
			Priority:     wo.Priority,
			PlanStartDate: func() int64 {
				if wo.PlanStartDate.Valid {
					return wo.PlanStartDate.Time.Unix()
				}
				return 0
			}(),
			PlanEndDate: func() int64 {
				if wo.PlanEndDate.Valid {
					return wo.PlanEndDate.Time.Unix()
				}
				return 0
			}(),
			ActualStartDate: func() int64 {
				if wo.ActualStartDate.Valid {
					return wo.ActualStartDate.Time.Unix()
				}
				return 0
			}(),
			ActualEndDate: func() int64 {
				if wo.ActualEndDate.Valid {
					return wo.ActualEndDate.Time.Unix()
				}
				return 0
			}(),
			Workshop:  wo.Workshop.String,
			Remark:    wo.Remark.String,
			CreatedAt: wo.CreatedAt.Unix(),
			CreatedBy: wo.CreatedBy.Int64,
			UpdatedAt: wo.UpdatedAt.Unix(),
			UpdatedBy: wo.UpdatedBy.Int64,
			BomVersion: func() string {
				bom, err := l.svcCtx.BomModel.FindOne(l.ctx, wo.BomId.Int64)
				if err != nil {
					return ""
				}
				return bom.Version
			}(),
		})
	}

	return &production.WorkOrderListResp{
		Total: total,
		List:  pbWorkOrderList,
	}, nil
}
