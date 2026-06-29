package userrolelogic

import (
	"context"

	"erp/app/auth/rpc/internal/code"
	"erp/app/auth/rpc/internal/svc"
	"erp/app/auth/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DelUserRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDelUserRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelUserRoleLogic {
	return &DelUserRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DelUserRoleLogic) DelUserRole(in *pb.DelUserRoleReq) (*pb.DelUserRoleResp, error) {
	err := l.svcCtx.UserRoleModel.Delete(l.ctx, in.Id)
	if err != nil {
		return nil, code.DeleteUserRoleFail
	}

	return &pb.DelUserRoleResp{}, nil
}
