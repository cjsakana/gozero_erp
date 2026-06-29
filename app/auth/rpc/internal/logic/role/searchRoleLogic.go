package rolelogic

import (
	"context"
	"erp/app/auth/rpc/internal/code"
	"erp/app/auth/rpc/internal/svc"
	types2 "erp/app/auth/rpc/internal/types"
	"erp/app/auth/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchRoleLogic {
	return &SearchRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// rpc GetRoleById(GetRoleByIdReq) returns (GetRoleByIdResp);
func (l *SearchRoleLogic) SearchRole(in *pb.SearchRoleReq) (*pb.SearchRoleResp, error) {

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
	var ret []*pb.Role
	for _, role := range roles {
		ret = append(ret, &pb.Role{
			Id:          role.Id,
			Code:        role.Code,
			Name:        role.Name,
			Description: role.Description.String,
		})
	}

	return &pb.SearchRoleResp{
		Total: total,
		Roles: ret,
	}, nil
}
