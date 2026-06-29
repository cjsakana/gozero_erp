package logic

import (
	"context"
	"erp/common/encrypt"
	"errors"

	"erp/app/user/rpc/internal/svc"
	"erp/app/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdatePasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdatePasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePasswordLogic {
	return &UpdatePasswordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdatePasswordLogic) UpdatePassword(in *pb.UpdatePasswordReq) (*pb.UpdatePasswordResp, error) {
	// 通过员工no查询用户
	user, err := l.svcCtx.UserModel.FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	// 验证密码
	if !encrypt.CheckPassword(in.OldPassword, user.PasswordHash) {
		return nil, errors.New("原始密码错误")
	}

	passwordHash, err := encrypt.HashPassword(in.NewPassword)
	if err != nil {
		return nil, err
	}
	user.PasswordHash = passwordHash

	err = l.svcCtx.UserModel.Update(l.ctx, user)
	if err != nil {
		return nil, err
	}

	return &pb.UpdatePasswordResp{}, nil
}
