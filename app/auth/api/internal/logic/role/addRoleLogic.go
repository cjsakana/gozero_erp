package role

import (
	"context"
	"erp/app/auth/api/internal/svc"
	"erp/app/auth/api/internal/types"
	"erp/app/auth/rpc/pb"
	"erp/common/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddRoleLogic {
	return &AddRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddRoleLogic) AddRole(req *types.AddRoleRequest) (resp *types.AddRoleResponse, err error) {
	ret, err := l.svcCtx.RoleRPC.AddRole(l.ctx, &pb.AddRoleReq{
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		return nil, err
	}
	resp = &types.AddRoleResponse{
		Id: util.Int64ToString(ret.Id),
	}
	return
}
