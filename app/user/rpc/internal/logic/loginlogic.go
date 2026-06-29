package logic

import (
	"context"
	"erp/common/encrypt"
	"errors"

	"erp/app/user/rpc/internal/svc"
	"erp/app/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *pb.LoginReq) (*pb.LoginResp, error) {
	// 通过员工no查询用户
	user, err := l.svcCtx.UserModel.FindOneByEmployeeNo(l.ctx, in.EmployeeNo)
	if err != nil {
		return nil, errors.New("用户不存在")
	}
	// 检查是否已离职
	if user.Resigned != 0 {
		return nil, errors.New("用户不存在")
	}

	// 验证密码
	if !encrypt.CheckPassword(in.Password, user.PasswordHash) {
		return nil, errors.New("密码错误")
	}

	return &pb.LoginResp{
		Id: user.Id,
	}, nil
}
