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

type GroupListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupListLogic {
	return &GroupListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GroupListLogic) GroupList(in *social.GroupListReq) (*social.GroupListResp, error) {

	// 根据用户ID查询用户所在的群
	userGroups, err := l.svcCtx.GroupMemberModel.ListByUserId(l.ctx, in.Uid)
	if err != nil {
		return nil, errors.Wrapf(errorx.NewDBErr(), "list group members err [%v] req [%v]", err, in.Uid)
	}
	if len(userGroups) == 0 {
		return &social.GroupListResp{}, nil
	}

	// 提取群ID
	ids := make([]string, 0, len(userGroups))
	for _, v := range userGroups {
		ids = append(ids, v.GroupId)
	}

	// 根据群ID查询群信息
	groups, err := l.svcCtx.GroupsModel.ListByGroupIds(l.ctx, ids)
	if err != nil {
		return nil, errors.Wrapf(errorx.NewDBErr(), "list group err [%v] req [%v]", err, ids)
	}

	// 将查询到的群信息转换为RPC响应对象
	var respList []*social.Group
	copier.Copy(&respList, &groups)

	return &social.GroupListResp{
		List: respList,
	}, nil
}
