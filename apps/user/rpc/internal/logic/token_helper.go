package logic

import (
	"context"
	"fmt"
	"time"

	"dicetales.com/apps/user/rpc/internal/svc"
	"dicetales.com/pkg/auth"
	"dicetales.com/pkg/constants"
	"dicetales.com/pkg/errorx"
	"github.com/google/uuid"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type TokenPack struct {
	AccessToken   string
	AccessExpire  int64
	RefreshToken  string
	RefreshExpire int64
}

// issueDualTokens helper for both login and register
func issueDualTokens(ctx context.Context, svcCtx *svc.ServiceContext, uid string, logger logx.Logger) (*TokenPack, error) {
	now := time.Now().Unix()

	// Access Token
	accessToken, err := auth.GenerateAccessToken(svcCtx.Config.Jwt.AccessSecret, now, svcCtx.Config.Jwt.AccessExpire, uid)
	if err != nil {
		return nil, errors.Wrapf(errorx.NewInternalErr(), "generate access token err %v", err)
	}

	// Refresh Token
	jti := uuid.New().String()
	refreshToken, err := auth.GenerateRefreshToken(svcCtx.Config.Jwt.RefreshSecret, now, svcCtx.Config.Jwt.RefreshExpire, uid, jti)
	if err != nil {
		return nil, errors.Wrapf(errorx.NewInternalErr(), "generate refresh token err %v", err)
	}

	// Redis Whitelist (Single Device Login)
	wlKey := fmt.Sprintf(constants.REDIS_REFRESH_TOKEN_WL, uid)

	// Single device policy: overwrite existing refresh hash directly
	hashedToken := auth.HashToken(refreshToken)
	err = svcCtx.Redis.SetexCtx(ctx, wlKey, hashedToken, int(svcCtx.Config.Jwt.RefreshExpire))
	if err != nil {
		logger.Errorf("redis set refresh token err: %v", err)
		// continue, although it limits token rotation later
	}

	return &TokenPack{
		AccessToken:   accessToken,
		AccessExpire:  now + svcCtx.Config.Jwt.AccessExpire,
		RefreshToken:  refreshToken,
		RefreshExpire: now + svcCtx.Config.Jwt.RefreshExpire,
	}, nil
}
