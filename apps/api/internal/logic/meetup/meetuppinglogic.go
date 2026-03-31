package meetup

import (
	"context"

	"dicetales.com/apps/api/internal/svc"
	"dicetales.com/apps/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MeetupPingLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Meetup Module Placeholder
func NewMeetupPingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MeetupPingLogic {
	return &MeetupPingLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MeetupPingLogic) MeetupPing(req *types.MeetupPlaceholderReq) (resp *types.MeetupPlaceholderResp, err error) {
	// todo: add your logic here and delete this line

	return
}
