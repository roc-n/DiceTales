package logic

import (
	"context"

	"dicetales.com/apps/game/rpc/game"
	"dicetales.com/apps/game/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListGamesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListGamesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListGamesLogic {
	return &ListGamesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 分页查询桌游列表
func (l *ListGamesLogic) ListGames(in *game.ListGamesReq) (*game.ListGamesResp, error) {
	// todo: add your logic here and delete this line

	return &game.ListGamesResp{}, nil
}
