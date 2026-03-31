package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ GameResourceModel = (*customGameResourceModel)(nil)

type (
	// GameResourceModel is an interface to be customized, add more methods here,
	// and implement the added methods in customGameResourceModel.
	GameResourceModel interface {
		gameResourceModel
		FindByGameId(ctx context.Context, gameId uint64) ([]*GameResource, error)
	}

	customGameResourceModel struct {
		*defaultGameResourceModel
	}
)

// NewGameResourceModel returns a model for the database table.
func NewGameResourceModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) GameResourceModel {
	return &customGameResourceModel{
		defaultGameResourceModel: newGameResourceModel(conn, c, opts...),
	}
}

func (m *customGameResourceModel) FindByGameId(ctx context.Context, gameId uint64) ([]*GameResource, error) {
	key := fmt.Sprintf("cache:gameResource:gameId:%v", gameId)
	var resp []*GameResource
	err := m.QueryRowCtx(ctx, &resp, key, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("select %s from %s where `game_id` = ?", gameResourceRows, m.table)
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
