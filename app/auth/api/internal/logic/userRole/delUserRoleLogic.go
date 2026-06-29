package userRole

import (
	"context"
	"erp/app/auth/api/internal/svc"
	"erp/app/auth/api/internal/types"
	"erp/app/auth/rpc/pb"
	"erp/common/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type DelUserRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDelUserRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelUserRoleLogic {
	return &DelUserRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DelUserRoleLogic) DelUserRole(req *types.DelUserRoleRequest) (resp *types.EmptyResponse, err error) {
	id, err := util.StringToInt64(req.Id)
	if err != nil {
		return nil, err
	}
	
	_, err = l.svcCtx.UserRoleRPC.DelUserRole(l.ctx, &pb.DelUserRoleReq{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	return
}
