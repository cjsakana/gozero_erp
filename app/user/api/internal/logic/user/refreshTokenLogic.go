package user

import (
	"context"
	"erp/app/auth/rpc/pb"
	"erp/app/user/api/internal/svc"
	"erp/app/user/api/internal/types"
	"erp/common/jwt"
	"erp/common/util"
	"erp/common/xcode"
	"erp/common/xtypes"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefreshTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRefreshTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext, w http.ResponseWriter) *RefreshTokenLogic {
	return &RefreshTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RefreshTokenLogic) RefreshToken() (resp *types.RefreshTokenResp, err error) {
	// 使用辅助函数获取userId和employeeId，避免float64精度丢失
	userId, err := util.GetInt64FromCtx(l.ctx, xtypes.UserIdKey)
	if err != nil {
		return nil, err
	}
	employeeId, err := util.GetInt64FromCtx(l.ctx, xtypes.EmployeeIdKey)
	if err != nil {
		return nil, err
	}

	ret, err := l.svcCtx.UserRoleRPC.GetUserRoleByUserId(l.ctx, &pb.GetUserRoleByUserIdReq{
		UserId: userId,
	})
	if err != nil {
		return nil, err
	}

	var roleIds []int64
	for _, v := range ret.UserRole {
		roleIds = append(roleIds, v.RoleId)
	}

	// 获取或生成 token version（用于版本控制，使旧 token 失效）
	var version int64
	versionStr, err := l.svcCtx.BizRedis.Get(fmt.Sprintf(xtypes.CacheJWTVersionKey, userId))
	if err != nil || versionStr == "" {
		// 如果 Redis 中没有 version，生成新的（使用时间戳）
		version = time.Now().Unix()
	} else {
		// 如果已有 version，递增它（使旧 token 失效）
		oldVersion, _ := strconv.ParseInt(versionStr, 10, 64)
		version = oldVersion + 1
	}

	// 将新 version 存储到 Redis（不设置过期时间，永久存储）
	_ = l.svcCtx.BizRedis.Set(fmt.Sprintf(xtypes.CacheJWTVersionKey, userId), fmt.Sprintf("%d", version))

	// 将userId和employeeId作为字符串存储，避免float64精度丢失
	token, err := jwt.BuildTokens(jwt.TokenOptions{
		AccessSecret: l.svcCtx.Config.Auth.AccessSecret,
		AccessExpire: l.svcCtx.Config.Auth.AccessExpire,
		Fields: map[string]interface{}{
			xtypes.UserIdKey:       strconv.FormatInt(userId, 10),
			xtypes.EmployeeIdKey:   strconv.FormatInt(employeeId, 10),
			xtypes.TokenVersionKey: version,
		},
	})
	if err != nil {
		return nil, xcode.New(500, err.Error())
	}

	// 标记旧token缓存，设为黑名单
	jti := l.ctx.Value(xtypes.JwtId).(string)
	key := fmt.Sprintf(xtypes.CacheJWTBlackKey, jti)
	_ = l.svcCtx.BizRedis.SetexCtx(l.ctx, key, "1", int(l.svcCtx.Config.Auth.AccessExpire))

	// OAuth 2.0方式：在响应体中返回token，不再使用Cookie
	resp = &types.RefreshTokenResp{
		AccessToken:  token.AccessToken,
		AccessExpire: token.AccessExpire,
	}
	return
}
