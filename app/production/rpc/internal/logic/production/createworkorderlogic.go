package productionlogic

import (
	"context"
	"database/sql"
	"erp/app/production/rpc/internal/model"
	"erp/common/util"
	"time"

	"erp/app/production/rpc/internal/svc"
	"erp/app/production/rpc/production"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateWorkOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateWorkOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateWorkOrderLogic {
	return &CreateWorkOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 生产工单管理
func (l *CreateWorkOrderLogic) CreateWorkOrder(in *production.CreateWorkOrderReq) (*production.CreateWorkOrderResp, error) {
	// 生成工单雪花ID
	orderId := util.GenerateSnowflake()
	orderNo := util.GenerateNo("WO")

	workOrder := &model.WorkOrder{
		Id:              orderId,
		OrderNo:         orderNo,
		ProductId:       in.ProductId,
		ProductName:     sql.NullString{String: in.ProductName, Valid: true},
		BomId:           sql.NullInt64{Int64: in.BomId, Valid: in.BomId > 0},
		Quantity:        in.Quantity,
		CompletedQty:    0,
		QualifiedQty:    0,
		DefectiveQty:    0,
		WarehouseId:     sql.NullInt64{Int64: in.WarehouseId, Valid: in.WarehouseId > 0},
		Status:          1, // 未开工
		Priority:        in.Priority,
		PlanStartDate:   sql.NullTime{Time: time.Unix(in.PlanStartDate, 0), Valid: in.PlanStartDate != 0},
		PlanEndDate:     sql.NullTime{Time: time.Unix(in.PlanEndDate, 0), Valid: in.PlanEndDate != 0},
		ActualStartDate: sql.NullTime{},
		ActualEndDate:   sql.NullTime{},
		Workshop:        sql.NullString{String: in.Workshop, Valid: in.Workshop != ""},
		Remark:          sql.NullString{String: in.Remark, Valid: in.Remark != ""},
		CreatedBy:       sql.NullInt64{Int64: in.CreatedBy, Valid: in.CreatedBy > 0},
		UpdatedBy:       sql.NullInt64{},
	}

	_, err := l.svcCtx.WorkOrderModel.Insert(l.ctx, workOrder)
	if err != nil {
		return nil, err
	}

	return &production.CreateWorkOrderResp{Id: orderId}, nil
}
