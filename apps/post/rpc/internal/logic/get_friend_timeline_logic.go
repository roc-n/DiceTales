package logic

import (
	"context"
	"fmt"

	"dicetales.com/apps/post/rpc/internal/svc"
	"dicetales.com/apps/post/rpc/post"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFriendTimelineLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFriendTimelineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFriendTimelineLogic {
	return &GetFriendTimelineLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFriendTimelineLogic) GetFriendTimeline(in *post.GetFriendTimelineReq) (*post.GetFriendTimelineResp, error) {
	// 1. ZREVRANGEBYSCORE from Redis Inbox `user:inbox:{uid}` to perform an O(logN) pull
	inboxKey := fmt.Sprintf("user:inbox:%d", in.UserId)

	limit := in.Limit
	if limit == 0 {
		limit = 20
	}

	// Logic outline:
	// redisClient.ZrevrangebyscoreWithLimitCtx(...)
	_ = inboxKey // Use inboxKey later for redis query

	// 2. We extract the array of PostIds
	var postIds []int64
	// e.g., []int64{101, 100, 99}

	if len(postIds) == 0 {
		return &post.GetFriendTimelineResp{
			List: []*post.PostDetail{},
		}, nil
	}

	// 3. Batch assemble posts using MR (MapReduce) conceptually or in parallel locally,
	// querying MySQL/Redis caches, attaching comments counting, etc.

	var postDetails []*post.PostDetail
	// Fill details mock
	// ...

	return &post.GetFriendTimelineResp{
		List: postDetails,
	}, nil
}
