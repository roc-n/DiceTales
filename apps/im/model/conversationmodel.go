package model

import (
"context"
"github.com/zeromicro/go-zero/core/stores/mon"
"go.mongodb.org/mongo-driver/bson"
)

var _ ConversationModel = (*customConversationModel)(nil)

type (
ConversationModel interface {
conversationModel
FindByConversationIds(ctx context.Context, conversationIds []string) ([]*Conversation, error)
FindOneByConversationId(ctx context.Context, conversationId string) (*Conversation, error)
UpdateMsg(ctx context.Context, conversationId string, maxSeq int64, latestMsg string, lastMsgTime int64) error
}

customConversationModel struct {
*defaultConversationModel
}
)

func NewConversationModel(url, db, collection string) ConversationModel {
conn := mon.MustNewModel(url, db, collection)
return &customConversationModel{
defaultConversationModel: newDefaultConversationModel(conn),
}
}

func MustConversationModel(url, db string) ConversationModel {
return NewConversationModel(url, db, "conversation")
}

func (m *customConversationModel) FindByConversationIds(ctx context.Context, conversationIds []string) ([]*Conversation, error) {
var data []*Conversation
err := m.conn.Find(ctx, &data, bson.M{"conversationId": bson.M{"$in": conversationIds}})
return data, err
}

func (m *customConversationModel) FindOneByConversationId(ctx context.Context, conversationId string) (*Conversation, error) {
var data Conversation
err := m.conn.FindOne(ctx, &data, bson.M{"conversationId": conversationId})
return &data, err
}

func (m *customConversationModel) UpdateMsg(ctx context.Context, conversationId string, maxSeq int64, latestMsg string, lastMsgTime int64) error {
_, err := m.conn.UpdateOne(ctx,
bson.M{"conversationId": conversationId},
bson.M{"$set": bson.M{"maxSeq": maxSeq, "latestMsg": latestMsg, "lastMsgTime": lastMsgTime}},
)
return err
}
