package requisition

import (
	"context"
	"erp/app/production/rpc/production"
	"erp/common/util"
	"erp/common/xtypes"

	"erp/app/production/api/internal/svc"
	"erp/app/production/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApproveRequisitionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 审批领料单
func NewApproveRequisitionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApproveRequisitionLogic {
	return &ApproveRequisitionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApproveRequisitionLogic) ApproveRequisition(req *types.ApproveRequisitionReq) (resp *types.ApproveRequisitionResp, err error) {
	approverId, err := util.GetInt64FromCtx(l.ctx, xtypes.EmployeeIdKey)
	if err != nil {
		return nil, err
	}

	id, err := util.StringToInt64(req.Id)
	if err != nil {
		return nil, err
	}
	_, err = l.svcCtx.ProductionRPC.ApproveRequisition(l.ctx, &production.ApproveRequisitionReq{
		Id:         id,
		Status:     req.Status,
		ApproverId: approverId,
	})
	if err != nil {
		return nil, err
	}

	resp = &types.ApproveRequisitionResp{}
	return
}
