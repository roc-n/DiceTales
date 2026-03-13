package conversation

import (
	"dicetales.com/apps/im/ws/aid"
	"dicetales.com/apps/im/ws/internal/svc"
	"dicetales.com/apps/im/ws/websocket"

	// mqAid "dicetales.com/apps/task/mq/aid"
	"dicetales.com/pkg/constants"
	// "dicetales.com/pkg/wuid"
	"github.com/mitchellh/mapstructure"
)

func Chat(svc *svc.ServiceContext) websocket.HandlerFunc {
	return func(srv *websocket.Server, cx *websocket.Connx, msg *websocket.Message) {
		// 接收websocket层的消息
		var data aid.ChatMessage
		if err := mapstructure.Decode(msg.Data, &data); err != nil {
			srv.Send(websocket.NewErrMessage(err), cx)
			return
		}

		if data.ConversationId == "" {
			switch data.ChatType {
			case constants.SingleChatType:
				// data.ConversationId = wuid.CombineId(cx.Uid, data.RecvId)
			case constants.GroupChatType:
				data.ConversationId = data.RecvId

			}
		}

		// 转化为mq消息格式并推送
		// err := svc.MessageTransferClient.Push(context.Background(), &mqAid.MessageTransfer{
		// 	ConversationId: data.ConversationId,
		// 	ChatType:       data.ChatType,
		// 	SendId:         cx.Uid,
		// 	RecvId:         data.RecvId,
		// 	SendTime:       time.Now().UnixMilli(),
		// 	MType:          data.Message.MType,
		// 	Content:        data.Message.Content,
		// })
		// if err != nil {
		// 	srv.Send(websocket.NewErrMessage(err), cx)
		// 	return
		// }
	}
}
