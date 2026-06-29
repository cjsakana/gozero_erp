package logic

import (
	"context"

	"erp/app/image/rpc/internal/code"
	"erp/app/image/rpc/internal/svc"
	"erp/app/image/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
)

type GetImageByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	Logger logx.Logger
}

func NewGetImageByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetImageByIdLogic {
	return &GetImageByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetImageByIdLogic) GetImageById(in *pb.GetImageByIdReq) (*pb.GetImageByIdResp, error) {
	img, err := l.svcCtx.ImageModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if err == sqlc.ErrNotFound {
			return nil, code.ImageNotFound
		}
		return nil, code.GetImageFail
	}
	resp := &pb.Image{
		Id:           img.Id,
		BusinessType: img.BusinessType,
		BusinessId:   img.BusinessId,
		ImageUrl:     img.ImageUrl,
		ImageOrder:   img.ImageOrder,
		IsMain:       img.IsMain,
		UploadedBy:   img.UploadedBy,
		CreatedAt:    img.CreatedAt.Unix(),
		UpdatedAt:    img.UpdatedAt.Unix(),
	}
	return &pb.GetImageByIdResp{Image: resp}, nil
}
