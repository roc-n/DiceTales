package model

import (
	"time"

	"dicetales.com/apps/im/ws/core"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Conversation struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`

	ConversationId string        `bson:"conversationId,omitempty"`
	ChatType       core.ChatType `bson:"chatType,omitempty"`
	IsShow         bool          `bson:"isShow,omitempty"`
	Total          int           `bson:"total,omitempty"`
	Seq            int64         `bson:"seq"`
	Msg            *ChatLog      `bson:"msg,omitempty"`

	// TODO: Fill your own fields
	UpdateAt time.Time `bson:"updateAt,omitempty" json:"updateAt,omitempty"`
	CreateAt time.Time `bson:"createAt,omitempty" json:"createAt,omitempty"`
}
