package user

import (
	"context"
	"erp/app/user/api/internal/code"
	"erp/app/user/rpc/user"
	"erp/common/util"
	"erp/common/xtypes"

	"erp/app/user/api/internal/svc"
	"erp/app/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdatePasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdatePasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePasswordLogic {
	return &UpdatePasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdatePasswordLogic) UpdatePassword(req *types.UpdatePasswordReq) (resp *types.EmptyResponse, err error) {
	if req.OldPassword == req.NewPassword {
		return nil, code.SamePasswordTwice
	}
	userid, err := util.GetInt64FromCtx(l.ctx, xtypes.UserIdKey)
	if err != nil {
		return nil, err
	}
	_, err = l.svcCtx.UserRPC.UpdatePassword(l.ctx, &user.UpdatePasswordReq{
		Id:          userid,
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	})
	if err != nil {
		return nil, err
	}
	return &types.EmptyResponse{}, nil
}
