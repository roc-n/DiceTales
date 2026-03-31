package svc

import (
	"dicetales.com/apps/post/model"
	"dicetales.com/apps/post/rpc/internal/config"

	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config                 config.Config
	PostModel              model.PostModel
	CommentModel           model.CommentModel
	SocialInteractionModel model.SocialInteractionModel
	RedisClient            *redis.Redis
	PostEventPusher        *kq.Pusher
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewMysql(c.Mysql.DataSource)
	return &ServiceContext{
		Config:                 c,
		PostModel:              model.NewPostModel(sqlConn, c.CacheRedis),
		CommentModel:           model.NewCommentModel(c.Mongo.Url, c.Mongo.Db, "comment"),
		SocialInteractionModel: model.NewSocialInteractionModel(c.Mongo.Url, c.Mongo.Db, "social_interaction"),
		RedisClient:            redis.MustNewRedis(c.Redisx),
		PostEventPusher:        kq.NewPusher(c.PostEventPusher.Brokers, c.PostEventPusher.Topic),
	}
}
