package model

import "github.com/zeromicro/go-zero/core/stores/mon"

var _ CommentModel = (*customCommentModel)(nil)

type (
	// CommentModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCommentModel.
	CommentModel interface {
		commentModel
	}

	customCommentModel struct {
		*defaultCommentModel
	}
)

// NewCommentModel returns a model for the mongo.
func NewCommentModel(url, db, collection string) CommentModel {
	conn := mon.MustNewModel(url, db, collection)
	return &customCommentModel{
		defaultCommentModel: newDefaultCommentModel(conn),
	}
}
