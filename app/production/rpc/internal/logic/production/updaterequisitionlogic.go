package productionlogic

import (
	"context"
	"database/sql"

	"erp/app/production/rpc/internal/svc"
	"erp/app/production/rpc/production"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateRequisitionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateRequisitionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateRequisitionLogic {
	return &UpdateRequisitionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateRequisitionLogic) UpdateRequisition(in *production.UpdateRequisitionReq) (*production.UpdateRequisitionResp, error) {
	req, err := l.svcCtx.MaterialRequisitionModel.FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	// 更新主表字段
	if in.Status > 0 {
		req.Status = in.Status
	}
	req.UpdatedBy = sql.NullInt64{Int64: in.UpdatedBy, Valid: in.UpdatedBy > 0}

	// 更新领料单主表（items 现在通过独立接口管理）
	err = l.svcCtx.MaterialRequisitionModel.Update(l.ctx, req)
	if err != nil {
		return nil, err
	}

	return &production.UpdateRequisitionResp{}, nil
}
