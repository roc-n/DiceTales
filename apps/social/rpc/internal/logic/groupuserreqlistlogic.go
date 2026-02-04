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

type GroupUserReqListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupUserReqListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupUserReqListLogic {
	return &GroupUserReqListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GroupUserReqListLogic) GroupUserReqList(in *social.GroupUserReqListReq) (*social.GroupUserReqListResp, error) {
	groupReqs, err := l.svcCtx.GroupRequestModel.ListByReqUid(l.ctx, in.Uid)
	if err != nil {
		return nil, errors.Wrapf(errorx.NewDBErr(), "list group req err [%v] req [%v]", err, in.Uid)
	}

	var respList []*social.GroupRequest
	copier.Copy(&respList, groupReqs)

	return &social.GroupUserReqListResp{
		List: respList,
	}, nil
}
