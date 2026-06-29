package rolePermission

import (
	"context"
	"erp/app/auth/api/internal/svc"
	"erp/app/auth/api/internal/types"
	"erp/app/auth/rpc/pb"
	"erp/common/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRolePermissionByRoleIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetRolePermissionByRoleIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRolePermissionByRoleIdLogic {
	return &GetRolePermissionByRoleIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetRolePermissionByRoleIdLogic) GetRolePermissionByRoleId(req *types.GetRolePermissionByRoleIdRequest) (resp *types.GetRolePermissionByRoleIdResponse, err error) {
	roleId, err := util.StringToInt64(req.RoleId)
	if err != nil {
		return nil, err
	}

	ret, err := l.svcCtx.RolePermissionRPC.GetRolePermissionByRoleId(l.ctx, &pb.GetRolePermissionByRoleIdReq{
		RoleId: roleId,
	})
	if err != nil {
		return nil, err
	}

	rp := types.RoleAndPermissions{
		Total: ret.RolePermissions.Total,
		Role: types.Role{
			Id:          util.Int64ToString(ret.RolePermissions.Role.Id),
			Code:        ret.RolePermissions.Role.Code,
			Name:        ret.RolePermissions.Role.Name,
			Description: ret.RolePermissions.Role.Description,
		},
	}
	for _, v := range ret.RolePermissions.Permissions {
		rp.RolePermissionDetail = append(rp.RolePermissionDetail, types.RolePermissionDetail{
			RolePermissionId: util.Int64ToString(v.RolePermissionId),
			PermissionId:     util.Int64ToString(v.PermissionId),
			ParentId:         util.Int64ToString(v.ParentId),
			Code:             v.Code,
			Description:      v.Description,
			Url:              v.Url,
			Method:           v.Method,
		})
	}
	resp = &types.GetRolePermissionByRoleIdResponse{
		RolePermissions: rp,
	}
	return
}
