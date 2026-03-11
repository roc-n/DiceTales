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
func (l *GetGameDetailLogic) GetGameDetail(in *game.GameDetailReq) (*game.GameDetailResp, error) {
	// 1. 获取基本信息（走缓存）
	gameModel, err := l.svcCtx.GameModel.FindOne(l.ctx, uint64(in.GameId))
	if err != nil {
		l.Logger.Errorf("查询游戏详情失败 gameId: %d, err: %v", in.GameId, err)
		return nil, err
	}

	// 2. 获取分类信息（走缓存，FindOneByGameId因为设置了唯一索引，goctl应该生成了此方法）
	categoryModel, err := l.svcCtx.GameCategoryInfoModel.FindOneByGameId(l.ctx, uint64(in.GameId))
	if err != nil {
		l.Logger.Errorf("查询游戏分类信息失败 gameId: %d, err: %v", in.GameId, err)
		// 如果分类信息获取失败，可以不中断主流程，或者按照业务需求处理
	}

	resp := &game.GameDetailResp{
		Id:                    int64(gameModel.Id),
		Name:                  gameModel.Name,
		NameEn:                gameModel.NameEn.String,
		CoverImg:              gameModel.CoverImg.String,
		Score:                 float64(gameModel.Score),
		ScoreCount:            int32(gameModel.ScoreCount),
		MinPlayers:            int32(gameModel.MinPlayers),
		MaxPlayers:            int32(gameModel.MaxPlayers),
		MinRecommendedPlayers: int32(gameModel.MinRecommendedPlayers.Int64),
		MaxRecommendedPlayers: int32(gameModel.MaxRecommendedPlayers.Int64),
		NeedHost:              int32(gameModel.NeedHost),
		RankPosition:          int32(gameModel.Rank.Int64),
		Year:                  int32(gameModel.Year.Int64),
		Description:           gameModel.Description.String,
		Difficulty:            float64(gameModel.Difficulty.Float64),
		DurationPerPlayer:     int32(gameModel.DurationPerPlayer.Int64),
		SetupTime:             gameModel.SetupTime.String,
		LanguageDependency:    gameModel.LanguageDependency.String,
	}

	if categoryModel != nil {
		resp.CategoryInfo = &game.CategoryInfo{
			Category:         categoryModel.Category.String,
			Mode:             categoryModel.Mode.String,
			Theme:            categoryModel.Theme.String,
			Mechanic:         categoryModel.Mechanic.String,
			Portability:      categoryModel.Portability.String,
			TableRequirement: categoryModel.TableRequirement.String,
			SuitableAge:      int32(categoryModel.SuitableAge.Int64),
		}
	}

	return resp, nil
}
