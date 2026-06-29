package image

import (
	"context"
	"erp/common/util"

	"erp/app/image/api/internal/svc"
	"erp/app/image/api/internal/types"
	"erp/app/image/rpc/image"

	"github.com/zeromicro/go-zero/core/logx"
)

type DelImageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDelImageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelImageLogic {
	return &DelImageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DelImageLogic) DelImage(req *types.DelImageReq) (resp *types.DelImageResp, err error) {

	id, err := util.StringToInt64(req.Id)
	if err != nil {
		return nil, err
	}

	_, err = l.svcCtx.ImageRPC.DelImage(l.ctx, &image.DelImageReq{Id: id})
	if err != nil {
		return nil, err
	}
	return &types.DelImageResp{}, nil
}
