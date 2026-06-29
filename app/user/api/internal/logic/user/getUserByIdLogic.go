package user

import (
	"context"
	"erp/app/user/api/internal/svc"
	"erp/app/user/api/internal/types"
	"erp/app/user/rpc/pb"
	"erp/common/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserByIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserByIdLogic {
	return &GetUserByIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserByIdLogic) GetUserById(req *types.GetUserByIdReq) (resp *types.GetUserByIdResp, err error) {
	// string ID -> int64
	userId, err := util.StringToInt64(req.Id)
	if err != nil {
		return nil, err
	}

	ret, err := l.svcCtx.UserRPC.GetUserById(l.ctx, &pb.GetUserByIdReq{
		Id: userId,
	})
	if err != nil {
		return nil, err
	}

	return &types.GetUserByIdResp{
		User: types.User{
			Id:         util.Int64ToString(ret.User.Id),        // int64 -> string
			EmployeeId: util.Int64ToString(ret.User.EmployeeId), // int64 -> string
			EmployeeNo: ret.User.EmployeeNo,
			Username:   ret.User.Username,
			RealName:   ret.User.RealName,
			Phone:      ret.User.Phone,
			Email:      ret.User.Email,
			Avatar:     ret.User.Avatar,
			Resigned:   ret.User.Resigned,
			CreatedBy:  util.Int64ToString(ret.User.CreatedBy), // int64 -> string
			CreatedAt:  ret.User.CreatedAt,
			UpdatedAt:  ret.User.UpdatedAt,
			UpdatedBy:  util.Int64ToString(ret.User.UpdatedBy), // int64 -> string
		},
	}, nil
}
