package rolepermissionlogic

import (
	"context"
	"erp/app/auth/rpc/internal/code"
	"erp/app/auth/rpc/internal/model"
	types2 "erp/app/auth/rpc/internal/types"

	"github.com/zeromicro/go-zero/core/mr"

	"erp/app/auth/rpc/internal/svc"
	"erp/app/auth/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchRolePermissionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchRolePermissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchRolePermissionLogic {
	return &SearchRolePermissionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SearchRolePermissionLogic) SearchRolePermission(in *pb.SearchRolePermissionReq) (*pb.SearchRolePermissionResp, error) {
	roles, total, err := l.svcCtx.RoleModel.SearchRoles(l.ctx, &types2.SearchRole{
		SearchCom: types2.SearchCom{
			Page:  in.Page,
			Limit: in.Limit,
		},
		Code:        in.Code,
		Name:        in.Name,
		Description: in.Description,
	})
	if err != nil {
		return nil, code.GetPermissionFail
	}

	type rp struct {
		RoleID      int64
		Permissions []*model.PermissionWithRolePermissionId
	}

	generate := func(source chan<- int64) {
		for _, role := range roles {
			source <- role.Id
		}
	}

	mapper := func(id int64, writer mr.Writer[*rp], cancel func(error)) {
		permissions, err := l.svcCtx.RolePermissionModel.FindPermissionsWithRolePermissionIdByRoleId(l.ctx, id)
		if err != nil {
			return
		}
		writer.Write(&rp{
			RoleID:      id,
			Permissions: permissions,
		})
	}

	reducer := func(pipe <-chan *rp, writer mr.Writer[map[int64][]*model.PermissionWithRolePermissionId], cancel func(error)) {
		result := make(map[int64][]*model.PermissionWithRolePermissionId)
		for p := range pipe {
			result[p.RoleID] = p.Permissions
		}
		writer.Write(result)
	}
	permissionsByRole, err := mr.MapReduce[int64, *rp, map[int64][]*model.PermissionWithRolePermissionId](generate, mapper, reducer)

	resp := pb.SearchRolePermissionResp{
		Total: total,
	}

	resp.RolePermissions = make([]*pb.RoleAndPermissions, 0, len(roles))
	for _, role := range roles {
		perms, exists := permissionsByRole[role.Id]
		if !exists {
			perms = []*model.PermissionWithRolePermissionId{} // 如果没权限，给个空数组，不要 nil
		}

		pbRole := &pb.Role{
			Id:          role.Id,
			Code:        role.Code,
			Name:        role.Name,
			Description: role.Description.String,
		}

		var pbPerms []*pb.PermissionWithRolePermissionId
		for _, perm := range perms {
			pbPerms = append(pbPerms, &pb.PermissionWithRolePermissionId{
				RolePermissionId: perm.RolePermissionId,
				PermissionId:     perm.Id,
				ParentId:         perm.ParentId,
				Code:             perm.Code.String,
				Description:      perm.Description.String,
				Url:              perm.Url.String,
				Method:           perm.Method.String,
			})
		}

		rolePerm := &pb.RoleAndPermissions{
			Role:        pbRole,
			Total:       int64(len(perms)), // 权限总数
			Permissions: pbPerms,
		}

		resp.RolePermissions = append(resp.RolePermissions, rolePerm)
	}

	return &resp, nil
}
