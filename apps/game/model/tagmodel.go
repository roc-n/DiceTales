package model

import (
	"context"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ TagModel = (*customTagModel)(nil)

type (
	// TagModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTagModel.
	TagModel interface {
		tagModel
		FindByIds(ctx context.Context, ids []uint64) ([]*Tag, error)
	}

	customTagModel struct {
		*defaultTagModel
	}
)

// NewTagModel returns a model for the database table.
func NewTagModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) TagModel {
	return &customTagModel{
		defaultTagModel: newTagModel(conn, c, opts...),
	}
}

func (m *customTagModel) FindByIds(ctx context.Context, ids []uint64) ([]*Tag, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	// Create placeholder array
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}

	query := fmt.Sprintf("select %s from %s where `id` in (%s)", tagRows, m.table, strings.Join(placeholders, ","))
	var resp []*Tag
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, args...)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
