package logic

import (
	"context"

	"dicetales.com/apps/user/model"
	"dicetales.com/apps/user/rpc/internal/svc"
	"dicetales.com/apps/user/rpc/user"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type FindUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindUserLogic {
	return &FindUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindUserLogic) FindUser(in *user.FindUserReq) (*user.FindUserResp, error) {

	var (
		userEntitys []*model.User
		userEntity  *model.User
		err         error
	)

	switch {
	case in.Id != "":
		userEntity, err = l.svcCtx.UserModel.FindById(l.ctx, in.Id)
	case in.Phone != "":
		userEntity, err = l.svcCtx.UserModel.FindByPhone(l.ctx, in.Phone)
	case in.Name != "":
		userEntitys, err = l.svcCtx.UserModel.ListByName(l.ctx, in.Name)
	}

	if err != nil {
		return nil, err
	}

	if userEntity != nil {
		userEntitys = append(userEntitys, userEntity)
	}

	var resp []*user.UserEntity
	copier.Copy(&resp, &userEntitys)

	return &user.FindUserResp{
		User: resp,
	}, nil
}
