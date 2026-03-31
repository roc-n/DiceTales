package post

import (
	"context"

	"dicetales.com/apps/api/internal/svc"
	"dicetales.com/apps/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PostPingLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Post Module Placeholder
func NewPostPingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PostPingLogic {
	return &PostPingLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PostPingLogic) PostPing(req *types.PostPlaceholderReq) (resp *types.PostPlaceholderResp, err error) {
	// todo: add your logic here and delete this line

	return
}
