package logic

import (
	"context"

	"erp/app/user/rpc/internal/svc"
	"erp/app/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserByIdLogic {
	return &GetUserByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserByIdLogic) GetUserById(in *pb.GetUserByIdReq) (*pb.GetUserByIdResp, error) {
	// 查询用户信息
	user, err := l.svcCtx.UserModel.FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	// 构造返回数据
	return &pb.GetUserByIdResp{
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
