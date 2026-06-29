package productionlogic

import (
	"context"

	"erp/app/production/rpc/internal/svc"
	"erp/app/production/rpc/production"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteWorkOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteWorkOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteWorkOrderLogic {
	return &DeleteWorkOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteWorkOrderLogic) DeleteWorkOrder(in *production.IdReq) (*production.EmptyResp, error) {
	err := l.svcCtx.WorkOrderModel.Delete(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	return &production.EmptyResp{}, nil
}
