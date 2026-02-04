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
	ErrGroupReqAlreadyPass   = errorx.New(errorx.BUSINESS_LOGIC_ERROR, "群请求已通过")
	ErrGroupReqAlreadyRefuse = errorx.New(errorx.BUSINESS_LOGIC_ERROR, "群请求已拒绝")
)

type GroupPutInHandleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupPutInHandleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupPutInHandleLogic {
	return &GroupPutInHandleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GroupPutInHandleLogic) GroupPutInHandle(in *social.GroupPutInHandleReq) (*social.GroupPutInHandleResp, error) {

	// 查找群申请
	groupReq, err := l.svcCtx.GroupRequestModel.FindOne(l.ctx, uint64(in.GroupReqId))
	if err != nil {
		return nil, errors.Wrapf(errorx.NewDBErr(), "find group req err [%v] req [%v]", err, in.GroupReqId)
	}

	// 检测是否已处理
	switch HandlerResult(groupReq.HandleResult.Int64) {
	case PassHandlerResult:
		return nil, errors.WithStack(ErrGroupReqAlreadyPass)
	case RefuseHandlerResult:
		return nil, errors.WithStack(ErrGroupReqAlreadyRefuse)
	}

	// 更新处理结果
	groupReq.HandleResult = sql.NullInt64{Int64: int64(in.HandleResult), Valid: true}
	groupReq.HandleUid = sql.NullString{String: in.HandleUid, Valid: true}
	groupReq.HandleTime = sql.NullTime{Time: time.Unix(time.Now().Unix(), 0), Valid: true}

	err = l.svcCtx.GroupRequestModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		if err := l.svcCtx.GroupRequestModel.UpdateTx(l.ctx, session, groupReq); err != nil {
			return errors.Wrapf(errorx.NewDBErr(), "update group req err %v req %v", err, groupReq)
		}

		// 未通过，直接返回
		if HandlerResult(groupReq.HandleResult.Int64) != PassHandlerResult {
			return nil
		}

		// 通过，插入群成员
		groupMember := &model.GroupMember{
			GroupId:     groupReq.GroupId,
			Uid:         groupReq.ReqUid,
			RoleLevel:   int64(MemberGroupRoleLevel),
			JoinTime:    sql.NullTime{Time: time.Unix(time.Now().Unix(), 0), Valid: true},
			JoinSource:  groupReq.JoinSource,
			OperatorUid: sql.NullString{String: in.HandleUid, Valid: true},
		}

		_, err = l.svcCtx.GroupMemberModel.InsertTx(l.ctx, session, groupMember)
		if err != nil {
			return errors.Wrapf(errorx.NewDBErr(), "insert group member err [%v] req [%v]", err, groupMember)
		}
		return nil
	})

	if HandlerResult(groupReq.HandleResult.Int64) != PassHandlerResult {
		return &social.GroupPutInHandleResp{}, err
	}

	return &social.GroupPutInHandleResp{
		GroupId: groupReq.GroupId,
	}, err
}
