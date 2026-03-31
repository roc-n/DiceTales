package logic

import (
	"context"

	"dicetales.com/apps/im/rpc/im"
	"dicetales.com/apps/im/rpc/internal/svc"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetChatLogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetChatLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetChatLogLogic {
	return &GetChatLogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetChatLogLogic) GetChatLog(in *im.GetChatLogReq) (*im.GetChatLogResp, error) {
	// 基于时间游标（SendTime）增量拉取消息内容。
	limit := in.Limit
	if limit <= 0 || limit > 100 {
		limit = 100
	}

	logs, err := l.svcCtx.ChatLogModel.FindBySendTime(l.ctx, in.ConversationId, in.CursorTime, limit)
	if err != nil {
		l.Errorf("find chat logs by time cursor err:%v", err)
		return nil, errors.Wrapf(err, "拉取消息列表失败")
	}

	var res []*im.ChatLog

	for _, log := range logs {
		res = append(res, &im.ChatLog{
			MsgId:          log.ID,
			ConversationId: log.ConversationId,
			SendId:         log.SendId,
			RecvId:         log.RecvId,
			MsgType:        log.MsgType,
			MsgContent:     log.MsgContent,
			SendTime:       log.SendTime,
		})
	}

	return &im.GetChatLogResp{
		ChatLogList: res,
	}, nil
}
