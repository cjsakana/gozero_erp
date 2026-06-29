package user

import (
	"context"
	"erp/app/user/api/internal/code"
	"erp/common/util"
	"fmt"

	"erp/app/user/api/internal/svc"
	"erp/app/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendVerifyCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendVerifyCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendVerifyCodeLogic {
	return &SendVerifyCodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

var (
	expireActivation = 60 * 30
)

func (l *SendVerifyCodeLogic) SendVerifyCode(req *types.SendVerifyCodeReq) (resp *types.EmptyResponse, err error) {
	randCode := util.RandomNumeric(6)
	key := fmt.Sprintf(types.VerificationCodeKey, req.Phone)
	err = l.svcCtx.BizRedis.Setex(key, randCode, expireActivation)
	if err != nil {
		return nil, code.SendVerificationCodeFail
	}

	return
}
