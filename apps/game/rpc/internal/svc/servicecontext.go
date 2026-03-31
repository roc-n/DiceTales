package svc

import (
	"dicetales.com/apps/game/model"
	"dicetales.com/apps/game/rpc/internal/config"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config
	*redis.Redis
	GameModel            model.GameModel
	GameResourceModel    model.GameResourceModel
	TagModel             model.TagModel
	GameTagRelationModel model.GameTagRelationModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewMysql(c.DataSource)
	return &ServiceContext{
		Config:               c,
		Redis:                redis.MustNewRedis(c.Redisx),
		GameModel:            model.NewGameModel(sqlConn, c.Cache),
		GameResourceModel:    model.NewGameResourceModel(sqlConn, c.Cache),
		TagModel:             model.NewTagModel(sqlConn, c.Cache),
		GameTagRelationModel: model.NewGameTagRelationModel(sqlConn, c.Cache),
	}
}
