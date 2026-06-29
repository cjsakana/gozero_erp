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

type UpdateRequisitionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新领料单
func NewUpdateRequisitionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateRequisitionLogic {
	return &UpdateRequisitionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateRequisitionLogic) UpdateRequisition(req *types.UpdateRequisitionReq) (resp *types.UpdateRequisitionResp, err error) {
	updatedBy, err := util.GetInt64FromCtx(l.ctx, xtypes.EmployeeIdKey)
	if err != nil {
		return nil, err
	}
	id, err := util.StringToInt64(req.Id)
	if err != nil {
		return nil, err
	}

	// 调用 RPC 更新领料单主表（items 现在通过独立接口管理）
	_, err = l.svcCtx.ProductionRPC.UpdateRequisition(l.ctx, &production.UpdateRequisitionReq{
		Id:        id,
		Status:    req.Status,
		UpdatedBy: updatedBy,
	})
	if err != nil {
		return nil, err
	}

	resp = &types.UpdateRequisitionResp{}
	return
}
