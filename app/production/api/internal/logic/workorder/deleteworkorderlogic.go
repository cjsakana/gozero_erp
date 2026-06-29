package workorder

import (
	"context"
	"erp/app/production/rpc/production"
	"erp/common/util"

	"erp/app/production/api/internal/svc"
	"erp/app/production/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteWorkOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除生产工单
func NewDeleteWorkOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteWorkOrderLogic {
	return &DeleteWorkOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteWorkOrderLogic) DeleteWorkOrder(req *types.IdReq) (resp *types.DeleteWorkOrderResp, err error) {
	id, err := util.StringToInt64(req.Id)
	if err != nil {
		return nil, err
	}
	_, err = l.svcCtx.ProductionRPC.DeleteWorkOrder(l.ctx, &production.IdReq{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	resp = &types.DeleteWorkOrderResp{}
	return
}
