package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Comment struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	PostId    int64              `bson:"postId,omitempty" json:"postId,omitempty"`
	UserId    int64              `bson:"userId,omitempty" json:"userId,omitempty"`
	Content   string             `bson:"content,omitempty" json:"content,omitempty"`
	ParentId  int64              `bson:"parentId,omitempty" json:"parentId,omitempty"`
	RootId    int64              `bson:"rootId,omitempty" json:"rootId,omitempty"`
	LikeCount int32              `bson:"likeCount,omitempty" json:"likeCount,omitempty"`
	Status    int32              `bson:"status,omitempty" json:"status,omitempty"`
	CreateAt  time.Time          `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	UpdateAt  time.Time          `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
}
