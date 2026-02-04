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

type GroupMemberLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupMemberLogic {
	return &GroupMemberLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GroupMemberLogic) GroupMember(in *social.GroupMemberReq) (*social.GroupMemberResp, error) {

	members, err := l.svcCtx.GroupMemberModel.ListByGroupId(l.ctx, in.GroupId)
	if err != nil {
		return nil, errors.Wrapf(errorx.NewDBErr(), "list group members by group_id err [%v] req [%v] ", err, in.GroupId)
	}

	var respList []*social.GroupMember
	copier.Copy(&respList, &members)

	return &social.GroupMemberResp{
		List: respList,
	}, nil
}
