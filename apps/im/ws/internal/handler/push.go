package handler

import (
	"dicetales.com/apps/im/ws/core"
	"dicetales.com/apps/im/ws/internal/svc"
	"github.com/mitchellh/mapstructure"
)

func Push(svc *svc.ServiceContext) core.HandlerFunc {
	return func(srv *core.Server, conn *core.Connx, msg *core.Message) {
		var data core.ChatMessage
		if err := mapstructure.Decode(msg.Data, &data); err != nil {
			srv.Send(core.NewErrMessage(err), conn)
			return
		}

		switch data.ChatType {
		case core.SingleChatType:
			single(srv, &data, data.RecvId)
		case core.GroupChatType:
			group(srv, &data)
		}
	}
}

func single(srv *core.Server, data *core.ChatMessage, recvId string) error {
	// 获取目标连接
	recvConn := srv.GetConn(recvId)
	if recvConn == nil {
		return nil
	}

	srv.Infof("push msg %v", data)

	return srv.Send(core.NewMessage(&core.ChatMessage{
		ConversationId: data.ConversationId,
		ChatType:       data.ChatType,
		SendId:         data.SendId,
		RecvId:         recvId,
		SendTime:       data.SendTime,
		Wrapper:        data.Wrapper,
	}), recvConn)
}

func group(srv *core.Server, data *core.ChatMessage) error {
	for _, id := range data.RecvIds {
		func(id string) {
			srv.Schedule(func() {
				single(srv, data, id)
			})
		}(id)
	}

	return nil
}
