package logic

import (
	"context"

	"erp/app/image/rpc/internal/code"
	"erp/app/image/rpc/internal/svc"
	"erp/app/image/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchImageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	Logger logx.Logger
}

func NewSearchImageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchImageLogic {
	return &SearchImageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SearchImageLogic) SearchImage(in *pb.SearchImageReq) (*pb.SearchImageResp, error) {
	list, err := l.svcCtx.ImageModel.FindByBiz(l.ctx, in.BusinessType, in.BusinessId, in.Page, in.Limit)
	if err != nil {
		return nil, code.SearchImageFail
	}
	items := make([]*pb.Image, 0, len(list))
	for _, v := range list {
		items = append(items, &pb.Image{
			Id:           v.Id,
			BusinessType: v.BusinessType,
			BusinessId:   v.BusinessId,
			ImageUrl:     v.ImageUrl,
			ImageOrder:   v.ImageOrder,
			IsMain:       v.IsMain,
			UploadedBy:   v.UploadedBy,
			CreatedAt:    v.UpdatedAt.Unix(),
			UpdatedAt:    v.UpdatedAt.Unix(),
		})
	}
	return &pb.SearchImageResp{Image: items}, nil
}
