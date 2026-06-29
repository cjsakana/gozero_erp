package userRole

import (
	"context"
	"erp/app/auth/api/internal/svc"
	"erp/app/auth/api/internal/types"
	"erp/app/auth/rpc/pb"
	"erp/common/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddUserRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddUserRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddUserRoleLogic {
	return &AddUserRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddUserRoleLogic) AddUserRole(req *types.AddUserRoleRequest) (resp *types.AddUserRoleResponse, err error) {
	userId, err := util.StringToInt64(req.UserId)
	if err != nil {
		return nil, err
	}
	roleId, err := util.StringToInt64(req.RoleId)
	if err != nil {
		return nil, err
	}
	
	res, err := l.svcCtx.UserRoleRPC.AddUserRole(l.ctx, &pb.AddUserRoleReq{
		UserId: userId,
		RoleId: roleId,
	})
	if err != nil {
		return nil, err
	}
	resp = &types.AddUserRoleResponse{
		Id: util.Int64ToString(res.Id),
	}
	return
}
