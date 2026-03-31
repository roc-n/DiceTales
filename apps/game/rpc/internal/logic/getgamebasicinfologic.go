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

// 批量拉取桌游轻量摘要
func (l *GetGameBasicInfoLogic) GetGameBasicInfo(in *game.GetGameBasicInfoReq) (*game.GetGameBasicInfoResp, error) {
	resp := &game.GetGameBasicInfoResp{
		GameBasics: make(map[int64]*game.GameBasic),
	}

	// 优化点：虽然可以写一个 FindByIds 批量查，但由于我们接入了缓存层
	// 简单的遍历 FindOne 利用高速缓存也是比较常见的做法
	// 如果性能要求极高，可以在模型层增加带 Redis MGET 的批量查询接口
	for _, id := range in.GameIds {
		info, err := l.svcCtx.GameModel.FindOne(l.ctx, uint64(id))
		if err == nil && info != nil {
			resp.GameBasics[id] = &game.GameBasic{
				Id:       int64(info.Id),
				Name:     info.Name,
				CoverImg: info.CoverImg,
				Score:    info.Score,
			}
		}
	}

	return resp, nil
}
