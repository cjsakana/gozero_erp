package logic

import (
	"context"
	"database/sql"
	"erp/common/encrypt"

	"erp/app/user/rpc/internal/svc"
	"erp/app/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ForgotPasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewForgotPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ForgotPasswordLogic {
	return &ForgotPasswordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ForgotPasswordLogic) ForgotPassword(in *pb.ForgotPasswordReq) (*pb.ForgotPasswordResp, error) {
	user, err := l.svcCtx.UserModel.FindOneByPhone(l.ctx, sql.NullString{String: in.Phone, Valid: true})
	if err != nil {
		return nil, err
	}

	passwordHash, err := encrypt.HashPassword(in.Password)
	if err != nil {
		return nil, err
	}
	user.PasswordHash = passwordHash

	err = l.svcCtx.UserModel.Update(l.ctx, user)
	if err != nil {
		return nil, err
	}

	return &pb.ForgotPasswordResp{}, nil
}
