package handler

import (
	"context"
	"time"

	"dicetales.com/apps/im/ws/core"
	"dicetales.com/apps/im/ws/internal/svc"
	"dicetales.com/pkg/id"

	"github.com/mitchellh/mapstructure"
)

func Chat(svc *svc.ServiceContext) core.HandlerFunc {
	return func(srv *core.Server, cx *core.Connx, msg *core.Message) {
		var data core.ChatMessage
		if err := mapstructure.Decode(msg.Data, &data); err != nil {
			srv.Send(core.NewErrMessage(err), cx)
			return
		}

		if data.ConversationId == "" {
			switch data.ChatType {
			case core.SingleChatType:
				data.ConversationId = id.CombineId(cx.Uid, data.RecvId)
			case core.GroupChatType:
				data.ConversationId = data.RecvId

			}
		}

		// 转化为mq消息格式并推送
		err := svc.TransferClient.Push(context.Background(), &core.ChatMessage{
			ConversationId: data.ConversationId,
			ChatType:       data.ChatType,
			SendId:         cx.Uid,
			RecvId:         data.RecvId,
			SendTime:       time.Now().UnixMilli(),
			Wrapper:        data.Wrapper,
		})
		if err != nil {
			srv.Send(core.NewErrMessage(err), cx)
			return
		}
	}
}
