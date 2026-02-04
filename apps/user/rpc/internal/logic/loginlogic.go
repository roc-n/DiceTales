package logic

import (
	"context"
	"time"

	"dicetales.com/apps/user/model"
	"dicetales.com/apps/user/rpc/internal/svc"
	"dicetales.com/apps/user/rpc/user"
	"dicetales.com/pkg/auth"
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

	// 生成token
	now := time.Now().Unix()
	token, err := auth.GetJwtToken(l.svcCtx.Config.Jwt.AccessSecret, now, l.svcCtx.Config.Jwt.AccessExpire, userInfo.Id)
	if err != nil {
		return nil, errors.Wrapf(errorx.NewInternalErr(), "ctxdata get jwt token err %v", err)

	}

	return &user.LoginResp{
		Token:  token,
		Expire: now + l.svcCtx.Config.Jwt.AccessExpire,
	}, nil
}
