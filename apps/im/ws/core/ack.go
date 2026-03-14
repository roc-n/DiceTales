package core

import (
	"sync"
	"time"
)

type AckType int

const (
	NoAck AckType = iota
	OnceAck
)

type AckManager struct {
	msgMu      sync.Mutex
	pendingMsg []*Message
	cursor     int
}

func (m *AckManager) Process(conn *Connx) {
	for {
		select {
		case <-conn.done:
			return
		default:
		}

		m.msgMu.Lock()
		if len(m.pendingMsg) == 0 {
			m.msgMu.Unlock()
			time.Sleep(100 * time.Microsecond)
			continue
		}

		// 逐条读取新消息
		message := m.pendingMsg[0]

		switch conn.s.opt.ack {
		case OnceAck:
			conn.s.Send(&Message{
				FrameType: FrameAck,
				Id:        message.Id,
			}, conn)

			m.pendingMsg = m.pendingMsg[1:]
			if message.Id > m.cursor {
				m.cursor++
			}
			m.msgMu.Unlock()
			conn.msgChan <- message
		}
	}
}

func (m *AckManager) IsAck(message *Message, ack AckType) bool {
	if message == nil {
		return ack != NoAck
	}
	return ack != NoAck && message.FrameType != FrameNoAck
}

func (t AckType) AckString() string {
	switch t {
	case OnceAck:
		return "OnceAck"
	}
	return "NoAck"
}
