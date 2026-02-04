package model

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ FriendRequestModel = (*customFriendRequestModel)(nil)

type (
	// FriendRequestModel is an interface to be customized, add more methods here,
	// and implement the added methods in customFriendRequestModel.
	FriendRequestModel interface {
		friendRequestModel

		FindByReqUidAndUserId(ctx context.Context, rid, uid string) (*FriendRequest, error)
		Trans(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error
		ListNoHandleOut(ctx context.Context, uid string) ([]*FriendRequest, error)
		ListNoHandleIn(ctx context.Context, uid string) ([]*FriendRequest, error)
		UpdateTx(ctx context.Context, session sqlx.Session, data *FriendRequest) error
	}

	customFriendRequestModel struct {
		*defaultFriendRequestModel
	}
)

func (m *customFriendRequestModel) FindByReqUidAndUserId(ctx context.Context, rid, uid string) (*FriendRequest, error) {
	query := fmt.Sprintf("select %s from %s where `req_uid` = ? and `uid` = ?", friendRequestRows, m.table)

	var resp FriendRequest
	err := m.QueryRowNoCacheCtx(ctx, &resp, query, rid, uid)

	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customFriendRequestModel) ListNoHandleOut(ctx context.Context, uid string) ([]*FriendRequest, error) {
	query := fmt.Sprintf("select %s from %s where `uid` = ? and `handle_result` = 1 ", friendRequestRows, m.table)
	var resp []*FriendRequest
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, uid)

	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *customFriendRequestModel) ListNoHandleIn(ctx context.Context, uid string) ([]*FriendRequest, error) {
	query := fmt.Sprintf("select %s from %s where `req_uid` = ? and `handle_result` = 1 ", friendRequestRows, m.table)
	var resp []*FriendRequest
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, uid)

	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *customFriendRequestModel) Trans(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error {
	return m.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		return fn(ctx, session)
	})
}

func (m *defaultFriendRequestModel) UpdateTx(ctx context.Context, session sqlx.Session, data *FriendRequest) error {
	friendRequestsIdKey := fmt.Sprintf("%s%v", cacheFriendRequestIdPrefix, data.Id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, friendRequestRowsWithPlaceHolder)
		return session.ExecCtx(ctx, query, data.Uid, data.ReqUid, data.ReqMsg, data.ReqTime, data.ReqRemark, data.HandleResult, data.HandleMsg, data.HandledAt, data.Id)
	}, friendRequestsIdKey)
	return err
}

// NewFriendRequestModel returns a model for the database table.
func NewFriendRequestModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) FriendRequestModel {
	return &customFriendRequestModel{
		defaultFriendRequestModel: newFriendRequestModel(conn, c, opts...),
	}
}
