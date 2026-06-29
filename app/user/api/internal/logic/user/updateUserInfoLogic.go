package user

import (
	"context"
	"erp/app/user/api/internal/svc"
	"erp/app/user/api/internal/types"
	"erp/app/user/rpc/pb"
	"erp/common/util"
	"erp/common/xtypes"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserInfoLogic {
	return &UpdateUserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateUserInfoLogic) UpdateUserInfo(req *types.UpdateUserInfoReq) (resp *types.EmptyResponse, err error) {
	userId, err := util.GetInt64FromCtx(l.ctx, xtypes.UserIdKey)
	if err != nil {
		return nil, err
	}
	updateReq := &pb.UpdateUserReq{
		Id:       userId,
		Username: req.Username,
		RealName: req.RealName,
		Phone:    req.Phone,
		Email:    req.Email,
	}

	_, err = l.svcCtx.UserRPC.UpdateUser(l.ctx, updateReq)
	if err != nil {
		return nil, err
	}

	return &types.EmptyResponse{}, nil
}
