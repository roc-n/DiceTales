package logic

import (
	"context"
	"fmt"

	"dicetales.com/apps/im/rpc/im"
	"dicetales.com/apps/im/rpc/internal/svc"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type AckMessageReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAckMessageReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AckMessageReadLogic {
	return &AckMessageReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AckMessageReadLogic) AckMessageRead(in *im.AckMessageReadReq) (*im.AckMessageReadResp, error) {
	// 精简后的架构：全量依赖客户端计算与 Redis 轻量存储游标，不再去触碰庞大笨重的 MongoDB
	userReadSeqKey := fmt.Sprintf("user_read_seq:%s", in.OwnerId)
	err := l.svcCtx.Redis.HsetCtx(l.ctx, userReadSeqKey, in.ConversationId, fmt.Sprintf("%d", in.ReadSeq))

	if err != nil {
		l.Errorf("redis update read seq err:%v", err)
		return nil, errors.Wrapf(err, "更新消息已读缓存失败")
	}

	return &im.AckMessageReadResp{
		Success: true,
	}, nil
}
