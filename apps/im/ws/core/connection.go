package core

import (
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

type Connx struct {
	*websocket.Conn
	mu sync.Mutex

	// 为调用Server相关属性和方法，直接持有Server实例指针
	s *Server

	Uid string

	// 心跳检测相关
	idleTime    time.Time     // 开始空闲时刻
	idleTimeOut time.Duration // 最大空闲时间

	// ACK机制
	ackM *AckManager

	// 布尔型活跃标记位，1表示活跃，0表示不活跃，专门用于时间轮策略
	isAlive int32

	// 通知handlerWrite
	msgChan chan *Message
	// 退出信号
	done chan struct{}
}

func NewConnx(s *Server, w http.ResponseWriter, r *http.Request) *Connx {
	c, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.Errorf("conn upgrade err %v", err)
		return nil
	}

	conn := &Connx{
		Conn: c,
		s:    s,

		idleTime:    time.Now(),
		idleTimeOut: s.opt.maxConnectionIdle,

		ackM: &AckManager{
			pendingMsg: make([]*Message, 0, 2),
		},

		isAlive: 1, // 初始默认活跃

		msgChan: make(chan *Message, 1),
		done:    make(chan struct{}),
	}

	// 根据配置的心跳策略选择不同的处理方式
	if s.opt.heartbeatStrategy == HeartbeatStrategyTimingWheel && s.tw != nil {
		s.tw.addConn(conn)
	} else {
		// 默认传统策略：每个连接开启单独的心跳检测协程
		go conn.heartBeat()
	}

	return conn
}

func (c *Connx) appendMessage(msg *Message) {
	c.ackM.msgMu.Lock()
	defer c.ackM.msgMu.Unlock()

	if msg.FrameType == FrameAck {
		return
	}

	c.ackM.pendingMsg = append(c.ackM.pendingMsg, msg)
}

func (c *Connx) ReadMessage() (messageType int, data []byte, err error) {
	messageType, data, err = c.Conn.ReadMessage()

	// 时间轮策略下，仅需极度轻量的原子更新
	if c.s.opt.heartbeatStrategy == HeartbeatStrategyTimingWheel {
		atomic.StoreInt32(&c.isAlive, 1)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.idleTime = time.Time{} // 零值标记活跃状态

	return
}

func (c *Connx) WriteMessage(messageType int, data []byte) error {
	// 时间轮策略下，仅需极度轻量的原子更新
	if c.s.opt.heartbeatStrategy == HeartbeatStrategyTimingWheel {
		atomic.StoreInt32(&c.isAlive, 1)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	err := c.Conn.WriteMessage(messageType, data)
	c.idleTime = time.Now()

	return err
}

func (c *Connx) heartBeat() {
	idleTimer := time.NewTimer(c.idleTimeOut)
	defer func() {
		idleTimer.Stop()
	}()

	for {
		select {
		case <-idleTimer.C:
			c.mu.Lock()
			idle := c.idleTime
			if idle.IsZero() {
				// The connection is non-idle.
				c.mu.Unlock()
				idleTimer.Reset(c.idleTimeOut)
				continue
			}
			val := c.idleTimeOut - time.Since(idle)
			c.mu.Unlock()
			if val <= 0 {
				// The connection has been idle for a duration of idleTimeOut or more.
				c.s.Close(c)
				return
			}
			idleTimer.Reset(val)
		case <-c.done:
			return
		}
	}
}

func (c *Connx) Close() error {
	select {
	case <-c.done:
	default:
		close(c.done)
	}

	return c.Conn.Close()
}
