package model

import (
	"errors"

	"github.com/zeromicro/go-zero/core/stores/mon"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var (
	ErrNotFound        = sqlx.ErrNotFound
	ErrNotFoundMon     = mon.ErrNotFound
	ErrInvalidObjectId = errors.New("invalid objectId")
)
