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

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

var (
	ErrPhoneRegistered = errorx.New(errorx.SERVER_COMMON_ERROR, "手机号已注册")
	ErrPasswordEmpty   = errorx.New(errorx.REQUEST_PARAM_ERROR, "密码不能为空")
)

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *user.RegisterReq) (*user.RegisterResp, error) {

	// 验证手机号是否已注册
	userEntity, err := l.svcCtx.UserModel.FindByPhone(l.ctx, in.Phone)
	if err != nil && err != model.ErrNotFound {
		return nil, errors.Wrapf(errorx.NewDBErr(), "find user by phone err: [%v], phone number: [%v]", err, in.Phone)
	}
	if userEntity != nil {
		return nil, errors.WithStack(ErrPhoneRegistered)
	}

	// 定义用户数据
	if len(in.Password) == 0 {
		return nil, errors.WithStack(ErrPasswordEmpty)
	}

	id, err := l.svcCtx.IDGen.Get(l.ctx)
	if err != nil {
		return nil, errors.Wrapf(errorx.NewInternalErr(), "generate user id err: [%v]", err)
	}

	userEntity = &model.User{
		Id:       id,
		Avatar:   in.Avatar,
		Nickname: in.Nickname,
		Phone:    in.Phone,
		Bio:      in.Bio,
		Sex:      int64(in.Sex),
		City:     in.City,
	}
	genPassword, err := encrypt.GenPasswordHash([]byte(in.Password))
	if err != nil {
		return nil, err
	}
	userEntity.Password = string(genPassword)

	// 插入用户数据
	_, err = l.svcCtx.UserModel.Insert(l.ctx, userEntity)
	if err != nil {
		return nil, errors.Wrapf(errorx.NewDBErr(), "insert user err [%v], req: [%v]", err, in)
	}

	// 确认ID已使用 (修改 id_pool 状态)
	if err := l.svcCtx.IDGen.Confirm(context.Background(), id); err != nil {
		l.Logger.Errorf("confirm user id failed, id: %s, err: [%v]", id, err)
	}

	// 生成token
	now := time.Now().Unix()
	token, err := auth.GetJwtToken(l.svcCtx.Config.Jwt.AccessSecret, now, l.svcCtx.Config.Jwt.AccessExpire, userEntity.Id)
	if err != nil {
		return nil, errors.Wrapf(errorx.NewInternalErr(), "auth get jwt token err: [%v]", err)
	}

	return &user.RegisterResp{
		Token:  token,
		Expire: now + l.svcCtx.Config.Jwt.AccessExpire,
		Id:     id,
	}, nil
}
