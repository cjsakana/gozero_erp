package rolePermission

import (
	"context"
	"erp/app/auth/api/internal/svc"
	"erp/app/auth/api/internal/types"
	"erp/app/auth/rpc/pb"
	"erp/common/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddRolePermissionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddRolePermissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddRolePermissionLogic {
	return &AddRolePermissionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddRolePermissionLogic) AddRolePermission(req *types.AddRolePermissionRequest) (resp *types.AddRolePermissionResponse, err error) {
	roleId, err := util.StringToInt64(req.RoleId)
	if err != nil {
		return nil, err
	}
	permissionId, err := util.StringToInt64(req.PermissionId)
	if err != nil {
		return nil, err
	}
	
	ret, err := l.svcCtx.RolePermissionRPC.AddRolePermission(l.ctx, &pb.AddRolePermissionReq{
		RoleId:       roleId,
		PermissionId: permissionId,
	})
	if err != nil {
		return nil, err
	}
	resp = &types.AddRolePermissionResponse{
		Id: util.Int64ToString(ret.Id),
	}
	return
}
