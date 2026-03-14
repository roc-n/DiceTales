package handler

import (
	"dicetales.com/apps/im/ws/core"
	"dicetales.com/apps/im/ws/internal/svc"
)

func Online(svc *svc.ServiceContext) core.HandlerFunc {
	return func(srv *core.Server, cx *core.Connx, msg *core.Message) {
		uids := srv.GetUsers()
		err := srv.Send(core.NewMessage(uids), cx)
		srv.Info("err ", err)
	}
}
