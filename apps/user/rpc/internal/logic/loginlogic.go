package logic

import (
	"context"

	"dicetales.com/apps/user/model"
	"dicetales.com/apps/user/rpc/internal/svc"
	"dicetales.com/apps/user/rpc/user"
	"dicetales.com/pkg/encrypt"
	"dicetales.com/pkg/errorx"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

var (
	ErrIdNotRegister = errorx.New(errorx.SERVER_COMMON_ERROR, "手机号未注册")
	ErrUserPwdError  = errorx.New(errorx.SERVER_COMMON_ERROR, "密码错误")
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *user.LoginReq) (*user.LoginResp, error) {

	// 根据id验证用户存在
	userInfo, err := l.svcCtx.UserModel.FindById(l.ctx, in.Id)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, errors.WithStack(ErrIdNotRegister)
		}
		return nil, errors.Wrapf(errorx.NewDBErr(), "find user by id err: [%v], req: [%v]", err, in.Id)
	}

	// 密码验证
	if !encrypt.ValidatePasswordHash(in.Password, userInfo.Password) {
		return nil, errors.WithStack(ErrUserPwdError)
	}

	pack, err := issueDualTokens(l.ctx, l.svcCtx, userInfo.Id, l.Logger)
	if err != nil {
		return nil, err
	}

	return &user.LoginResp{
		Token:         pack.AccessToken,  // backwards compat
		Expire:        pack.AccessExpire, // backwards compat
		AccessToken:   pack.AccessToken,
		AccessExpire:  pack.AccessExpire,
		RefreshToken:  pack.RefreshToken,
		RefreshExpire: pack.RefreshExpire,
	}, nil
}
