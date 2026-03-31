package model

import "github.com/zeromicro/go-zero/core/stores/mon"

var _ SocialInteractionModel = (*customSocialInteractionModel)(nil)

type (
	// SocialInteractionModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSocialInteractionModel.
	SocialInteractionModel interface {
		socialInteractionModel
	}

	customSocialInteractionModel struct {
		*defaultSocialInteractionModel
	}
)

// NewSocialInteractionModel returns a model for the mongo.
func NewSocialInteractionModel(url, db, collection string) SocialInteractionModel {
	conn := mon.MustNewModel(url, db, collection)
	return &customSocialInteractionModel{
		defaultSocialInteractionModel: newDefaultSocialInteractionModel(conn),
	}
}
