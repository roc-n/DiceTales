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

var (
	ErrGroupNotFound = errorx.New(errorx.BUSINESS_LOGIC_ERROR, "群聊未找到")
)

type GroupModifyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupModifyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupModifyLogic {
	return &GroupModifyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GroupModifyLogic) GroupModify(in *social.GroupModifyReq) (*social.GroupModifyResp, error) {

	group, err := l.svcCtx.GroupsModel.FindOne(l.ctx, in.GroupId)
	if err != nil && err != model.ErrNotFound {
		return nil, errors.Wrapf(errorx.NewDBErr(), "find group err [%v] req [%v]", err, in.GroupId)
	}
	if err == model.ErrNotFound {
		return nil, errors.WithStack(ErrGroupNotFound)
	}

	group.Name = in.Name
	group.Icon = in.Icon
	group.IsVerify = in.IsVerify
	group.Description = in.Description

	err = l.svcCtx.GroupsModel.Update(l.ctx, group)
	if err != nil {
		return nil, errors.Wrapf(errorx.NewDBErr(), "update group info err [%v] req [%v]", err, in)
	}

	return &social.GroupModifyResp{
		GroupId: in.GroupId,
	}, nil
}
