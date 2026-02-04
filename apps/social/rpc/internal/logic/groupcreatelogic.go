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

type GroupCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupCreateLogic {
	return &GroupCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GroupCreateLogic) GroupCreate(in *social.GroupCreateReq) (*social.GroupCreateResp, error) {

	g_id, err := l.svcCtx.IDGen.Get(l.ctx)
	if err != nil {
		return nil, errors.Wrapf(errorx.NewInternalErr(), "generate group id err: [%v]", err)
	}
	group := &model.Groups{
		Id:              g_id,
		Name:            in.Name,
		Icon:            in.Icon,
		Status:          1,
		CreatorUid:      in.CreatorUid,
		GroupType:       int64(in.GroupType),
		IsVerify:        in.IsVerify,
		Notification:    sql.NullString{String: "", Valid: true},
		NotificationUid: sql.NullString{String: "", Valid: true},
		Description:     in.Description,
	}

	err = l.svcCtx.GroupsModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		_, err := l.svcCtx.GroupsModel.InsertTx(l.ctx, session, group)
		if err != nil {
			return errors.Wrapf(errorx.NewDBErr(), "insert group err [%v], req %v", err, in)
		}

		_, err = l.svcCtx.GroupMemberModel.InsertTx(l.ctx, session, &model.GroupMember{
			GroupId:   group.Id,
			Uid:       in.CreatorUid,
			RoleLevel: int64(CreatorGroupRoleLevel),
			JoinTime:  sql.NullTime{Time: time.Now(), Valid: true},
		})
		if err != nil {
			return errors.Wrapf(errorx.NewDBErr(), "insert group member err [%v], req [%v]", err, in)
		}

		// 确认ID已使用 (修改 id_pool 状态)
		if err := l.svcCtx.IDGen.Confirm(context.Background(), g_id); err != nil {
			l.Logger.Errorf("confirm group id failed, id: %s, err: [%v]", g_id, err)
		}

		return nil
	})

	return &social.GroupCreateResp{
		Id: g_id,
	}, err
}
