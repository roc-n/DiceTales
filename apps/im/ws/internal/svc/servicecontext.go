package svc

import (
	"dicetales.com/apps/im/model"
	"dicetales.com/apps/im/mq/client"
	"dicetales.com/apps/im/ws/internal/config"
)

type ServiceContext struct {
	Config config.Config

	model.ChatLogModel
	client.TransferClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:         c,
		TransferClient: client.NewChatTransferClient(c.MessageTransfer.Addrs, c.MessageTransfer.Topic),
		ChatLogModel:   model.MustChatLogModel(c.Mongo.Url, c.Mongo.Db),
	}
}
