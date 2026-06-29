package rolePermission

import (
	"context"
	"erp/app/auth/api/internal/svc"
	"erp/app/auth/api/internal/types"
	"erp/app/auth/rpc/pb"
	"erp/common/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchRolePermissionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchRolePermissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchRolePermissionLogic {
	return &SearchRolePermissionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchRolePermissionLogic) SearchRolePermission(req *types.SearchRolePermissionRequest) (resp *types.SearchRolePermissionResponse, err error) {
	ret, err := l.svcCtx.RolePermissionRPC.SearchRolePermission(l.ctx, &pb.SearchRolePermissionReq{
		Page:        req.Page,
		Limit:       req.Limit,
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		return nil, err
	}
	resp = &types.SearchRolePermissionResponse{
		Total: ret.Total,
	}
	for _, rolePermission := range ret.RolePermissions {
		rp := types.RoleAndPermissions{
			Total: rolePermission.Total,
			Role: types.Role{
				Id:          util.Int64ToString(rolePermission.Role.Id),
				Code:        rolePermission.Role.Code,
				Name:        rolePermission.Role.Name,
				Description: rolePermission.Role.Description,
			},
		}
		for _, v := range rolePermission.Permissions {
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
		resp.RolePermissions = append(resp.RolePermissions, rp)
	}

	return
}
