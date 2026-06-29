package user

import (
	"context"
	"erp/app/user/api/internal/svc"
	"erp/app/user/api/internal/types"
	"erp/app/user/rpc/pb"
	"erp/common/util"
	"erp/common/xtypes"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMyUserByIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMyUserByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMyUserByIdLogic {
	return &GetMyUserByIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMyUserByIdLogic) GetMyUserById() (resp *types.GetMyUserByIdResp, err error) {
	id, err := util.GetInt64FromCtx(l.ctx, xtypes.UserIdKey)
	if err != nil {
		return nil, err
	}
	ret, err := l.svcCtx.UserRPC.GetUserById(l.ctx, &pb.GetUserByIdReq{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	return &types.GetMyUserByIdResp{
		User: types.User{
			Id:         util.Int64ToString(ret.User.Id),
			EmployeeId: util.Int64ToString(ret.User.EmployeeId),
			EmployeeNo: ret.User.EmployeeNo,
			Username:   ret.User.Username,
			RealName:   ret.User.RealName,
			Phone:      ret.User.Phone,
			Email:      ret.User.Email,
			Avatar:     ret.User.Avatar,
			Resigned:   false,
			CreatedBy:  util.Int64ToString(ret.User.CreatedBy),
			CreatedAt:  ret.User.CreatedAt,
			UpdatedAt:  ret.User.UpdatedAt,
			UpdatedBy:  util.Int64ToString(ret.User.UpdatedBy),
		},
	}, nil
}
