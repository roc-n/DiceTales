package logic

import (
	"context"

	"dicetales.com/apps/post/rpc/internal/svc"
	"dicetales.com/apps/post/rpc/post"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPublicTimelineLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetPublicTimelineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPublicTimelineLogic {
	return &GetPublicTimelineLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetPublicTimelineLogic) GetPublicTimeline(in *post.GetPublicTimelineReq) (*post.GetPublicTimelineResp, error) {
	// todo: add your logic here and delete this line

	return &post.GetPublicTimelineResp{}, nil
}
