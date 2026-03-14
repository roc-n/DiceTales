package core

type FrameType uint8

const (
	FrameData  FrameType = 0x0
	FramePing  FrameType = 0x1
	FrameAck   FrameType = 0x2
	FrameNoAck FrameType = 0x3
	FrameErr   FrameType = 0x9
)

type Message struct {
	FrameType `json:"frameType"`
	Id        int `json:"id"`

	Method string `json:"method"`
	Data   any    `json:"data"` // map[string]any
}

func NewMessage(data any) *Message {
	return &Message{
		FrameType: FrameData,
		Data:      data,
	}
}

func NewErrMessage(err error) *Message {
	return &Message{
		FrameType: FrameErr,
		Data:      err.Error(),
	}
}
