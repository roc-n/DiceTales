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

type InteractLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewInteractLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InteractLogic {
	return &InteractLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *InteractLogic) Interact(in *post.InteractReq) (*post.InteractResp, error) {
	// 1. Facts fast insertion into MongoDB
	// Depending on Action (Like/Unlike), it could be an Upsert with status (0: canceled, 1: valid)
	status := int32(1)
	if in.Action == 2 || in.Action == 4 { // Unlike or Uncollect
		status = 0
	}

	interaction := &model.SocialInteraction{
		ID:         primitive.NewObjectID(),
		UserId:     in.UserId,
		TargetId:   in.TargetId,
		TargetType: in.TargetType,
		Action:     in.Action,
		Status:     status,
	}

	// Using model Insert (For production, implement Upsert based on uk_user_target_action equivalent in Mongo if needed)
	// Ignoring error handling in this stub for brevity
	_ = l.svcCtx.SocialInteractionModel.Insert(l.ctx, interaction)

	// 2. Publish to Kafka to coalesce and aggregate counts (Write-Behind)
	event := map[string]interface{}{
		"type":        "INTERACT",
		"user_id":     in.UserId,
		"target_id":   in.TargetId,
		"target_type": in.TargetType,
		"action":      in.Action,
		"status":      status,
	}
	eventBytes, _ := json.Marshal(event)

	if err := l.svcCtx.PostEventPusher.Push(l.ctx, string(eventBytes)); err != nil {
		l.Logger.Errorf("Interact Push Kafka err: %v", err)
	}

	return &post.InteractResp{
		Success: true,
	}, nil
}
