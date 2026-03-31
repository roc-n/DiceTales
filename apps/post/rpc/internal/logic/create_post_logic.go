package logic

import (
	"context"
	"database/sql"
	"encoding/json"

	"dicetales.com/apps/post/model"
	"dicetales.com/apps/post/rpc/internal/svc"
	"dicetales.com/apps/post/rpc/post"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreatePostLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreatePostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatePostLogic {
	return &CreatePostLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreatePostLogic) CreatePost(in *post.CreatePostReq) (*post.CreatePostResp, error) {
	imagesBytes, _ := json.Marshal(in.Images)

	// 1. Insert into MySQL Post table
	newPost := &model.Post{
		UserId:        uint64(in.UserId),
		Content:       sql.NullString{String: in.Content, Valid: in.Content != ""},
		Images:        sql.NullString{String: string(imagesBytes), Valid: true},
		Visibility:    int64(in.Visibility),
		RelatedGameId: uint64(in.RelatedGameId),
	}

	res, err := l.svcCtx.PostModel.Insert(l.ctx, newPost)
	if err != nil {
		l.Logger.Errorf("CreatePost Insert DB err: %v", err)
		return nil, err
	}

	postId, _ := res.LastInsertId()

	// 2. Publish to Kafka to fan-out to friend's timeline Inbox (Push mode)
	event := map[string]interface{}{
		"type":       "CREATE_POST",
		"post_id":    postId,
		"user_id":    in.UserId,
		"visibility": in.Visibility,
	}
	eventBytes, _ := json.Marshal(event)

	if err := l.svcCtx.PostEventPusher.Push(l.ctx, string(eventBytes)); err != nil {
		l.Logger.Errorf("CreatePost Push Kafka err: %v", err)
		// Note: Not returning error here to maintain high availability; push failure shouldn't fail the primary request
	}

	return &post.CreatePostResp{
		PostId: postId,
	}, nil
}
