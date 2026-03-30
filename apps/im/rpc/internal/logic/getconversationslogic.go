package logic

import (
"context"
"fmt"

"dicetales.com/apps/im/rpc/im"
"dicetales.com/apps/im/rpc/internal/svc"

"github.com/zeromicro/go-zero/core/logx"
)

type GetConversationsLogic struct {
ctx    context.Context
svcCtx *svc.ServiceContext
logx.Logger
}

func NewGetConversationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetConversationsLogic {
return &GetConversationsLogic{
ctx:    ctx,
svcCtx: svcCtx,
Logger: logx.WithContext(ctx),
}
}

func (l *GetConversationsLogic) GetConversations(in *im.GetConversationsReq) (*im.GetConversationsResp, error) {
// 1. 无需从 DB 查询，直接根据用户的 Zset 拉取他参与过的最活跃的所有会话 ID (0, -1 为倒序拉取所有)
userZsetKey := fmt.Sprintf("user_conversations:%s", in.OwnerId)
convIds, err := l.svcCtx.Redis.ZrevrangeCtx(l.ctx, userZsetKey, 0, -1)
if err != nil {
l.Errorf("redis zrevrange err:%v", err)
return nil, err
}

if len(convIds) == 0 {
return &im.GetConversationsResp{}, nil
}

// 2. 根据获取到的 conversationIds 批量前往 MongoDB 拉取这批会话的全局事实记录
conversations, err := l.svcCtx.ConversationModel.FindByConversationIds(l.ctx, convIds)
if err != nil {
l.Errorf("find conversations by ids err:%v", err)
return nil, err
}

// 把拿到的 conversations 利用 map 排序或直接装配返回
var res []*im.Conversation

// redis拉出来是有序的，但 mongodb in 查询并不是按指定顺序返回，所以需要以 convId 为基准排装配一下
convMap := make(map[string]*im.Conversation)
for _, conv := range conversations {
convMap[conv.ConversationId] = &im.Conversation{
ConversationId: conv.ConversationId,
ChatType:       int32(conv.ChatType),
MaxSeq:         conv.MaxSeq,
LatestMsg:      conv.LatestMsg,
LastMsgTime:    conv.LastMsgTime,
}
}

// 重新按照 Redis Zset 倒序出来的正确的时间顺序放入 res 中
for _, cid := range convIds {
if c, ok := convMap[cid]; ok {
res = append(res, c)
}
}

return &im.GetConversationsResp{
ConversationList: res,
}, nil
}
