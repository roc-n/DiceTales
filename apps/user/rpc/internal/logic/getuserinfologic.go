package logic

import (
	"context"

	"dicetales.com/apps/user/model"
	"dicetales.com/apps/user/rpc/internal/svc"
	"dicetales.com/apps/user/rpc/user"
	"dicetales.com/pkg/errorx"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

var ErrUserNotFound = errorx.New(errorx.BUSINESS_LOGIC_ERROR, "未找到用户")

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserInfoLogic) GetUserInfo(in *user.GetUserInfoReq) (*user.GetUserInfoResp, error) {

	userEntiy, err := l.svcCtx.UserModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	var resp user.UserEntity
	copier.Copy(&resp, userEntiy)

	return &user.GetUserInfoResp{
		User: &resp,
	}, nil
}
