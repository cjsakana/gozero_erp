package position

import (
	"context"
	"erp/app/hr/api/internal/svc"
	"erp/app/hr/api/internal/types"
	"erp/app/hr/rpc/pb"
	"erp/common/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPositionDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPositionDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPositionDetailLogic {
	return &GetPositionDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPositionDetailLogic) GetPositionDetail(req *types.GetPositionDetailRequest) (resp *types.GetPositionDetailResponse, err error) {
	id, err := util.StringToInt64(req.Id)
	if err != nil {
		return nil, err
	}
	
	ret, err := l.svcCtx.HrRPC.PositionZrpcClient.GetPositionById(l.ctx, &pb.GetPositionByIdReq{Id: id})
	if err != nil {
		return nil, err
	}
	resp = &types.GetPositionDetailResponse{
		Position: types.Position{
			Id:          util.Int64ToString(ret.Position.Id),
			Name:        ret.Position.Name,
			Description: ret.Position.Description,
		},
	}

	return
}
