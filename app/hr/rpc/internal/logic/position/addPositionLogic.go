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

type AddPositionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddPositionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddPositionLogic {
	return &AddPositionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// -----------------------岗位表-----------------------
func (l *AddPositionLogic) AddPosition(in *pb.AddPositionReq) (*pb.AddPositionResp, error) {
	res, err := l.svcCtx.PositionModel.Insert(l.ctx, &model.Position{
		Name:        in.Name,
		Description: sql.NullString{String: in.Description, Valid: true},
	})
	if err != nil {

		return nil, code.AddPositionFail

	}
	id, _ := res.LastInsertId()
	return &pb.AddPositionResp{
		Id: id,
	}, nil
}
