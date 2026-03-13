package handler

import (
	"dicetales.com/apps/im/ws/internal/handler/conversation"
	"dicetales.com/apps/im/ws/internal/handler/push"
	"dicetales.com/apps/im/ws/internal/handler/user"
	"dicetales.com/apps/im/ws/internal/svc"
	"dicetales.com/apps/im/ws/websocket"
)

func RegisterHandlers(srv *websocket.Server, svc *svc.ServiceContext) {
	srv.AddRoutes([]websocket.Route{
		{
			Method:  "user.online",
			Handler: user.Online(svc),
		},
		{
			Method:  "conversation.chat",
			Handler: conversation.Chat(svc),
		},
		{
			Method:  "push",
			Handler: push.Push(svc),
		},
	})
}
