package user

import (
	"dicetales.com/apps/im/ws/internal/svc"
	"dicetales.com/apps/im/ws/websocket"
)

func Online(svc *svc.ServiceContext) websocket.HandlerFunc {
	return func(srv *websocket.Server, cx *websocket.Connx, msg *websocket.Message) {
		uids := srv.GetUsers()
		u := srv.GetUsers(cx)
		err := srv.Send(websocket.NewMessage(u[0], uids), cx)
		srv.Info("err ", err)
	}
}
