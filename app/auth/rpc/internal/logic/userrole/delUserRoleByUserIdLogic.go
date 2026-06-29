package userrolelogic

import (
	"context"
	"fmt"

	"erp/app/auth/rpc/internal/svc"
	"erp/app/auth/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DelUserRoleByUserIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDelUserRoleByUserIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelUserRoleByUserIdLogic {
	return &DelUserRoleByUserIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DelUserRoleByUserIdLogic) DelUserRoleByUserId(in *pb.DelUserRoleByUserIdReq) (*pb.DelUserRoleByUserIdReq, error) {
	ids, err := l.svcCtx.UserRoleModel.DeleteByUserId(l.ctx, in.UserId)
	if err != nil {
		return nil, err
	}
	for _, id := range ids {
		cacheErpAuthUserRoleIdPrefix := "cache:erpAuth:userRole:id:"
		erpAuthUserRoleIdKey := fmt.Sprintf("%s%v", cacheErpAuthUserRoleIdPrefix, id)
		_, _ = l.svcCtx.BizRedis.Del(erpAuthUserRoleIdKey)
	}

	return &pb.DelUserRoleByUserIdReq{}, nil
}
