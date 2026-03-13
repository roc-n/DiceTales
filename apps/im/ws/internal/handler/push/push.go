package push

import (
	"dicetales.com/apps/im/ws/aid"
	"dicetales.com/apps/im/ws/internal/svc"
	"dicetales.com/apps/im/ws/websocket"
	"dicetales.com/pkg/constants"
	"github.com/mitchellh/mapstructure"
)

func Push(svc *svc.ServiceContext) websocket.HandlerFunc {
	return func(srv *websocket.Server, conn *websocket.Connx, msg *websocket.Message) {
		var data aid.PushMessage
		if err := mapstructure.Decode(msg.Data, &data); err != nil {
			srv.Send(websocket.NewErrMessage(err))
			return
		}

		switch data.ChatType {
		case constants.SingleChatType:
			single(srv, &data, data.RecvId)
		case constants.GroupChatType:
			group(srv, &data)
		}
	}
}

func single(srv *websocket.Server, data *aid.PushMessage, recvId string) error {
	// 获取目标连接
	recvConn := srv.GetConn(recvId)
	if recvConn == nil {
		return nil
	}

	srv.Infof("push msg %v", data)

	return srv.Send(websocket.NewMessage(data.SendId, &aid.ChatMessage{
		ConversationId: data.ConversationId,
		ChatType:       data.ChatType,
		SendTime:       data.SendTime,
		Message: aid.Message{
			MType:   data.MType,
			Content: data.Content,
		},
	}), recvConn)
}

func group(srv *websocket.Server, data *aid.PushMessage) error {
	// 闭包捕获变量，避免并发数据竞争
	for _, id := range data.RecvIds {
		func(id string) {
			srv.Schedule(func() {
				single(srv, data, id)
			})
		}(id)
	}

	return nil
}
