package user

import (
	"context"

	"dicetales.com/apps/api/internal/svc"
	"dicetales.com/apps/api/internal/types"
	"dicetales.com/apps/user/rpc/userclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefreshTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 刷新Token
func NewRefreshTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshTokenLogic {
	return &RefreshTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RefreshTokenLogic) RefreshToken(req *types.RefreshTokenReq) (resp *types.RefreshTokenResp, err error) {
	rpcResp, err := l.svcCtx.UserRpc.RefreshToken(l.ctx, &userclient.RefreshTokenReq{
		RefreshToken: req.RefreshToken,
	})
	if err != nil {
		return nil, err
	}

	return &types.RefreshTokenResp{
		AccessToken:   rpcResp.AccessToken,
		AccessExpire:  rpcResp.AccessExpire,
		RefreshToken:  rpcResp.RefreshToken,
		RefreshExpire: rpcResp.RefreshExpire,
	}, nil
}
