package image

import (
	"context"
	"erp/common/util"

	"erp/app/image/api/internal/svc"
	"erp/app/image/api/internal/types"
	"erp/app/image/rpc/image"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetImageByIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetImageByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetImageByIdLogic {
	return &GetImageByIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetImageByIdLogic) GetImageById(req *types.GetImageByIdReq) (resp *types.GetImageByIdResp, err error) {
	id, err := util.StringToInt64(req.Id)
	if err != nil {
		return nil, err
	}

	ret, err := l.svcCtx.ImageRPC.GetImageById(l.ctx, &image.GetImageByIdReq{Id: id})
	if err != nil {
		return nil, err
	}
	if ret == nil || ret.Image == nil {
		return &types.GetImageByIdResp{}, nil
	}
	iv := ret.Image

	return &types.GetImageByIdResp{Image: types.Image{
		Id:           util.Int64ToString(iv.Id),
		BusinessType: iv.BusinessType,
		BusinessId:   util.Int64ToString(iv.BusinessId),
		ImageUrl:     iv.ImageUrl,
		ImageOrder:   iv.ImageOrder,
		IsMain:       iv.IsMain,
		UploadedBy:   iv.UploadedBy,
		CreatedAt:    iv.CreatedAt,
		UpdatedAt:    iv.UpdatedAt,
	}}, nil
}
