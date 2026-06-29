package role

import (
	"context"
	"erp/app/auth/api/internal/svc"
	"erp/app/auth/api/internal/types"
	"erp/app/auth/rpc/pb"
	"erp/common/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type DelRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDelRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelRoleLogic {
	return &DelRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DelRoleLogic) DelRole(req *types.DelRoleRequest) (resp *types.EmptyResponse, err error) {
	roleId, err := util.StringToInt64(req.Id)
	if err != nil {
		return nil, err
	}
	
	_, err = l.svcCtx.RoleRPC.DelRole(l.ctx, &pb.DelRoleReq{
		Id: roleId,
	})
	if err != nil {
		return nil, err
	}
	return
}
