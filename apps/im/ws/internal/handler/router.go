package handler

import (
	"dicetales.com/apps/im/ws/core"
	"dicetales.com/apps/im/ws/internal/svc"
)

func RegisterHandlers(srv *core.Server, svc *svc.ServiceContext) {
	srv.AddRoutes([]core.Route{
		{
			Method:  "user.online",
			Handler: Online(svc),
		},
		{
			Method:  "conversation.chat",
			Handler: Chat(svc),
		},
		{
			Method:  "push",
			Handler: Push(svc),
		},
	})
}
