package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ GameModel = (*customGameModel)(nil)

type (
	// GameModel is an interface to be customized, add more methods here,
	// and implement the added methods in customGameModel.
	GameModel interface {
		gameModel
	}

	customGameModel struct {
		*defaultGameModel
	}
)

// NewGameModel returns a model for the database table.
func NewGameModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) GameModel {
	return &customGameModel{
		defaultGameModel: newGameModel(conn, c, opts...),
	}
}
