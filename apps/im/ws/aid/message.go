package aid

import (
	"dicetales.com/pkg/constants"
)

type (
	//聊天信息的基础结构体
	Message struct {
		constants.MType `mapstructure:"mType"`
		Content         string `mapstructure:"content"`
	}

	ChatMessage struct {
		ConversationId string `mapstructure:"conversationId"`
		SendId         string `mapstructure:"sendId"`
		RecvId         string `mapstructure:"recvId"`
		SendTime       int64  `mapstructure:"sendTime"`

		constants.ChatType `mapstructure:"chatType"`
		Message            `mapstructure:"message"`
	}

	PushMessage struct {
		ConversationId string `mapstructure:"conversationId"`
		SendId         string `mapstructure:"sendId"`
		RecvId         string `mapstructure:"recvId"`
		SendTime       int64  `mapstructure:"sendTime"`

		// TODO Content&MType可合并为Message
		constants.ChatType `mapstructure:"chatType"`
		Content            string `mapstructure:"content"`
		constants.MType    `mapstructure:"mType"`

		RecvIds []string `mapstructure:"recvIds"`
	}
)
