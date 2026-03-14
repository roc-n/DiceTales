package logic

import (
	"context"
	"time"

	"dicetales.com/apps/im/model"
	"dicetales.com/apps/im/ws/core"
	"dicetales.com/apps/im/ws/internal/svc"
	"dicetales.com/pkg/id"
	"golang.org/x/net/websocket"
)

type Conversation struct {
	ctx context.Context
	srv *websocket.Server
	svc *svc.ServiceContext
}

func NewConversation(ctx context.Context, srv *websocket.Server, svc *svc.ServiceContext) *Conversation {
	return &Conversation{
		ctx: ctx,
		srv: srv,
		svc: svc,
	}
}

func (l *Conversation) SingleChat(data *core.ChatMessage, userId string) error {
	if data.ConversationId == "" {
		data.ConversationId = id.CombineId(userId, data.RecvId)
	}

	chatLog := model.ChatLog{
		ConversationId: data.ConversationId,
		SendId:         userId,
		RecvId:         data.RecvId,
		ChatType:       data.ChatType,
		MsgType:        data.Wrapper.MType,
		MsgContent:     data.Wrapper.Content,
		SendTime:       time.Now().UnixNano(),
	}
	err := l.svc.ChatLogModel.Insert(l.ctx, &chatLog)

	return err
}
