package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"dicetales.com/apps/im/model"
	"dicetales.com/apps/im/mq/internal/svc"
	"dicetales.com/apps/im/ws/core"
	ws "dicetales.com/apps/im/ws/core"
	"dicetales.com/apps/social/rpc/social"
	"github.com/zeromicro/go-zero/core/logx"
)

type MessageTransfer struct {
	svc *svc.ServiceContext
	logx.Logger
}

func NewMessageTransfer(svc *svc.ServiceContext) *MessageTransfer {
	return &MessageTransfer{
		svc:    svc,
		Logger: logx.WithContext(context.Background()),
	}
}

func (t *MessageTransfer) Consume(ctx context.Context, key, value string) error {

	fmt.Println("key : ", key, " value : ", value)

	var data ws.ChatMessage
	if err := json.Unmarshal([]byte(value), &data); err != nil {
		return err
	}

	// 记录数据
	if err := t.addChatLog(ctx, &data); err != nil {
		return err
	}

	switch data.ChatType {
	case ws.SingleChatType:
		return t.single(&data)
	case ws.GroupChatType:
		return t.group(ctx, &data)
	}

	return nil
}

func (t *MessageTransfer) single(data *ws.ChatMessage) error {
	return t.svc.WsClient.Send(ws.Message{
		FrameType: ws.FrameData,
		Method:    "push",
		Data:      data,
	})
}

func (t *MessageTransfer) group(ctx context.Context, data *ws.ChatMessage) error {

	users, err := t.svc.Social.GroupMember(ctx, &social.GroupMemberReq{
		GroupId: data.RecvId,
	})
	if err != nil {
		return err
	}

	data.RecvIds = make([]string, 0, len(users.List))
	for _, member := range users.List {
		// 过滤自己，不需要给自己发送消息
		if member.Uid == data.SendId {
			continue
		}

		data.RecvIds = append(data.RecvIds, member.Uid)
	}

	return t.svc.WsClient.Send(core.Message{
		FrameType: core.FrameData,
		Method:    "push",
		Data:      data,
	})
}

func (t *MessageTransfer) addChatLog(ctx context.Context, data *ws.ChatMessage) error {

	// MongoDB消息
	chatLog := model.ChatLog{
		ConversationId: data.ConversationId,
		SendId:         data.SendId,
		RecvId:         data.RecvId,
		ChatType:       data.ChatType,
		MsgType:        data.Wrapper.MType,
		MsgContent:     data.Wrapper.Content,
		SendTime:       data.SendTime,
	}
	err := t.svc.ChatLogModel.Insert(ctx, &chatLog)

	if err != nil {
		return err
	}

	return t.svc.ConversationModel.UpdateMsg(ctx, &chatLog)
}
