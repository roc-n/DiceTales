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

type FriendPutInListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFriendPutInListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutInListLogic {
	return &FriendPutInListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FriendPutInListLogic) FriendPutInList(in *social.FriendPutInListReq) (*social.FriendPutInListResp, error) {
	friendReqsOut, err := l.svcCtx.FriendRequestModel.ListNoHandleOut(l.ctx, in.Uid)
	if err != nil {
		return nil, errors.Wrapf(errorx.NewDBErr(), "list friend req err [%v] req [%v]", err, in.Uid)
	}
	friendReqsIn, err := l.svcCtx.FriendRequestModel.ListNoHandleIn(l.ctx, in.Uid)
	if err != nil {
		return nil, errors.Wrapf(errorx.NewDBErr(), "list friend req err [%v] req [%v]", err, in.Uid)
	}

	var reqList []*social.FriendRequest
	var recvList []*social.FriendRequest
	copier.Copy(&reqList, friendReqsOut)
	copier.Copy(&recvList, friendReqsIn)

	return &social.FriendPutInListResp{
		ReqList:  reqList,
		RecvList: recvList,
	}, nil
}
