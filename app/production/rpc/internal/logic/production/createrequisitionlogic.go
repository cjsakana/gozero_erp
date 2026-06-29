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

type CreateRequisitionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateRequisitionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateRequisitionLogic {
	return &CreateRequisitionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 领料单管理
func (l *CreateRequisitionLogic) CreateRequisition(in *production.CreateRequisitionReq) (*production.CreateRequisitionResp, error) {
	// 生成领料单雪花ID
	requisitionId := util.GenerateSnowflake()
	requisitionNo := util.GenerateNo("MR")

	var requisitionDate time.Time
	if in.RequisitionDate != 0 {
		requisitionDate = time.Unix(in.RequisitionDate, 0)
	} else {
		requisitionDate = time.Now()
	}

	workOrder, err := l.svcCtx.WorkOrderModel.FindOne(l.ctx, in.WorkOrderId)
	if err != nil {
		return nil, err
	}

	// 创建领料单主表（items 现在通过独立接口创建）
	requisition := &model.MaterialRequisition{
		Id:              requisitionId,
		RequisitionNo:   requisitionNo,
		WorkOrderId:     in.WorkOrderId,
		WorkOrderNo:     sql.NullString{String: workOrder.OrderNo, Valid: true},
		WarehouseId:     in.WarehouseId,
		RequisitionDate: requisitionDate,
		Status:          1, // 待审批
		ApprovedBy:      sql.NullInt64{Int64: in.ApprovedBy, Valid: in.ApprovedBy > 0},
		CreatedBy:       sql.NullInt64{Int64: in.CreatedBy, Valid: in.CreatedBy > 0},
		UpdatedBy:       sql.NullInt64{},
	}

	_, err = l.svcCtx.MaterialRequisitionModel.Insert(l.ctx, requisition)
	if err != nil {
		return nil, err
	}

	return &production.CreateRequisitionResp{Id: requisitionId}, nil
}
