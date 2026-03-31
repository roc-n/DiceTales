package logic

import (
	"context"

	"dicetales.com/apps/game/rpc/game"
	"dicetales.com/apps/game/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGameDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGameDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGameDetailLogic {
	return &GetGameDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取桌游详细信息
func (l *GetGameDetailLogic) GetGameDetail(in *game.GetGameDetailReq) (*game.GetGameDetailResp, error) {
	// 1. 从缓存/DB获取游戏主记录
	gameInfo, err := l.svcCtx.GameModel.FindOne(l.ctx, uint64(in.GameId))
	if err != nil {
		l.Errorf("FindOne game_id: %d, err: %v", in.GameId, err)
		return nil, err
	}

	// 2. 从缓存/DB获取关联标签ID
	tagIds, err := l.svcCtx.GameTagRelationModel.FindTagIdsByGameId(l.ctx, uint64(in.GameId))
	if err != nil {
		l.Errorf("FindTagIdsByGameId game_id: %d, err: %v", in.GameId, err)
		// 非致命错误，可以继续返回基本信息
	}

	var tags []*game.Tag
	if len(tagIds) > 0 {
		tagRecords, err := l.svcCtx.TagModel.FindByIds(l.ctx, tagIds)
		if err == nil {
			for _, t := range tagRecords {
				tags = append(tags, &game.Tag{
					Id:       int64(t.Id),
					Name:     t.Name,
					Category: int32(t.Category),
				})
			}
		}
	}

	// 3. 从缓存/DB获取关联资源
	resources, err := l.svcCtx.GameResourceModel.FindByGameId(l.ctx, uint64(in.GameId))
	if err != nil {
		l.Errorf("FindByGameId resources game_id: %d, err: %v", in.GameId, err)
	}

	var resList []*game.GameResource
	for _, r := range resources {
		resList = append(resList, &game.GameResource{
			Id:    int64(r.Id),
			Type:  int32(r.Type),
			Title: r.Title,
			Url:   r.Url,
		})
	}

	// 4. 组装响应
	return &game.GetGameDetailResp{
		Id:          int64(gameInfo.Id),
		Name:        gameInfo.Name,
		NameEn:      gameInfo.NameEn,
		CoverImg:    gameInfo.CoverImg,
		Score:       gameInfo.Score,
		ScoreCount:  int32(gameInfo.ScoreCount),
		PublishYear: int32(gameInfo.PublishYear.Int64),
		MinPlayers:  int32(gameInfo.MinPlayers),
		MaxPlayers:  int32(gameInfo.MaxPlayers),
		DurationMin: int32(gameInfo.DurationMin),
		DurationMax: int32(gameInfo.DurationMax),
		Complexity:  gameInfo.Complexity,
		Description: gameInfo.Description.String,
		Tags:        tags,
		Resources:   resList,
	}, nil
}
