package userRole

import (
	"context"
	"erp/app/auth/api/internal/svc"
	"erp/app/auth/api/internal/types"
	"erp/app/auth/rpc/pb"
	"erp/common/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateUserRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserRoleLogic {
	return &UpdateUserRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateUserRoleLogic) UpdateUserRole(req *types.UpdateUserRoleRequest) (resp *types.EmptyResponse, err error) {
	id, err := util.StringToInt64(req.Id)
	if err != nil {
		return nil, err
	}
	userId, err := util.StringToInt64(req.UserId)
	if err != nil {
		return nil, err
	}
	roleId, err := util.StringToInt64(req.RoleId)
	if err != nil {
		return nil, err
	}
	
	_, err = l.svcCtx.UserRoleRPC.UpdateUserRole(l.ctx, &pb.UpdateUserRoleReq{
		Id:     id,
		UserId: userId,
		RoleId: roleId,
	})
	if err != nil {
		return nil, err
	}

	return
}
