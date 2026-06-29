package image

import (
	"context"
	"erp/common/util"

	"erp/app/image/api/internal/svc"
	"erp/app/image/api/internal/types"
	"erp/app/image/rpc/image"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateImageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateImageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateImageLogic {
	return &UpdateImageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateImageLogic) UpdateImage(req *types.UpdateImageReq) (resp *types.UpdateImageResp, err error) {
	id, err := util.StringToInt64(req.Id)
	if err != nil {
		return nil, err
	}
	_, err = l.svcCtx.ImageRPC.UpdateImage(l.ctx, &image.UpdateImageReq{
		Id:         id,
		ImageOrder: req.ImageOrder,
		IsMain:     req.IsMain,
	})
	if err != nil {
		return nil, err
	}
	return &types.UpdateImageResp{}, nil
}
