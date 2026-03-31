package game

import (
	"context"

	"dicetales.com/apps/api/internal/svc"
	"dicetales.com/apps/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GamePingLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Game Module Placeholder
func NewGamePingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GamePingLogic {
	return &GamePingLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GamePingLogic) GamePing(req *types.GamePlaceholderReq) (resp *types.GamePlaceholderResp, err error) {
	// todo: add your logic here and delete this line

	return
}
