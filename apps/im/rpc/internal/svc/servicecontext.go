package svc

import (
	"dicetales.com/apps/im/model"
	"dicetales.com/apps/im/rpc/internal/config"
	"github.com/bwmarrin/snowflake"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type ServiceContext struct {
	Config            config.Config
	ConversationModel model.ConversationModel
	ChatLogModel      model.ChatLogModel
	Redis             *redis.Redis
	Snowflake         *snowflake.Node
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化 Snowflake 节点
	node, err := snowflake.NewNode(c.Snowflake.Node)
	if err != nil {
		panic("snowflake generate start err: " + err.Error())
	}

	return &ServiceContext{
		Config:            c,
		ConversationModel: model.NewConversationModel(c.Mongo.Url, c.Mongo.Db, "conversation"),
		ChatLogModel:      model.NewChatLogModel(c.Mongo.Url, c.Mongo.Db, "chat_log"),
		Redis:             redis.MustNewRedis(c.Redis),
		Snowflake:         node,
	}
}
