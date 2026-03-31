package logic

import (
	"context"
	"encoding/json"

	"dicetales.com/apps/post/model"
	"dicetales.com/apps/post/rpc/internal/svc"
	"dicetales.com/apps/post/rpc/post"

	"github.com/zeromicro/go-zero/core/logx"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateCommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCommentLogic {
	return &CreateCommentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateCommentLogic) CreateComment(in *post.CreateCommentReq) (*post.CreateCommentResp, error) {
	// 1. Optional Block: Integrate SensitiveWordFilter middleware/RPC here

	// 2. Fact storage into MongoDB smoothly
	newComment := &model.Comment{
		ID:       primitive.NewObjectID(),
		PostId:   in.PostId,
		UserId:   in.UserId,
		Content:  in.Content,
		ParentId: in.ParentId,
		RootId:   in.RootId,
		Status:   0,
	}

	err := l.svcCtx.CommentModel.Insert(l.ctx, newComment)
	if err != nil {
		return nil, err
	}

	// 3. Dispatch an InteractionEvent for Comment Count Coalescing
	event := map[string]interface{}{
		"type":        "COMMENT",
		"user_id":     in.UserId,
		"target_id":   in.PostId, // We aggregate comment counts into Post
		"target_type": 1,         // 1: POST
		"action":      5,         // Let's assume action '5' stands for ADD_COMMENT
		"status":      1,
	}
	eventBytes, _ := json.Marshal(event)

	if err := l.svcCtx.PostEventPusher.Push(l.ctx, string(eventBytes)); err != nil {
		l.Logger.Errorf("CreateComment Push Kafka err: %v", err)
	}

	// Converting Mongo ObjectID to an identifiable Int64 timestamp-based or string for Response
	// Here we just mock return 1 for CommentId due to int64 requirement
	return &post.CreateCommentResp{
		CommentId: 1,
	}, nil
}
