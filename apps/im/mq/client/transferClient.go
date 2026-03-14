package client

import (
	"context"
	"encoding/json"

	ws "dicetales.com/apps/im/ws/core"
	"github.com/zeromicro/go-queue/kq"
)

type TransferClient interface {
	Push(ctx context.Context, msg *ws.ChatMessage) error
}

type transferClient struct {
	pusher *kq.Pusher
}

func (c *transferClient) Push(ctx context.Context, msg *ws.ChatMessage) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return c.pusher.Push(ctx, string(body))
}

func NewChatTransferClient(addr []string, topic string, opts ...kq.PushOption) *transferClient {
	return &transferClient{
		pusher: kq.NewPusher(addr, topic),
	}
}
