package social

import (
	"context"

	"dicetales.com/apps/api/internal/svc"
	"dicetales.com/apps/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SocialPingLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Social Module Placeholder
func NewSocialPingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SocialPingLogic {
	return &SocialPingLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SocialPingLogic) SocialPing(req *types.SocialPlaceholderReq) (resp *types.SocialPlaceholderResp, err error) {
	// todo: add your logic here and delete this line

	return
}
