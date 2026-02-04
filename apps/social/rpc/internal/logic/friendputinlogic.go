package logic

import (
	"context"
	"database/sql"
	"time"

	"dicetales.com/apps/social/model"
	"dicetales.com/apps/social/rpc/internal/svc"
	"dicetales.com/apps/social/rpc/social"
	"dicetales.com/pkg/errorx"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type FriendPutInLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFriendPutInLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutInLogic {
	return &FriendPutInLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 好友业务：请求好友、通过或拒绝申请、好友列表
func (l *FriendPutInLogic) FriendPutIn(in *social.FriendPutInReq) (*social.FriendPutInResp, error) {
	// todo: add your logic here and delete this line

	// 申请人是否与目标是好友关系
	friends, err := l.svcCtx.FriendModel.FindByUidAndFid(l.ctx, in.Uid, in.ReqUid)
	if err != nil && err != model.ErrNotFound {
		return nil, errors.Wrapf(errorx.NewDBErr(), "find friends by uid and fid err [%v], req [%v]", err, in)
	}
	if friends != nil {
		return &social.FriendPutInResp{}, err
	}

	// 是否已经有过申请，申请是不成功，没有完成
	friendReqs, err := l.svcCtx.FriendRequestModel.FindByReqUidAndUserId(l.ctx, in.ReqUid, in.Uid)
	if err != nil && err != model.ErrNotFound {
		return nil, errors.Wrapf(errorx.NewDBErr(), "find friendsRequest by rid and uid err [%v], req [%v] ", err, in)
	}
	if friendReqs != nil {
		return &social.FriendPutInResp{}, err
	}

	// 创建申请记录
	_, err = l.svcCtx.FriendRequestModel.Insert(l.ctx, &model.FriendRequest{
		Uid:    in.Uid,
		ReqUid: in.ReqUid,
		ReqMsg: sql.NullString{
			Valid:  true,
			String: in.ReqMsg,
		},
		ReqTime: time.Unix(in.ReqTime, 0),
		HandleResult: sql.NullInt64{
			Int64: int64(NoHandlerResult),
			Valid: true,
		},
		ReqRemark: in.Remark,
	})

	if err != nil {
		return nil, errors.Wrapf(errorx.NewDBErr(), "insert friendRequest err [%v], req [%v] ", err, in)
	}

	return &social.FriendPutInResp{}, nil
}
