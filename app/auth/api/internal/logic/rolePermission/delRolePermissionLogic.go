package rolePermission

import (
	"context"
	"erp/app/auth/api/internal/svc"
	"erp/app/auth/api/internal/types"
	"erp/app/auth/rpc/pb"
	"erp/common/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type DelRolePermissionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDelRolePermissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelRolePermissionLogic {
	return &DelRolePermissionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DelRolePermissionLogic) DelRolePermission(req *types.DelRolePermissionRequest) (resp *types.EmptyResponse, err error) {
	id, err := util.StringToInt64(req.Id)
	if err != nil {
		return nil, err
	}
	
	_, err = l.svcCtx.RolePermissionRPC.DelRolePermission(l.ctx, &pb.DelRolePermissionReq{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	return
}
