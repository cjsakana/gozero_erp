package jwt

import (
	"context"
	"erp/common/xcode"
	"erp/common/xtypes"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/dgrijalva/jwt-go"
)

var (
	GenTokenFail = xcode.New(500, "Token生成失败，请稍后重新登录")
)

type (
	TokenOptions struct {
		AccessSecret string
		AccessExpire int64
		Fields       map[string]interface{}
	}

	Token struct {
		AccessToken  string `json:"access_token"`
		AccessExpire int64  `json:"access_expire"`
	}
)

// BuildTokens 生成Token
func BuildTokens(opt TokenOptions) (Token, error) {
	var token Token
	now := time.Now().Add(-time.Minute).Unix()
	accessToken, err := genToken(now, opt.AccessSecret, opt.Fields, opt.AccessExpire)
	if err != nil {
		return token, GenTokenFail
	}
	token.AccessToken = accessToken
	token.AccessExpire = now + opt.AccessExpire

	return token, nil
}

func genToken(iat int64, secretKey string, payloads map[string]interface{}, seconds int64) (string, error) {
	claims := make(jwt.MapClaims)
	claims[xtypes.JwtExpire] = iat + seconds
	claims[xtypes.JwtIssueAt] = iat

	// 唯一令牌 ID
	newUUID, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}
	claims[xtypes.JwtId] = newUUID.String()

	for k, v := range payloads {
		claims[k] = v
	}
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims

	return token.SignedString([]byte(secretKey))
}

// ParseToken 解析JWT令牌
func ParseToken(tokenString, secretKey string, r *http.Request) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return err
	}
	if !token.Valid {
		return jwt.ErrInvalidKey
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return jwt.ErrInvalidKeyType
	}

	ctx := r.Context()
	for k, v := range claims {
		switch k {
		case xtypes.JwtAudience, xtypes.JwtExpire, xtypes.JwtIssueAt, xtypes.JwtIssuer, xtypes.JwtNotBefore, xtypes.JwtSubject:
			// ignore the standard claims
		case xtypes.JwtId:
			// jti (JWT ID) 需要被设置到 context 中，用于黑名单检查
			ctx = context.WithValue(ctx, k, v)
		default:
			ctx = context.WithValue(ctx, k, v)
		}
	}
	*r = *r.WithContext(ctx)
	return nil
}
