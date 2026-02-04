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
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var (
	ErrFriendReqAlreadyPass   = errorx.New(errorx.BUSINESS_LOGIC_ERROR, "好友申请已通过")
	ErrFriendReqAlreadyRefuse = errorx.New(errorx.BUSINESS_LOGIC_ERROR, "好友申请已拒绝")
)

type FriendPutInHandleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFriendPutInHandleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutInHandleLogic {
	return &FriendPutInHandleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FriendPutInHandleLogic) FriendPutInHandle(in *social.FriendPutInHandleReq) (*social.FriendPutInHandleResp, error) {

	// 获取好友申请记录
	friendReq, err := l.svcCtx.FriendRequestModel.FindOne(l.ctx, uint64(in.FriendReqId))
	if err != nil {
		return nil, errors.Wrapf(errorx.NewDBErr(), "find friendsRequest by friendReqid err [%v] req [%v] ", err,
			in.FriendReqId)
	}

	// 验证是否有处理
	switch HandlerResult(friendReq.HandleResult.Int64) {
	case PassHandlerResult:
		return nil, errors.WithStack(ErrFriendReqAlreadyPass)
	case RefuseHandlerResult:
		return nil, errors.WithStack(ErrFriendReqAlreadyRefuse)
	}

	friendReq.HandleResult = sql.NullInt64{Int64: int64(in.HandleResult), Valid: true}
	friendReq.HandleMsg = sql.NullString{String: in.HandleMessage, Valid: true}
	friendReq.HandledAt = sql.NullTime{Time: time.Unix(in.HandleTime, 0), Valid: true}

	// 修改申请结果 -> 建立两条好友关系记录 -> 事务
	err = l.svcCtx.FriendRequestModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		if err := l.svcCtx.FriendRequestModel.UpdateTx(l.ctx, session, friendReq); err != nil {
			return errors.Wrapf(errorx.NewDBErr(), "update friend request err [%v], req [%v]", err, friendReq)
		}
		if HandlerResult(in.HandleResult) != PassHandlerResult {
			return nil
		}

		friends := []*model.Friend{
			{
				Uid:       friendReq.Uid,
				FriendUid: friendReq.ReqUid,
				Remark: sql.NullString{
					String: friendReq.ReqRemark,
					Valid:  true,
				},
				AddSource: sql.NullInt64{
					Int64: int64(in.AddSource),
					Valid: true,
				},
			}, {
				Uid:       friendReq.ReqUid,
				FriendUid: friendReq.Uid,
				Remark: sql.NullString{
					String: in.Remark,
					Valid:  true,
				},
				AddSource: sql.NullInt64{
					Int64: int64(in.AddSource),
					Valid: true,
				},
			},
		}

		_, err = l.svcCtx.FriendModel.Inserts(l.ctx, session, friends...)
		if err != nil {
			return errors.Wrapf(errorx.NewDBErr(), "friend inserts err %v, req %v", err, friends)
		}

		return nil
	})

	return &social.FriendPutInHandleResp{}, err
}
