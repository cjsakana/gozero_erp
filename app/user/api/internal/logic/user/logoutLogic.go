package user

import (
	"context"
	"erp/common/xtypes"
	"fmt"
	"net/http"

	"erp/app/user/api/internal/svc"
	"erp/app/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LogoutLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	writer http.ResponseWriter
}

func NewLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext, writer http.ResponseWriter) *LogoutLogic {
	return &LogoutLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		writer: writer,
	}
}

func (l *LogoutLogic) Logout() (resp *types.EmptyResponse, err error) {
	// 标记缓存
	jti := l.ctx.Value(xtypes.JwtId).(string)
	key := fmt.Sprintf(xtypes.CacheJWTBlackKey, jti)
	err = l.svcCtx.BizRedis.SetexCtx(l.ctx, key, "1", int(l.svcCtx.Config.Auth.AccessExpire))

	if err != nil {
		return nil, err
	}
	return &types.EmptyResponse{}, nil
}
