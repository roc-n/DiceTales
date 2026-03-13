package handler

import (
	"context"
	"net/http"

	"dicetales.com/apps/im/ws/internal/svc"
	"dicetales.com/pkg/ctxdata"
	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/token"
)

type JwtAuth struct {
	svc    *svc.ServiceContext
	parser *token.TokenParser
	logx.Logger
}

func NewJwtAuth(svc *svc.ServiceContext) *JwtAuth {
	return &JwtAuth{
		svc:    svc,
		parser: token.NewTokenParser(),
		Logger: logx.WithContext(context.Background()),
	}
}

func (j *JwtAuth) Auth(w http.ResponseWriter, r *http.Request) bool {
	tok, err := j.parser.ParseToken(r, j.svc.Config.JwtAuth.AccessSecret, "")
	if err != nil {
		j.Logger.Errorf("parse token error: %v", err)
		return false
	}

	if !tok.Valid {
		return false
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok {
		return false
	}

	// 注意这里不要用自定义类型
	*r = *r.WithContext(context.WithValue(r.Context(), ctxdata.Identify, claims[ctxdata.Identify]))

	return true
}

func (j *JwtAuth) UserId(r *http.Request) string {
	return ctxdata.GetUId(r.Context())
}
