package positionlogic

import (
	"context"
	"erp/app/hr/rpc/internal/svc"
	types2 "erp/app/hr/rpc/internal/types"
	"erp/app/hr/rpc/pb"

	"erp/app/hr/rpc/internal/code"
	"github.com/zeromicro/go-zero/core/logx"
)

type SearchPositionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchPositionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchPositionLogic {
	return &SearchPositionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SearchPositionLogic) SearchPosition(in *pb.SearchPositionReq) (*pb.SearchPositionResp, error) {
	positions, total, err := l.svcCtx.PositionModel.Search(l.ctx, &types2.SearchPositionParams{
		SearchCom: types2.SearchCom{
			Page:  in.Page,
			Limit: in.Limit,
		},
		Name:        in.Name,
		Description: in.Description,
	})
	if err != nil {

		return nil, code.GetPositionFail

	}

	var pbPositions []*pb.Position
	for _, position := range positions {
		pbPositions = append(pbPositions, &pb.Position{
			Id:          position.Id,
			Name:        position.Name,
			Description: position.Description.String,
		})
	}

	return &pb.SearchPositionResp{
		Total:    total,
		Position: pbPositions,
	}, nil
}
