package logic

import (
	"context"

	"erp/app/user/rpc/internal/svc"
	"erp/app/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserByEmployeeIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserByEmployeeIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserByEmployeeIdLogic {
	return &GetUserByEmployeeIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserByEmployeeIdLogic) GetUserByEmployeeId(in *pb.GetUserByEmployeeIdReq) (*pb.GetUserByEmployeeIdResp, error) {
	user, err := l.svcCtx.UserModel.FindOneByEmployeeId(l.ctx, in.EmployeeId)
	if err != nil {
		return nil, err
	}

	return &pb.GetUserByEmployeeIdResp{
		User: &pb.User{
			Id:         user.Id,
			EmployeeId: user.EmployeeId,
			EmployeeNo: user.EmployeeNo,
			Username:   user.Username,
			RealName:   user.RealName,
			Phone:      user.Phone.String,
			Email:      user.Email.String,
			Avatar:     user.Avatar.String,
			Resigned:   user.Resigned != 0,
			CreatedBy:  user.CreatedBy.Int64,
			CreatedAt:  user.CreatedAt.Unix(),
			UpdatedAt:  user.UpdatedAt.Unix(),
			UpdatedBy:  user.UpdatedBy.Int64,
		},
	}, nil
}
