package position

import (
	"context"
	"erp/app/hr/api/internal/svc"
	"erp/app/hr/api/internal/types"
	"erp/app/hr/rpc/pb"
	"erp/common/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPositionListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPositionListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPositionListLogic {
	return &GetPositionListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPositionListLogic) GetPositionList(req *types.GetPositionListRequest) (resp *types.GetPositionListResponse, err error) {
	ret, err := l.svcCtx.HrRPC.PositionZrpcClient.SearchPosition(l.ctx, &pb.SearchPositionReq{
		Page:        req.Page,
		Limit:       req.Limit,
		Name:        req.Keyword,
		Description: req.Keyword,
	})
	if err != nil {
		return nil, err
	}

	resp = &types.GetPositionListResponse{
		Total: ret.Total,
	}
	for _, v := range ret.Position {

		resp.List = append(resp.List, &types.Position{
			Id:          util.Int64ToString(v.Id),
			Name:        v.Name,
			Description: v.Description,
		})
	}

	return
}
