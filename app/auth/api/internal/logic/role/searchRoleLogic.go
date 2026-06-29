package role

import (
	"context"
	"erp/app/auth/api/internal/svc"
	"erp/app/auth/api/internal/types"
	"erp/app/auth/rpc/pb"
	"erp/common/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchRoleLogic {
	return &SearchRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchRoleLogic) SearchRole(req *types.SearchRoleRequest) (resp *types.SearchRoleResponse, err error) {
	ret, err := l.svcCtx.RoleRPC.SearchRole(l.ctx, &pb.SearchRoleReq{
		Page:        req.Page,
		Limit:       req.Limit,
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		return nil, err
	}
	resp = &types.SearchRoleResponse{
		Total: ret.Total,
		Roles: make([]types.Role, 0, len(ret.Roles)),
	}
	for _, role := range ret.Roles {
		resp.Roles = append(resp.Roles, types.Role{
			Id:          util.Int64ToString(role.Id),
			Code:        role.Code,
			Name:        role.Name,
			Description: role.Description,
		})
	}

	return
}
