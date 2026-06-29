package role

import (
	"context"
	"erp/app/auth/api/internal/svc"
	"erp/app/auth/api/internal/types"
	"erp/app/auth/rpc/pb"
	"erp/common/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateRoleLogic {
	return &UpdateRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateRoleLogic) UpdateRole(req *types.UpdateRoleRequest) (resp *types.EmptyResponse, err error) {
	roleId, err := util.StringToInt64(req.Id)
	if err != nil {
		return nil, err
	}
	
	_, err = l.svcCtx.RoleRPC.UpdateRole(l.ctx, &pb.UpdateRoleReq{
		Id:          roleId,
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		return nil, err
	}

	return
}
