package logic

import (
	"context"

	"dicetales.com/apps/game/rpc/game"
	"dicetales.com/apps/game/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGameBasicInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGameBasicInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGameBasicInfoLogic {
	return &GetGameBasicInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取桌游基础摘要
func (l *GetGameBasicInfoLogic) GetGameBasicInfo(in *game.GameBasicInfoReq) (*game.GameBasicInfoResp, error) {
	// 查询底层数据库，携带 Redis 缓存（缓存由 go-zero model 自动管理）
	gameModel, err := l.svcCtx.GameModel.FindOne(l.ctx, uint64(in.GameId))
	if err != nil {
		l.Logger.Errorf("查询游戏基础信息失败 gameId: %d, err: %v", in.GameId, err)
		return nil, err
	}

	return &game.GameBasicInfoResp{
		Id:       int64(gameModel.Id),
		Name:     gameModel.Name,
		CoverImg: gameModel.CoverImg.String,
		Score:    float64(gameModel.Score),
	}, nil
}
