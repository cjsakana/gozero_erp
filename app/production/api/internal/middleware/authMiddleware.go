package middleware

import (
	"erp/app/auth/rpc/client/auth"
	"erp/app/auth/rpc/pb"
	"erp/app/user/rpc/user"
	"erp/common/interceptors"
	"erp/common/jwt"
	"erp/common/xcode"
	"erp/common/xtypes"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest/httpx"
	"github.com/zeromicro/go-zero/zrpc"
)

type AuthConf struct {
	AccessSecret string
	AccessExpire int64
}

type AuthMiddlewareConfig struct {
	AuthConf    AuthConf
	AuthRPCConf zrpc.RpcClientConf
	RedisConf   redis.RedisConf
}

type AuthMiddleware struct {
	AuthRPC  auth.AuthZrpcClient
	UserRPC  user.UserZrpcClient
	BizRedis *redis.Redis
	Auth     struct {
		AccessSecret string
		AccessExpire int64
	}
}

func NewAuthMiddleware(c AuthMiddlewareConfig) *AuthMiddleware {
	authRPC := zrpc.MustNewClient(c.AuthRPCConf, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))

	return &AuthMiddleware{
		AuthRPC:  auth.NewAuthZrpcClient(authRPC),
		BizRedis: redis.New(c.RedisConf.Host, redis.WithPass(c.RedisConf.Pass)),
		Auth:     AuthConf{AccessSecret: c.AuthConf.AccessSecret, AccessExpire: c.AuthConf.AccessExpire},
	}
}

func (m *AuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// 从 Header 的 Authorization 中获取 Token（OAuth 2.0规范）
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			httpx.ErrorCtx(ctx, w, xcode.Unauthorized)
			return
		}

		// OAuth 2.0规范的要求：Bearer <token>
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader { // 没有Bearer前缀
			httpx.ErrorCtx(ctx, w, xcode.Unauthorized)
			return
		}

		// 黑名单制度，先解析
		err := jwt.ParseToken(tokenString, m.Auth.AccessSecret, r)
		if err != nil {
			httpx.ErrorCtx(ctx, w, xcode.Unauthorized)
			return
		}

		// ParseToken 更新了 r.Context()，需要重新获取更新后的 context
		ctx = r.Context()

		// 后Redis
		jtiValue := ctx.Value(xtypes.JwtId)
		if jtiValue == nil {
			httpx.ErrorCtx(ctx, w, xcode.Unauthorized)
			return
		}
		jti, ok := jtiValue.(string)
		if !ok {
			httpx.ErrorCtx(ctx, w, xcode.Unauthorized)
			return
		}
		exists, _ := m.BizRedis.Exists(fmt.Sprintf(xtypes.CacheJWTBlackKey, jti))
		if exists {
			httpx.ErrorCtx(ctx, w, xcode.Unauthorized)
			return
		}

		// 校验version版本，解决 roleIds问题
		versionValue := ctx.Value(xtypes.TokenVersionKey)
		if versionValue == nil {
			httpx.ErrorCtx(ctx, w, xcode.Unauthorized)
			return
		}
		version, ok := versionValue.(float64)
		if !ok {
			httpx.ErrorCtx(ctx, w, xcode.Unauthorized)
			return
		}

		userIdValue := ctx.Value(xtypes.UserIdKey)
		if userIdValue == nil {
			httpx.ErrorCtx(ctx, w, xcode.Unauthorized)
			return
		}
		// userId现在作为字符串存储在token中，避免float64精度丢失
		userIdStr, ok := userIdValue.(string)
		if !ok {
			httpx.ErrorCtx(ctx, w, xcode.Unauthorized)
			return
		}
		userIdInt64, err := strconv.ParseInt(userIdStr, 10, 64)
		if err != nil {
			httpx.ErrorCtx(ctx, w, xcode.Unauthorized)
			return
		}

		versionInt64 := int64(version)

		// 从 Redis 获取用户当前的 token version
		v, err := m.BizRedis.Get(fmt.Sprintf(xtypes.CacheJWTVersionKey, userIdInt64))
		if err != nil || v == "" {
			// Redis 中没有 version，说明是异常情况，拒绝访问
			httpx.ErrorCtx(ctx, w, xcode.Unauthorized)
			return
		}

		vi, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			httpx.ErrorCtx(ctx, w, xcode.Unauthorized)
			return
		}

		// 如果 token 中的 version 与 Redis 中的不一致，说明 token 已失效（用户重新登录）
		if vi != versionInt64 {
			httpx.ErrorCtx(ctx, w, xcode.Unauthorized)
			return
		}

		// 校验Permission
		ret1, err := m.AuthRPC.GetUserRoleByUserId(ctx, &pb.GetUserRoleByUserIdReq{UserId: userIdInt64})
		if err != nil {
			httpx.ErrorCtx(ctx, w, xcode.Unauthorized)
			return
		}

		path := strings.SplitN(r.URL.String(), "?", 2)[0]
		path = strings.SplitN(path, "#", 2)[0]
		// 处理尾部斜杠，我的请求不存在 /
		path = strings.TrimRight(path, "/")
		method := r.Method

		matched := false

		for _, userRole := range ret1.UserRole {
			roleId := userRole.RoleId
			ret2, err := m.AuthRPC.GetRolePermissionByRoleId(ctx, &pb.GetRolePermissionByRoleIdReq{RoleId: roleId})
			if err != nil {
				httpx.ErrorCtx(ctx, w, xcode.Unauthorized)
				return
			}
			for _, permission := range ret2.RolePermissions.Permissions {
				// 支持动态路径参数匹配，如 /api/v1/role/{id} 可以匹配 /api/v1/role/123
				if m.matchPath(permission.Url, path) && permission.Method == method {
					matched = true
					next.ServeHTTP(w, r.WithContext(ctx))
					return
				}
			}
		}

		// 如果没有匹配到权限，拒绝访问
		if !matched {
			httpx.ErrorCtx(ctx, w, xcode.Unauthorized)
			return
		}
	}
}

// matchPath 匹配URL路径，支持动态路径参数
// pattern: 权限表中的URL模式，如 /api/v1/role/{id}
// path: 实际请求的URL路径，如 /api/v1/role/123
// 返回: 是否匹配
func (m *AuthMiddleware) matchPath(pattern, path string) bool {
	// 完全匹配
	if pattern == path {
		return true
	}

	// 处理动态路径参数，将 {param} 转换为正则表达式
	// 例如：/api/v1/role/{id} -> /api/v1/role/([^/]+)
	if strings.Contains(pattern, "{") && strings.Contains(pattern, "}") {
		// 转义特殊字符
		patternRegex := regexp.QuoteMeta(pattern)
		// 将 {xxx} 替换为匹配任意非斜杠字符的正则表达式
		patternRegex = regexp.MustCompile(`\\\{[^}]+\}`).ReplaceAllString(patternRegex, `([^/]+)`)

		// 添加开始和结束锚点，确保完全匹配
		patternRegex = "^" + patternRegex + "$"

		matched, err := regexp.MatchString(patternRegex, path)
		if err == nil && matched {
			return true
		}
	}

	return false
}
