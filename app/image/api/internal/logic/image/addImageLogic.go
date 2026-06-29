package image

import (
	"context"
	"erp/common/util"

	"erp/app/image/api/internal/svc"
	"erp/app/image/api/internal/types"
	"erp/app/image/rpc/image"
	"erp/common/xtypes"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddImageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddImageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddImageLogic {
	return &AddImageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddImageLogic) AddImage(req *types.AddImageReq) (resp *types.AddImageResp, err error) {
	employeeIdKey, err := util.GetInt64FromCtx(l.ctx, xtypes.EmployeeIdKey)

	businessId, err := util.StringToInt64(req.BusinessId)
	if err != nil {
		return nil, err
	}

	in := &image.AddImageReq{
		BusinessType: req.BusinessType,
		BusinessId:   businessId,
		ImageUrl:     req.ImageUrl,
		ImageOrder:   req.ImageOrder,
		IsMain:       req.IsMain,
		UploadedBy:   employeeIdKey,
	}
	_, err = l.svcCtx.ImageRPC.AddImage(l.ctx, in)
	if err != nil {
		return nil, err
	}

	return &types.AddImageResp{}, nil
}
