package logic

import (
"context"
"fmt"
"time"

"dicetales.com/apps/im/model"
"dicetales.com/apps/im/rpc/im"
"dicetales.com/apps/im/rpc/internal/svc"

"github.com/pkg/errors"
"github.com/zeromicro/go-zero/core/logx"
)

type PutChatLogLogic struct {
ctx    context.Context
svcCtx *svc.ServiceContext
logx.Logger
}

func NewPutChatLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PutChatLogLogic {
return &PutChatLogLogic{
ctx:    ctx,
svcCtx: svcCtx,
Logger: logx.WithContext(ctx),
}
}

func (l *PutChatLogLogic) PutChatLog(in *im.PutChatLogReq) (*im.PutChatLogResp, error) {
msgId := in.MsgId
if msgId == "" {
msgId = l.svcCtx.Snowflake.Generate().String()
}
nowTime := time.Now().UnixMilli()

// 获取全局唯一递增Seq，用于维护会话的最大可用Seq（供未读数计算使用）
seqKey := fmt.Sprintf("conv_seq:%s", in.ConversationId)
seq, err := l.svcCtx.Redis.IncrCtx(l.ctx, seqKey)
if err != nil {
l.Errorf("redis incr seq err: %v", err)
return nil, errors.Wrapf(err, "获取消息递增序列失败")
}

chatLog := &model.ChatLog{
ID:             msgId,
ConversationId: in.ConversationId,
SendId:         in.SendId,
RecvId:         in.RecvId,
ChatType:       in.ChatType,
MsgType:        in.MsgType,
MsgContent:     in.MsgContent,
SendTime:       nowTime,
}

err = l.svcCtx.ChatLogModel.Insert(l.ctx, chatLog)
if err != nil {
l.Errorf("insert chatlog err: %v", err)
return nil, errors.Wrapf(err, "持久化消息记录失败")
}

// 更新这张全局记录表上的最新消息和最新 Seq，以便首页可以直接预览拉取
err = l.svcCtx.ConversationModel.UpdateMsg(l.ctx, in.ConversationId, seq, in.MsgContent, nowTime)
if err != nil {
l.Errorf("update conversation msg err: %v", err)
}

// 更新发送端的 Zset (时间戳为 score) 维护列表排序
senderZsetKey := fmt.Sprintf("user_conversations:%s", in.SendId)
_, _ = l.svcCtx.Redis.ZaddCtx(l.ctx, senderZsetKey, nowTime, in.ConversationId)

// 更新接收端的 Zset 维护列表排序
if in.ChatType == 1 {
recvZsetKey := fmt.Sprintf("user_conversations:%s", in.RecvId)
_, _ = l.svcCtx.Redis.ZaddCtx(l.ctx, recvZsetKey, nowTime, in.ConversationId)
} else {
// todo 群聊需要发送到消息队列异步分发给群内的其他所有组员（因为群员可能有几万人）
}

return &im.PutChatLogResp{
MsgId:    msgId,
Seq:      seq,
SendTime: nowTime,
}, nil
}
