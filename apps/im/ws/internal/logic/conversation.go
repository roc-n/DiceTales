package logic

import (
	"context"

	// "dicetales.com/apps/im/immodels"
	"dicetales.com/apps/im/ws/aid"
	"dicetales.com/apps/im/ws/internal/svc"
	"dicetales.com/apps/im/ws/websocket"
	"github.com/zeromicro/go-zero/core/errorx"
	// "dicetales.com/pkg/wuid"
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

func (l *Conversation) SingleChat(data *aid.ChatMessage, userId string) error {
	if data.ConversationId == "" {
		// data.ConversationId = wuid.CombineId(userId, data.RecvId)
	}

	// time.Sleep(time.Minute)
	// chatLog := immodels.ChatLog{
	// 	ConversationId: data.ConversationId,
	// 	SendId:         userId,
	// 	RecvId:         data.RecvId,
	// 	ChatType:       data.ChatType,
	// 	MsgFrom:        0,
	// 	MsgType:        data.MType,
	// 	MsgContent:     data.Content,
	// 	SendTime:       time.Now().UnixNano(),
	// }
	// err := l.svc.ChatLogModel.Insert(l.ctx, &chatLog)
	err := errorx.Wrapf(context.Background().Err(), "insert chat log error")

	return err
}
