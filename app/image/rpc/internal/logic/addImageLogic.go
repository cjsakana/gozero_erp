package logic

import (
	"context"
	"erp/app/image/rpc/internal/code"
	"erp/app/image/rpc/internal/model"
	"erp/app/image/rpc/internal/svc"
	"erp/app/image/rpc/pb"
	"erp/common/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddImageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	Logger logx.Logger
}

func NewAddImageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddImageLogic {
	return &AddImageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// -----------------------image-----------------------
func (l *AddImageLogic) AddImage(in *pb.AddImageReq) (*pb.AddImageResp, error) {
	id := util.GenerateSnowflake()
	data := &model.Image{
		Id:           id,
		BusinessType: in.BusinessType,
		BusinessId:   in.BusinessId,
		ImageUrl:     in.ImageUrl,
		ImageOrder:   in.ImageOrder,
		IsMain:       in.IsMain,
		UploadedBy:   in.UploadedBy,
	}
	_, err := l.svcCtx.ImageModel.Insert(l.ctx, data)
	if err != nil {
		return nil, code.AddImageFail
	}
	return &pb.AddImageResp{}, nil
}
