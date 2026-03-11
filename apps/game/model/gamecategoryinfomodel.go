package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ GameCategoryInfoModel = (*customGameCategoryInfoModel)(nil)

type (
	// GameCategoryInfoModel is an interface to be customized, add more methods here,
	// and implement the added methods in customGameCategoryInfoModel.
	GameCategoryInfoModel interface {
		gameCategoryInfoModel
	}

	customGameCategoryInfoModel struct {
		*defaultGameCategoryInfoModel
	}
)

// NewGameCategoryInfoModel returns a model for the database table.
func NewGameCategoryInfoModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) GameCategoryInfoModel {
	return &customGameCategoryInfoModel{
		defaultGameCategoryInfoModel: newGameCategoryInfoModel(conn, c, opts...),
	}
}
