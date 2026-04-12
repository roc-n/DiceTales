package user

import (
	"context"

	"dicetales.com/apps/api/internal/svc"
	"dicetales.com/apps/api/internal/types"
	"dicetales.com/apps/user/rpc/userclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 用户登入
func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	rpcResp, err := l.svcCtx.UserRpc.Login(l.ctx, &userclient.LoginReq{
		Account:  req.Phone, // Map Phone appropriately if needed, or update LoginReq based on rpc definition
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}

	return &types.LoginResp{
		Token:         rpcResp.Token,
		Expire:        rpcResp.Expire,
		AccessToken:   rpcResp.AccessToken,
		AccessExpire:  rpcResp.AccessExpire,
		RefreshToken:  rpcResp.RefreshToken,
		RefreshExpire: rpcResp.RefreshExpire,
	}, nil
}
