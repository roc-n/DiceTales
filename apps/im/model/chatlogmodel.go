package model

import (
	"context"

	"github.com/zeromicro/go-zero/core/stores/mon"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var _ ChatLogModel = (*customChatLogModel)(nil)

type (
	ChatLogModel interface {
		chatLogModel
		FindBySendTime(ctx context.Context, conversationId string, cursorTime int64, limit int32) ([]*ChatLog, error)
	}

	customChatLogModel struct {
		*defaultChatLogModel
	}
)

func NewChatLogModel(url, db, collection string) ChatLogModel {
	conn := mon.MustNewModel(url, db, collection)
	return &customChatLogModel{
		defaultChatLogModel: newDefaultChatLogModel(conn),
	}
}

func MustChatLogModel(url, db, collection string) ChatLogModel {
	return NewChatLogModel(url, db, collection)
}

func (m *customChatLogModel) FindBySendTime(ctx context.Context, conversationId string, cursorTime int64, limit int32) ([]*ChatLog, error) {
	var resp []*ChatLog
	err := m.conn.Find(ctx, &resp,
		bson.M{"conversationId": conversationId, "sendTime": bson.M{"$gt": int64(cursorTime)}},
		options.Find().SetSort(bson.M{"sendTime": 1}).SetLimit(int64(limit)),
	)
	return resp, err
}
