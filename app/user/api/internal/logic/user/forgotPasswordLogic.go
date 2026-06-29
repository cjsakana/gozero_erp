package user

import (
	"context"
	"erp/app/user/api/internal/code"
	"erp/app/user/rpc/pb"
	"fmt"

	"erp/app/user/api/internal/svc"
	"erp/app/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ForgotPasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewForgotPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ForgotPasswordLogic {
	return &ForgotPasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ForgotPasswordLogic) ForgotPassword(req *types.ForgotPasswordReq) (resp *types.EmptyResponse, err error) {
	key := fmt.Sprintf(types.VerificationCodeKey, req.Phone)
	cacheCode, _ := l.svcCtx.BizRedis.Get(key)
	if cacheCode == "" && cacheCode != req.VerifyCode {
		return nil, code.VerificationCodeInvalid
	}
	l.svcCtx.BizRedis.Del(key)

	_, err = l.svcCtx.UserRPC.ForgotPassword(l.ctx, &pb.ForgotPasswordReq{
		Phone:    req.Phone,
		Password: req.Password,
	})
	if err != nil {
		return nil, code.ForgotPasswordFail
	}

	return
}
