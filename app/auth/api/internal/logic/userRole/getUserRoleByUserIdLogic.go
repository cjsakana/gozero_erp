package userRole

import (
	"context"
	"erp/app/auth/api/internal/svc"
	"erp/app/auth/api/internal/types"
	"erp/app/auth/rpc/pb"
	"erp/common/util"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserRoleByUserIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserRoleByUserIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserRoleByUserIdLogic {
	return &GetUserRoleByUserIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserRoleByUserIdLogic) GetUserRoleByUserId(req *types.GetUserRoleByUserIdRequest) (resp *types.GetUserRoleByUserIdResponse, err error) {
	userId, err := util.StringToInt64(req.UserId)
	if err != nil {
		return nil, err
	}

	ret, err := l.svcCtx.UserRoleRPC.GetUserRoleByUserId(l.ctx, &pb.GetUserRoleByUserIdReq{
		UserId: userId,
	})

	if err != nil {
		return nil, err
	}
	resp = &types.GetUserRoleByUserIdResponse{
		Total:    ret.Total,
		UserRole: make([]types.UserRole, 0), // 初始化为空切片而不是 nil
	}
	for _, v := range ret.UserRole {
		byIdResp, err := l.svcCtx.RoleRPC.GetRoleById(l.ctx, &pb.GetRoleByIdReq{
			Id: v.RoleId,
		})
		if err != nil {
			return nil, err
		}
		resp.UserRole = append(resp.UserRole, types.UserRole{
			Id:          util.Int64ToString(v.Id),
			UserId:      util.Int64ToString(v.UserId),
			RoleId:      util.Int64ToString(v.RoleId),
			Code:        byIdResp.Role.Code,
			Name:        byIdResp.Role.Name,
			Description: byIdResp.Role.Description,
		})
	}

	return
}
