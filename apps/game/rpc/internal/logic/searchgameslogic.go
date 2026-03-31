package logic

import (
	"context"

	"dicetales.com/apps/game/rpc/game"
	"dicetales.com/apps/game/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchGamesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchGamesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchGamesLogic {
	return &SearchGamesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 模糊搜索
func (l *SearchGamesLogic) SearchGames(in *game.SearchGamesReq) (*game.SearchGamesResp, error) {
	// todo: add your logic here and delete this line

	return &game.SearchGamesResp{}, nil
}
