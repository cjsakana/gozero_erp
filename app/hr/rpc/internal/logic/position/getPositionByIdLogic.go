package positionlogic

import (
	"context"

	"erp/app/hr/rpc/internal/svc"
	"erp/app/hr/rpc/pb"

	"erp/app/hr/rpc/internal/code"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
)

type GetPositionByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetPositionByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPositionByIdLogic {
	return &GetPositionByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetPositionByIdLogic) GetPositionById(in *pb.GetPositionByIdReq) (*pb.GetPositionByIdResp, error) {
	one, err := l.svcCtx.PositionModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if err == sqlc.ErrNotFound {
			return nil, code.PositionNotFound
		}
		return nil, code.PositionNotFound
	}
	return &pb.GetPositionByIdResp{
		Position: &pb.Position{
			Id:          one.Id,
			Name:        one.Name,
			Description: one.Description.String,
		},
	}, nil
}
