package model

import (
	"dicetales.com/apps/im/ws/core"
)

type ChatLog struct {
	ID int64 `bson:"_id,omitempty" json:"id,omitempty"`

	// TODO: Fill your own fields
	ConversationId string `bson:"conversationId"`
	SendId         string `bson:"sendId"`
	RecvId         string `bson:"recvId"`

	ChatType   core.ChatType `bson:"chatType"`
	MsgType    core.MType    `bson:"msgType"`
	MsgContent string        `bson:"msgContent"`
	SendTime   int64         `bson:"sendTime"`
	Status     int           `bson:"status"`
}
