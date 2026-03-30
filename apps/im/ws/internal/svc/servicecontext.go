package svc

import (
	"dicetales.com/apps/im/model"
	"dicetales.com/apps/im/mq/client"
	"dicetales.com/apps/im/ws/internal/config"

	"github.com/bwmarrin/snowflake"
)

type ServiceContext struct {
	Config config.Config

	model.ChatLogModel
	client.TransferClient
	Snowflake *snowflake.Node
}

func NewServiceContext(c config.Config) *ServiceContext {
	node, _ := snowflake.NewNode(c.Snowflake.Node)

	return &ServiceContext{
		Config:         c,
		TransferClient: client.NewChatTransferClient(c.MessageTransfer.Addrs, c.MessageTransfer.Topic),
		ChatLogModel:   model.MustChatLogModel(c.Mongo.Url, c.Mongo.Db, "chatlog"),
		Snowflake:      node,
	}
}
