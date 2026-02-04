package logic

import (
	"context"

	"dicetales.com/apps/social/model"
	"dicetales.com/apps/social/rpc/internal/svc"
	"dicetales.com/apps/social/rpc/social"
	"dicetales.com/pkg/errorx"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type FriendDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFriendDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendDeleteLogic {
	return &FriendDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FriendDeleteLogic) FriendDelete(in *social.FriendDeleteReq) (*social.FriendDeleteResp, error) {

	friend, err := l.svcCtx.FriendModel.FindByUidAndFid(l.ctx, in.Uid, in.FriendUid)
	if err != nil && err != model.ErrNotFound {
		return nil, errors.Wrapf(errorx.NewDBErr(), "find friend err [%v] req [%v]", err, in)
	}
	if err == model.ErrNotFound {
		return nil, errors.WithStack(ErrGroupMemberQuit)
	}

	err = l.svcCtx.GroupMemberModel.Delete(l.ctx, friend.Id)
	if err != nil {
		return nil, errors.Wrapf(errorx.NewDBErr(), "delete friend err [%v] req [%v]", err, in)
	}

	return &social.FriendDeleteResp{}, nil
}
