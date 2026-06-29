package rolepermissionlogic

import (
	"context"
	"erp/app/auth/rpc/internal/code"
	"erp/app/auth/rpc/internal/svc"
	"erp/app/auth/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
)

type GetRolePermissionByRoleIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetRolePermissionByRoleIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRolePermissionByRoleIdLogic {
	return &GetRolePermissionByRoleIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetRolePermissionByRoleIdLogic) GetRolePermissionByRoleId(in *pb.GetRolePermissionByRoleIdReq) (*pb.GetRolePermissionByRoleIdResp, error) {
	role, err := l.svcCtx.RoleModel.FindOne(l.ctx, in.RoleId)
	if err != nil {
		if err == sqlc.ErrNotFound {
			return nil, code.RoleNotFound
		}
		return nil, code.GetPermissionFail
	}
	permissions, err := l.svcCtx.RolePermissionModel.FindPermissionsWithRolePermissionIdByRoleId(l.ctx, in.RoleId)
	if err != nil {
		return nil, code.GetPermissionFail
	}
	var rp pb.RoleAndPermissions

	pbRole := &pb.Role{
		Id:          role.Id,
		Code:        role.Code,
		Name:        role.Name,
		Description: role.Description.String,
	}

	rp.Role = pbRole
	rp.Total = int64(len(permissions))
	for _, permission := range permissions {
		rp.Permissions = append(rp.Permissions, &pb.PermissionWithRolePermissionId{
			RolePermissionId: permission.RolePermissionId,
			PermissionId:     permission.Id,
			ParentId:         permission.ParentId,
			Code:             permission.Code.String,
			Description:      permission.Description.String,
			Url:              permission.Url.String,
			Method:           permission.Method.String,
		})
	}

	return &pb.GetRolePermissionByRoleIdResp{
		RolePermissions: &rp,
	}, nil
}
