package logic

import (
	"context"

	"dicetales.com/apps/social/rpc/internal/svc"
	"dicetales.com/apps/social/rpc/social"
	"dicetales.com/pkg/errorx"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type FriendListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFriendListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendListLogic {
	return &FriendListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FriendListLogic) FriendList(in *social.FriendListReq) (*social.FriendListResp, error) {

	friendsList, err := l.svcCtx.FriendModel.ListByUserid(l.ctx, in.Uid)
	if err != nil {
		return nil, errors.Wrapf(errorx.NewDBErr(), "list friend by uid err [%v] req [%v] ", err, in.Uid)
	}

	var respList []*social.Friend
	copier.Copy(&respList, &friendsList)

	return &social.FriendListResp{
		List: respList,
	}, nil
}
