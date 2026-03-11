package svc

import (
	"dicetales.com/apps/game/model"
	"dicetales.com/apps/game/rpc/internal/config"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config                config.Config
	GameModel             model.GameModel
	GameCategoryInfoModel model.GameCategoryInfoModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewMysql(c.DataSource)
	return &ServiceContext{
		Config:                c,
		GameModel:             model.NewGameModel(sqlConn, c.Cache),
		GameCategoryInfoModel: model.NewGameCategoryInfoModel(sqlConn, c.Cache),
	}
}
