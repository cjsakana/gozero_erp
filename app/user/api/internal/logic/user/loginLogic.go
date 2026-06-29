package user

import (
	"context"
	"erp/app/user/api/internal/code"
	"erp/app/user/api/internal/svc"
	"erp/app/user/api/internal/types"
	"erp/app/user/rpc/pb"
	"erp/common/jwt"
	"erp/common/util"
	"erp/common/xtypes"
	"fmt"
	"net/http"
	"strconv"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx     context.Context
	svcCtx  *svc.ServiceContext
	writer  http.ResponseWriter
	request *http.Request
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext, writer http.ResponseWriter, request *http.Request) *LoginLogic {
	return &LoginLogic{
		Logger:  logx.WithContext(ctx),
		ctx:     ctx,
		svcCtx:  svcCtx,
		writer:  writer,
		request: request,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	ret, err := l.svcCtx.UserRPC.Login(l.ctx, &pb.LoginReq{
		EmployeeNo: req.EmployeeNo,
		Password:   req.Password,
	})
	if err != nil {
		return nil, code.LoginFail
	}
	if ret.Id == -1 {
		return nil, code.LoginFail
	}

	// 获取或生成 token version（用于版本控制，使旧 token 失效）
	var version int64
	versionStr, err := l.svcCtx.BizRedis.Get(fmt.Sprintf(xtypes.CacheJWTVersionKey, ret.Id))
	if err != nil || versionStr == "" {
		// 如果 Redis 中没有 version，生成新的（使用时间戳）
		version = 0
	} else {
		// 如果已有 version，递增它（使旧 token 失效）
		oldVersion, _ := strconv.ParseInt(versionStr, 10, 64)
		version = oldVersion + 1
	}

	// 将新 version 存储到 Redis（不设置过期时间，永久存储）
	_ = l.svcCtx.BizRedis.Set(fmt.Sprintf(xtypes.CacheJWTVersionKey, ret.Id), fmt.Sprintf("%d", version))

	userById, err := l.svcCtx.UserRPC.GetUserById(l.ctx, &pb.GetUserByIdReq{
		Id: ret.Id,
	})
	if err != nil {
		return nil, code.LoginFail
	}
	// 将ID转换为字符串存储，避免JWT解析时float64精度丢失
	token, err := jwt.BuildTokens(jwt.TokenOptions{
		AccessSecret: l.svcCtx.Config.Auth.AccessSecret,
		AccessExpire: l.svcCtx.Config.Auth.AccessExpire,
		Fields: map[string]interface{}{
			xtypes.UserIdKey:       strconv.FormatInt(ret.Id, 10),
			xtypes.EmployeeIdKey:   strconv.FormatInt(userById.User.EmployeeId, 10),
			xtypes.TokenVersionKey: version,
		},
	})

	if err != nil {
		return nil, code.LoginFail
	}

	// OAuth 2.0方式：在响应体中返回token，不再使用Cookie
	resp = &types.LoginResp{
		Id:           util.Int64ToString(ret.Id), // int64 -> string
		AccessToken:  token.AccessToken,
		AccessExpire: token.AccessExpire,
	}
	return
}
