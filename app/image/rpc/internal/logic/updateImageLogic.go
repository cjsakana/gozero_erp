package logic

import (
	"context"

	"erp/app/image/rpc/internal/code"
	"erp/app/image/rpc/internal/svc"
	"erp/app/image/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
)

type UpdateImageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	Logger logx.Logger
}

func NewUpdateImageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateImageLogic {
	return &UpdateImageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateImageLogic) UpdateImage(in *pb.UpdateImageReq) (*pb.UpdateImageResp, error) {
	// 先查出现有记录
	img, err := l.svcCtx.ImageModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if err == sqlc.ErrNotFound {
			return nil, code.ImageNotFound
		}
		return nil, code.GetImageFail
	}
	// 更新可变更字段
	img.ImageOrder = in.ImageOrder
	img.IsMain = in.IsMain
	// 提交更新
	if err := l.svcCtx.ImageModel.Update(l.ctx, img); err != nil {
		return nil, code.UpdateImageFail
	}
	return &pb.UpdateImageResp{}, nil
}
