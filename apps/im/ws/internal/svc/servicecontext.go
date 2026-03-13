package svc

import (
	// "dicetales.com/apps/im/immodels"
	"dicetales.com/apps/im/ws/internal/config"
	// "dicetales.com/apps/task/mq/mqclient"
)

type ServiceContext struct {
	Config config.Config

	// immodels.ChatLogModel
	// mqclient.MessageTransferClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		// MessageTransferClient: mqclient.NewMsgChatTransferClient(c.MessageTransfer.Addrs, c.MessageTransfer.Topic),
		// ChatLogModel:          immodels.MustChatLogModel(c.Mongo.Url, c.Mongo.Db),
	}
}
