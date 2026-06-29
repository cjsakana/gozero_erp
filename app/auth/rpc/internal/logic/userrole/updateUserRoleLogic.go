package userrolelogic

import (
	"context"

	"erp/app/auth/rpc/internal/code"
	"erp/app/auth/rpc/internal/svc"
	"erp/app/auth/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserRoleLogic {
	return &UpdateUserRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateUserRoleLogic) UpdateUserRole(in *pb.UpdateUserRoleReq) (*pb.UpdateUserRoleResp, error) {
	err := l.svcCtx.UserRoleModel.XUpdate(l.ctx, in.Id, in.RoleId, in.UserId)
	if err != nil {
		return nil, code.UpdateUserRoleFail
	}
	return &pb.UpdateUserRoleResp{}, nil
}
