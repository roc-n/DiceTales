package websocket

type Route struct {
	Method  string
	Handler HandlerFunc
}

type HandlerFunc func(srv *Server, cx *Connx, msg *Message)
