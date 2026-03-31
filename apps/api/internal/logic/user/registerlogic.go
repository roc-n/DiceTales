package user

import (
	"context"

	"dicetales.com/apps/api/internal/svc"
	"dicetales.com/apps/api/internal/types"
	"dicetales.com/apps/user/rpc/userclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 用户注册
func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.RegisterResp, err error) {
	rpcResp, err := l.svcCtx.UserRpc.Register(l.ctx, &userclient.RegisterReq{
		Phone:    req.Phone,
		Nickname: req.Nickname,
		Password: req.Password,
		Avatar:   req.Avatar,
		Sex:      int32(req.Sex),
	})
	if err != nil {
		return nil, err
	}

	return &types.RegisterResp{
		Token:         rpcResp.Token,
		Expire:        rpcResp.Expire,
		AccessToken:   rpcResp.AccessToken,
		AccessExpire:  rpcResp.AccessExpire,
		RefreshToken:  rpcResp.RefreshToken,
		RefreshExpire: rpcResp.RefreshExpire,
	}, nil
}
