package model

import (
	"time"
)

type ChatLog struct {
	ID             string    `bson:"_id"`            // MsgId
	ConversationId string    `bson:"conversationId"` // 会话ID
	SendId         string    `bson:"sendId"`         // 发送人
	RecvId         string    `bson:"recvId"`         // 接收人
	ChatType       int32     `bson:"chatType"`       // 1-单聊 2-群聊
	MsgType        int32     `bson:"msgType"`        // 1-文本 2-图片 3-语音
	MsgContent     string    `bson:"msgContent"`     // 消息正文
	SendTime       int64     `bson:"sendTime"`       // 发送时间
	Status         int       `bson:"status"`         // 状态
	CreateAt       time.Time `bson:"createAt,omitempty"`
	UpdateAt       time.Time `bson:"updateAt,omitempty"`
}
