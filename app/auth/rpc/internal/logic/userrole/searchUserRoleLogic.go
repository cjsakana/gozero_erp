package userrolelogic

import (
	"context"
	"erp/app/auth/rpc/internal/types"

	"erp/app/auth/rpc/internal/svc"
	"erp/app/auth/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchUserRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchUserRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchUserRoleLogic {
	return &SearchUserRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SearchUserRoleLogic) SearchUserRole(in *pb.SearchUserRoleReq) (*pb.SearchUserRoleResp, error) {
	userRoles, count, err := l.svcCtx.UserRoleModel.SearchUserRoles(l.ctx, &types.SearchUserRoleParams{
		SearchCom: types.SearchCom{
			Page:  in.Page,
			Limit: in.Limit,
		},
		RoleId: in.RoleId,
	})
	if err != nil {
		return nil, err
	}

	var pbUR []*pb.UserRole
	for _, v := range userRoles {
		pbUR = append(pbUR, &pb.UserRole{
			Id:     v.Id,
			UserId: v.UserId.Int64,
			RoleId: v.RoleId.Int64,
		})
	}

	return &pb.SearchUserRoleResp{
		Total:    count,
		UserRole: pbUR,
	}, nil
}
