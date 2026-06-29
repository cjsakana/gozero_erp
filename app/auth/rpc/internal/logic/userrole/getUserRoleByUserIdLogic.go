package userrolelogic

import (
	"context"
	"erp/app/auth/rpc/internal/code"
	"erp/app/auth/rpc/internal/svc"
	"erp/app/auth/rpc/pb"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserRoleByUserIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserRoleByUserIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserRoleByUserIdLogic {
	return &GetUserRoleByUserIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserRoleByUserIdLogic) GetUserRoleByUserId(in *pb.GetUserRoleByUserIdReq) (*pb.GetUserRoleByUserIdResp, error) {
	userRoles, err := l.svcCtx.UserRoleModel.FindRolesByUserId(l.ctx, in.UserId)
	if err != nil {
		return nil, code.GetPermissionFail
	}
	var userRole []*pb.UserRole
	for _, ur := range userRoles {
		// 检查 Valid 字段，确保数据有效
		if !ur.UserId.Valid || !ur.RoleId.Valid {
			l.Errorf("Invalid user role data: id=%d, userId.Valid=%v, roleId.Valid=%v", ur.Id, ur.UserId.Valid, ur.RoleId.Valid)
			continue
		}
		userRole = append(userRole, &pb.UserRole{
			Id:     ur.Id,
			UserId: ur.UserId.Int64,
			RoleId: ur.RoleId.Int64,
		})
	}

	return &pb.GetUserRoleByUserIdResp{
		Total:    int64(len(userRole)),
		UserRole: userRole,
	}, nil
}
