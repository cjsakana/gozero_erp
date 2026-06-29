package rolepermissionlogic

import (
	"context"

	"erp/app/auth/rpc/internal/code"
	"erp/app/auth/rpc/internal/svc"
	"erp/app/auth/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DelRolePermissionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDelRolePermissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelRolePermissionLogic {
	return &DelRolePermissionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DelRolePermissionLogic) DelRolePermission(in *pb.DelRolePermissionReq) (*pb.DelRolePermissionResp, error) {
	err := l.svcCtx.RolePermissionModel.Delete(l.ctx, in.Id)
	if err != nil {
		return nil, code.DeleteRolePermissionFail
	}

	return &pb.DelRolePermissionResp{}, nil
}
