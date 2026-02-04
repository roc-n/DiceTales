package logic

import (
	"context"

	"dicetales.com/apps/social/rpc/internal/svc"
	"dicetales.com/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupNotificationModifyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupNotificationModifyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupNotificationModifyLogic {
	return &GroupNotificationModifyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GroupNotificationModifyLogic) GroupNotificationModify(in *social.GroupNotificationModifyReq) (*social.GroupNotificationModifyResp, error) {

	err := l.svcCtx.GroupsModel.UpdateNotification(l.ctx, in.GroupId, in.Notification, in.NotificationUid)
	if err != nil {
		return nil, err
	}

	return &social.GroupNotificationModifyResp{
		GroupId: in.GroupId,
	}, nil
}
