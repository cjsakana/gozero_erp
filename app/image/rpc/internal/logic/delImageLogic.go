package logic

import (
	"context"

	"erp/app/image/rpc/internal/code"
	"erp/app/image/rpc/internal/svc"
	"erp/app/image/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DelImageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	Logger logx.Logger
}

func NewDelImageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelImageLogic {
	return &DelImageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DelImageLogic) DelImage(in *pb.DelImageReq) (*pb.DelImageResp, error) {
	if err := l.svcCtx.ImageModel.Delete(l.ctx, in.Id); err != nil {
		return nil, code.DeleteImageFail
	}
	return &pb.DelImageResp{}, nil
}
