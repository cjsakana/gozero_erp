package productionlogic

import (
	"context"
	"database/sql"
	"time"

	"erp/app/production/rpc/internal/svc"
	"erp/app/production/rpc/production"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApproveRequisitionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewApproveRequisitionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApproveRequisitionLogic {
	return &ApproveRequisitionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ApproveRequisitionLogic) ApproveRequisition(in *production.ApproveRequisitionReq) (*production.ApproveRequisitionResp, error) {
	req, err := l.svcCtx.MaterialRequisitionModel.FindOne(l.ctx, in.Id)
	if err != nil {
		l.Logger.Errorf("领料单不存在, %v \n", err)
		return nil, err
	}

	// 更新审批状态
	req.Status = in.Status
	
	// 设置审批人和审批时间
	req.ApprovedBy = sql.NullInt64{Int64: in.ApproverId, Valid: in.ApproverId > 0}
	req.ApprovedAt = sql.NullTime{Time: time.Now(), Valid: true}
	
	// 更新更新人
	req.UpdatedBy = sql.NullInt64{Int64: in.ApproverId, Valid: in.ApproverId > 0}

	err = l.svcCtx.MaterialRequisitionModel.Update(l.ctx, req)
	if err != nil {
		l.Logger.Errorf("审批失败, %v \n", err)
		return nil, err
	}

	return &production.ApproveRequisitionResp{}, nil
}
