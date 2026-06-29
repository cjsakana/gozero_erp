package positionlogic

import (
	"context"
	"database/sql"
	"erp/app/hr/rpc/internal/model"

	"erp/app/hr/rpc/internal/svc"
	"erp/app/hr/rpc/pb"

	"erp/app/hr/rpc/internal/code"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdatePositionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdatePositionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePositionLogic {
	return &UpdatePositionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdatePositionLogic) UpdatePosition(in *pb.UpdatePositionReq) (*pb.UpdatePositionResp, error) {
	err := l.svcCtx.PositionModel.XUpdate(l.ctx, &model.Position{
		Id:          in.Id,
		Name:        in.Name,
		Description: sql.NullString{String: in.Description, Valid: true},
	})
	if err != nil {

		return nil, code.UpdatePositionFail

	}

	return &pb.UpdatePositionResp{}, nil
}
