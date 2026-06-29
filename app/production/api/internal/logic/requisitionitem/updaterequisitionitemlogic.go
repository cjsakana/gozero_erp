package requisitionitem

import (
	"context"

	"erp/app/production/api/internal/svc"
	"erp/app/production/api/internal/types"
	"erp/app/production/rpc/production"
	"erp/common/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateRequisitionItemLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新领料单明细
func NewUpdateRequisitionItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateRequisitionItemLogic {
	return &UpdateRequisitionItemLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateRequisitionItemLogic) UpdateRequisitionItem(req *types.UpdateRequisitionItemReq) (resp *types.UpdateRequisitionItemResp, err error) {
	// 转换 ID
	id, err := util.StringToInt64(req.Id)
	if err != nil {
		return nil, err
	}

	materialId, err := util.StringToInt64(req.MaterialId)
	if err != nil {
		return nil, err
	}

	// 调用 RPC 更新领料单明细
	_, err = l.svcCtx.ProductionRPC.UpdateRequisitionItem(l.ctx, &production.UpdateRequisitionItemReq{
		Id:             id,
		MaterialId:     materialId,
		PlanQuantity:   req.PlanQuantity,
		ActualQuantity: req.ActualQuantity,
		Unit:           req.Unit,
		BatchNo:        req.BatchNo,
		Remark:         req.Remark,
	})
	if err != nil {
		return nil, err
	}

	resp = &types.UpdateRequisitionItemResp{}
	return
}
