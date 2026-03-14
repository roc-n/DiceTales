package core

import (
	"context"
	"encoding/json"
	"fmt"

	"net/http"

	"sync"

	"dicetales.com/pkg/constants"
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/limit"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/threading"
)

type Server struct {
	sync.RWMutex
	logx.Logger

	opt            *serverOption
	authentication Authentication

	// 路由表
	patten string
	routes map[string]HandlerFunc
	addr   string

	// websocket连接管理
	upgrader   websocket.Upgrader
	connToUser map[*Connx]string
	userToConn map[string]*Connx

	// 群聊业务相关
	*threading.TaskRunner

	// 用户级消息限流器
	limiter sync.Map
}

/*
WebSocket服务初始化相关，包括配置加载、服务启动停止等
*/

func NewServer(addr string, opts ...ServerOptions) *Server {
	opt := newServerOptions(opts...)

	return &Server{
		routes:   make(map[string]HandlerFunc),
		addr:     addr,
		patten:   opt.patten,
		opt:      &opt,
		upgrader: websocket.Upgrader{},

		authentication: opt.Authentication,

		connToUser: make(map[*Connx]string),
		userToConn: make(map[string]*Connx),

		Logger: logx.WithContext(context.Background()),

		TaskRunner: threading.NewTaskRunner(opt.groupMsgConcurrency),
	}
}

func (s *Server) AddRoutes(rs []Route) {
	for _, r := range rs {
		s.routes[r.Method] = r.Handler
	}
}

func (s *Server) Start() {
	http.HandleFunc(s.patten, s.ServeWs)
	s.Info(http.ListenAndServe(s.addr, nil))
}

func (s *Server) Stop() {
	fmt.Println("停止服务")
}

func (s *Server) ServeWs(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			s.Errorf("server handler ws recover err %v", r)
		}
	}()

	cx := NewConnx(s, w, r)
	if cx == nil {
		return
	}

	// 根据连接对象获取请求
	if !s.authentication.Auth(w, r) {
		s.Send(&Message{FrameType: FrameData, Data: "不具备访问权限"}, cx)
		cx.Close()
		return
	}

	// 记录连接
	s.addConn(cx, r)

	// 处理连接
	go s.handlerConn(cx)
}

/*
连接管理相关，包括连接的建立、关闭，以及连接与用户的映射关系维护
*/
func (s *Server) addConn(cx *Connx, req *http.Request) {
	uid := s.authentication.UserId(req)

	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	// 验证用户是否已在线
	if c := s.userToConn[uid]; c != nil {
		c.Close()
	}

	cx.Uid = uid
	s.connToUser[cx] = uid
	s.userToConn[uid] = cx
}

func (s *Server) GetConn(uid string) *Connx {
	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()

	return s.userToConn[uid]
}

func (s *Server) GetConns(uids ...string) []*Connx {
	if len(uids) == 0 {
		return nil
	}

	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()

	res := make([]*Connx, 0, len(uids))
	for _, uid := range uids {
		res = append(res, s.userToConn[uid])
	}
	return res
}

func (s *Server) GetUsers(cxs ...*Connx) []string {
	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()

	var res []string
	if len(cxs) == 0 {
		// 获取全部
		res = make([]string, 0, len(s.connToUser))
		for _, uid := range s.connToUser {
			res = append(res, uid)
		}
	} else {
		// 获取部分
		res = make([]string, 0, len(cxs))
		for _, conn := range cxs {
			res = append(res, s.connToUser[conn])
		}
	}

	return res
}

func (s *Server) Close(cx *Connx) {
	cx.Close()

	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	uid := s.connToUser[cx]
	if uid == "" {
		return
	}

	delete(s.connToUser, cx)
	delete(s.userToConn, uid)

	// 删除用户限流器
	s.limiter.Delete(uid)
}

func (s *Server) SendByUserId(msg any, sendIds ...string) error {
	if len(sendIds) == 0 {
		return nil
	}

	return s.Send(msg, s.GetConns(sendIds...)...)
}

func (s *Server) Send(msg any, cxs ...*Connx) error {
	if len(cxs) == 0 {
		return nil
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	for _, conn := range cxs {
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			return err
		}
	}

	return nil
}

// 基于Conn处理消息
func (s *Server) handlerConn(cx *Connx) {

	// 处理任务
	go s.handlerWrite(cx)

	if cx.ackM.IsAck(nil, s.opt.ack) {
		go cx.ackM.Process(cx)
	}

	for {
		// 获取请求消息，暂时仅处理文本消息
		_, msg, err := cx.ReadMessage()
		if err != nil {
			s.Errorf("websocket conn read message err %v", err)
			s.Close(cx)
			return
		}
		// 解析消息
		var message Message
		if err = json.Unmarshal(msg, &message); err != nil {
			s.Errorf("json unmarshal err %v, msg %v", err, string(msg))
			s.Close(cx)
			return
		}

		// 消息限流，当前仅针对发送消息限流
		// if message.Method == "conversation.chat" && !s.getLimiter(cx.uid).Allow() {
		// 	s.Send(&Message{
		// 		FrameType: FrameErr,
		// 		Data:      "消息发送过快，请稍后再试",
		// 	}, cx)
		// 	continue // 拦截消息，不再做进一步处理
		// }

		// // 敏感词过滤
		// if s.opt.sensitiveFilter != nil {
		// 	text, _ := message.Data.(string)
		// 	if ok, word := s.opt.sensitiveFilter.ContainsSensitive(text); ok {
		// 		s.Send(&Message{
		// 			FrameType: FrameErr,
		// 			Data:      fmt.Sprintf("消息包含敏感词：%s", word),
		// 		}, cx)
		// 		continue // 拦截消息，不再做进一步处理
		// 	}
		// }

		if cx.ackM.IsAck(&message, s.opt.ack) {
			cx.appendMessage(&message)
		} else {
			cx.msgChan <- &message
		}
	}
}

// 任务处理
func (s *Server) handlerWrite(c *Connx) {
	for {
		select {
		case <-c.done:
			return
		case message := <-c.msgChan:
			switch message.FrameType {
			case FramePing:
				s.Send(&Message{FrameType: FramePing}, c)
			case FrameData:
				// 根据请求的method路由到具体的handler
				if handler, ok := s.routes[message.Method]; ok {
					handler(s, c, message)
				} else {
					s.Send(&Message{FrameType: FrameData, Data: fmt.Sprintf("不存在执行的方法 %v 请检查", message.Method)}, c)
				}
			}
		}
	}
}

// 获取用户消息限流器
func (s *Server) getLimiter(uid string) *limit.TokenLimiter {
	if uid == constants.SYSTEM_ROOT_UID {
		// 系统用户不做限流
		return nil
	}
	limiter, ok := s.limiter.Load(uid)
	if !ok {
		limiter = limit.NewTokenLimiter(5, 20, nil, "")
		s.limiter.Store(uid, limiter)
	}
	return limiter.(*limit.TokenLimiter)
}
