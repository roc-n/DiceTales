package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserModel = (*customUserModel)(nil)

type (
	// UserModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserModel.
	UserModel interface {
		userModel

		FindByPhone(ctx context.Context, phone string) (*User, error)
		FindById(ctx context.Context, id string) (*User, error)
		ListByName(ctx context.Context, name string) ([]*User, error)
	}

	customUserModel struct {
		*defaultUserModel
	}
)

func (m *customUserModel) FindByPhone(ctx context.Context, phone string) (*User, error) {
	userPhoneKey := fmt.Sprintf("%s%v", cacheUserIdPrefix, phone)
	var resp User
	err := m.QueryRowCtx(ctx, &resp, userPhoneKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("select %s from %s where `phone` = ? limit 1", userRows, m.table)
		return conn.QueryRowCtx(ctx, v, query, phone)
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customUserModel) ListByName(ctx context.Context, name string) ([]*User, error) {
	query := fmt.Sprintf("select %s from %s where `nickname` like ? ", userRows, m.table)

	var resp []*User
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, fmt.Sprint("%", name, "%"))

	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *customUserModel) FindById(ctx context.Context, id string) (*User, error) {
	userIdKey := fmt.Sprintf("%s%v", cacheUserIdPrefix, id)
	var resp User
	err := m.QueryRowCtx(ctx, &resp, userIdKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", userRows, m.table)
		return conn.QueryRowCtx(ctx, v, query, id)
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// NewUserModel returns a model for the database table.
func NewUserModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) UserModel {
	return &customUserModel{
		defaultUserModel: newUserModel(conn, c, opts...),
	}
}
