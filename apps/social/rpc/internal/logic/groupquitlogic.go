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
	ErrGroupMemberQuit = errorx.New(errorx.BUSINESS_LOGIC_ERROR, "群成员已退出")
)

type GroupQuitLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupQuitLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupQuitLogic {
	return &GroupQuitLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GroupQuitLogic) GroupQuit(in *social.GroupQuitReq) (*social.GroupQuitResp, error) {

	groupMember, err := l.svcCtx.GroupMemberModel.FindByGroudIdAndUserId(l.ctx, in.Uid, in.GroupId)
	if err != nil && err != model.ErrNotFound {
		return nil, errors.Wrapf(errorx.NewDBErr(), "find group member err [%v] req [%v]", err, in.GroupId)
	}
	if err == model.ErrNotFound {
		return nil, errors.WithStack(ErrGroupMemberQuit)
	}

	err = l.svcCtx.GroupMemberModel.Delete(l.ctx, groupMember.Id)
	if err != nil {
		return nil, errors.Wrapf(errorx.NewDBErr(), "delete group member err [%v] req [%v]", err, in.GroupId)
	}

	return &social.GroupQuitResp{
		GroupId: in.GroupId,
	}, nil

}
