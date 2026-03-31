package logic

import (
	"context"

	"dicetales.com/apps/im/model"
	"dicetales.com/apps/im/rpc/im"
	"dicetales.com/apps/im/rpc/internal/svc"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/zeromicro/go-zero/core/logx"
)

type SetUpUserConversationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSetUpUserConversationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetUpUserConversationLogic {
	return &SetUpUserConversationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// SetUpUserConversation 建立或唤醒会话全局记录（无冗余复制版）
func (l *SetUpUserConversationLogic) SetUpUserConversation(in *im.SetUpUserConversationReq) (*im.SetUpUserConversationResp, error) {
	// 1. 根据类型生成唯一的 ConversationId
	var conversationId string
	switch in.ChatType {
	case 1: // 单聊: 按照字典序拼接防止 u1_u2 和 u2_u1 生成不同的会话
		if in.SendId < in.RecvId {
			conversationId = in.SendId + "_" + in.RecvId
		} else {
			conversationId = in.RecvId + "_" + in.SendId
		}
	case 2: // 群聊: 直接使用接收方ID（即群组的ID）作为会话ID
		conversationId = in.RecvId
	default:
		return nil, errors.New("不支持的聊天类型")
	}

	// 2. 检查全局是否已存在这一唯一的会话记录
	conv, err := l.svcCtx.ConversationModel.FindOneByConversationId(l.ctx, conversationId)

	if err != nil && err != mongo.ErrNoDocuments {
		l.Errorf("查询会话记录失败: %v", err)
		return nil, errors.Wrapf(err, "查询会话记录失败")
	}

	if err == mongo.ErrNoDocuments || conv == nil {
		// 3. 真正极简的表：只记录会话本身事实存在，脱离与个体的绑定，不再保留 Participants
		newConv := &model.Conversation{
			ConversationId: conversationId,
			ChatType:       in.ChatType,
		}

		if err := l.svcCtx.ConversationModel.Insert(l.ctx, newConv); err != nil {
			l.Errorf("创建新会话底层记录失败: %v", err)
			return nil, errors.Wrapf(err, "创建新会话失败")
		}
	}

	return &im.SetUpUserConversationResp{
		ConversationId: conversationId,
	}, nil
}
