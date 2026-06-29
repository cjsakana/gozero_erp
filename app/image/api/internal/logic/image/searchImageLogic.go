package image

import (
	"context"
	"erp/common/util"

	"erp/app/image/api/internal/svc"
	"erp/app/image/api/internal/types"
	"erp/app/image/rpc/image"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchImageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchImageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchImageLogic {
	return &SearchImageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchImageLogic) SearchImage(req *types.SearchImageReq) (resp *types.SearchImageResp, err error) {
	businessId, err := util.StringToInt64(req.BusinessId)
	if err != nil {
		return nil, err
	}
	ret, err := l.svcCtx.ImageRPC.SearchImage(l.ctx, &image.SearchImageReq{
		Page:         req.Page,
		Limit:        req.Limit,
		BusinessType: req.BusinessType,
		BusinessId:   businessId,
	})
	if err != nil {
		return nil, err
	}
	items := make([]*types.Image, 0)
	if ret != nil && len(ret.Image) > 0 {
		items = make([]*types.Image, 0, len(ret.Image))
		for _, v := range ret.Image {
			items = append(items, &types.Image{
				Id:           util.Int64ToString(v.Id),
				BusinessType: v.BusinessType,
				BusinessId:   util.Int64ToString(v.BusinessId),
				ImageUrl:     v.ImageUrl,
				ImageOrder:   v.ImageOrder,
				IsMain:       v.IsMain,
				UploadedBy:   v.UploadedBy,
				CreatedAt:    v.CreatedAt,
				UpdatedAt:    v.UpdatedAt,
			})
		}
	}
	return &types.SearchImageResp{Image: items}, nil
}
