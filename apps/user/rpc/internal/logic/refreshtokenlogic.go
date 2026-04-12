package logic

import (
	"context"
	"fmt"

	"dicetales.com/apps/user/rpc/internal/svc"
	"dicetales.com/apps/user/rpc/user"
	"dicetales.com/pkg/auth"
	"dicetales.com/pkg/constants"
	"dicetales.com/pkg/ctxdata"
	"dicetales.com/pkg/errorx"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

var (
	ErrInvalidRefreshToken = errorx.New(errorx.SERVER_COMMON_ERROR, "invalid refresh token")
)

type RefreshTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRefreshTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshTokenLogic {
	return &RefreshTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RefreshTokenLogic) RefreshToken(in *user.RefreshTokenReq) (*user.RefreshTokenResp, error) {
	// Parse Token
	claims, err := auth.ParseRefreshToken(in.RefreshToken, l.svcCtx.Config.Jwt.RefreshSecret)
	if err != nil {
		return nil, errors.WithStack(ErrInvalidRefreshToken)
	}

	uid, ok := claims[string(ctxdata.Identify)].(string)
	if !ok || uid == "" {
		return nil, errors.WithStack(ErrInvalidRefreshToken)
	}

	// Verify Whitelist
	wlKey := fmt.Sprintf(constants.REDIS_REFRESH_TOKEN_WL, uid)
	hashedToken := auth.HashToken(in.RefreshToken)

	storedHash, err := l.svcCtx.Redis.GetCtx(l.ctx, wlKey)
	if err != nil || storedHash != hashedToken {
		return nil, errors.WithStack(ErrInvalidRefreshToken)
	}

	// Rotate tokens
	pack, err := issueDualTokens(l.ctx, l.svcCtx, uid, l.Logger)
	if err != nil {
		return nil, err
	}

	return &user.RefreshTokenResp{
		AccessToken:   pack.AccessToken,
		AccessExpire:  pack.AccessExpire,
		RefreshToken:  pack.RefreshToken,
		RefreshExpire: pack.RefreshExpire,
	}, nil
}
