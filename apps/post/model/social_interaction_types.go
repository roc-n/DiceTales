package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SocialInteraction struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserId     int64              `bson:"userId,omitempty" json:"userId,omitempty"`
	TargetId   int64              `bson:"targetId,omitempty" json:"targetId,omitempty"`
	TargetType int32              `bson:"targetType,omitempty" json:"targetType,omitempty"`
	Action     int32              `bson:"action,omitempty" json:"action,omitempty"`
	Status     int32              `bson:"status,omitempty" json:"status,omitempty"`
	CreateAt   time.Time          `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	UpdateAt   time.Time          `bson:"updateAt,omitempty" json:"updateAt,omitempty"`
}
