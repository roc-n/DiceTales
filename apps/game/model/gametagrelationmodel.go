package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ GameTagRelationModel = (*customGameTagRelationModel)(nil)

type (
	// GameTagRelationModel is an interface to be customized, add more methods here,
	// and implement the added methods in customGameTagRelationModel.
	GameTagRelationModel interface {
		gameTagRelationModel
		FindTagIdsByGameId(ctx context.Context, gameId uint64) ([]uint64, error)
	}

	customGameTagRelationModel struct {
		*defaultGameTagRelationModel
	}
)

// NewGameTagRelationModel returns a model for the database table.
func NewGameTagRelationModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) GameTagRelationModel {
	return &customGameTagRelationModel{
		defaultGameTagRelationModel: newGameTagRelationModel(conn, c, opts...),
	}
}

func (m *customGameTagRelationModel) FindTagIdsByGameId(ctx context.Context, gameId uint64) ([]uint64, error) {
	key := fmt.Sprintf("cache:gameTagRelation:gameId:%v", gameId)
	var resp []uint64
	err := m.QueryRowCtx(ctx, &resp, key, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("select `tag_id` from %s where `game_id` = ?", m.table)
		return conn.QueryRowsCtx(ctx, v, query, gameId)
	})
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
