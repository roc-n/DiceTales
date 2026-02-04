package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ FriendModel = (*customFriendModel)(nil)

type (
	// FriendModel is an interface to be customized, add more methods here,
	// and implement the added methods in customFriendModel.
	FriendModel interface {
		friendModel

		FindByUidAndFid(ctx context.Context, uid, fid string) (*Friend, error)
		Inserts(ctx context.Context, session sqlx.Session, data ...*Friend) (sql.Result, error)
		ListByUserid(ctx context.Context, userId string) ([]*Friend, error)
	}

	customFriendModel struct {
		*defaultFriendModel
	}
)

func (m *customFriendModel) FindByUidAndFid(ctx context.Context, uid, fid string) (*Friend, error) {
	query := fmt.Sprintf(`select %s from %s where uid = ? and friend_uid = ?`, friendRows, m.table)

	var resp Friend
	err := m.QueryRowNoCacheCtx(ctx, &resp, query, uid, fid)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultFriendModel) Inserts(ctx context.Context, session sqlx.Session, data ...*Friend) (sql.Result, error) {
	var (
		sql  strings.Builder
		args []any
	)

	if len(data) == 0 {
		return nil, nil
	}

	// insert into tablename values(数据), (数据)
	sql.WriteString(fmt.Sprintf("insert into %s (%s) values ", m.table, friendRowsExpectAutoSet))

	for i, v := range data {
		sql.WriteString("(?, ?, ?, ?)")
		args = append(args, v.Uid, v.FriendUid, v.Remark, v.AddSource)
		if i == len(data)-1 {
			break
		}

		sql.WriteString(",")
	}

	return session.ExecCtx(ctx, sql.String(), args...)
}

func (m *defaultFriendModel) ListByUserid(ctx context.Context, userId string) ([]*Friend, error) {
	query := fmt.Sprintf("select %s from %s where `uid` = ? ", friendRows, m.table)

	var resp []*Friend
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, userId)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

// NewFriendModel returns a model for the database table.
func NewFriendModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) FriendModel {
	return &customFriendModel{
		defaultFriendModel: newFriendModel(conn, c, opts...),
	}
}
