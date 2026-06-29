package positionlogic

import (
	"context"

	"erp/app/hr/rpc/internal/svc"
	"erp/app/hr/rpc/pb"

	"erp/app/hr/rpc/internal/code"
	"github.com/zeromicro/go-zero/core/logx"
)

type DelPositionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDelPositionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelPositionLogic {
	return &DelPositionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DelPositionLogic) DelPosition(in *pb.DelPositionReq) (*pb.DelPositionResp, error) {
	err := l.svcCtx.PositionModel.Delete(l.ctx, in.Id)
	if err != nil {

		return nil, code.DeletePositionFail

	}

	return &pb.DelPositionResp{}, nil
}
