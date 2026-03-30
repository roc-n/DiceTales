package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Conversation struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	ConversationId string             `bson:"conversationId"`
	ChatType       int32              `bson:"chatType"`
	MaxSeq         int64              `bson:"maxSeq"`
	LatestMsg      string             `bson:"latestMsg"`
	LastMsgTime    int64              `bson:"lastMsgTime"`
	CreateAt       time.Time          `bson:"createAt,omitempty" json:"createAt,omitempty"`
	UpdateAt       time.Time          `bson:"updateAt,omitempty" json:"updateAt,omitempty"`
}
