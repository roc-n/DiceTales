package model

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ GroupRequestModel = (*customGroupRequestModel)(nil)

type (
	// GroupRequestModel is an interface to be customized, add more methods here,
	// and implement the added methods in customGroupRequestModel.
	GroupRequestModel interface {
		groupRequestModel

		Trans(ctx context.Context, fn func(context.Context, sqlx.Session) error) error
		UpdateTx(ctx context.Context, session sqlx.Session, data *GroupRequest) error

		FindByGroupIdAndReqId(ctx context.Context, groupId, reqId string) (*GroupRequest, error)
		ListNoHandle(ctx context.Context, groupId string) ([]*GroupRequest, error)
		ListByReqUid(ctx context.Context, reqUid string) ([]*GroupRequest, error)
	}

	customGroupRequestModel struct {
		*defaultGroupRequestModel
	}
)

func (m *customGroupRequestModel) Trans(ctx context.Context, fn func(context.Context, sqlx.Session) error) error {
	return m.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		return fn(ctx, session)
	})
}

func (m *customGroupRequestModel) FindByGroupIdAndReqId(ctx context.Context, groupId, reqId string) (*GroupRequest, error) {
	query := fmt.Sprintf("select %s from %s where `req_uid` = ? and `group_id` = ?", groupRequestRows, m.table)
	var resp GroupRequest
	err := m.QueryRowNoCacheCtx(ctx, &resp, query, reqId, groupId)

	switch err {
	case nil:
		return &resp, nil
	default:
		return nil, err
	}
}

func (m *customGroupRequestModel) ListNoHandle(ctx context.Context, groupId string) ([]*GroupRequest, error) {
	query := fmt.Sprintf("select %s from %s where `group_id` = ? and `handle_result` = 1 ", groupRequestRows, m.table)
	var resp []*GroupRequest
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, groupId)

	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *customGroupRequestModel) UpdateTx(ctx context.Context, session sqlx.Session, data *GroupRequest) error {
	groupRequestsIdKey := fmt.Sprintf("%s%v", cacheGroupRequestIdPrefix, data.Id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, groupRequestRowsWithPlaceHolder)
		return session.ExecCtx(ctx, query, data.ReqUid, data.GroupId, data.ReqMsg, data.ReqTime, data.JoinSource,
			data.InviterUid, data.HandleUid, data.HandleTime, data.HandleResult, data.Id)
	}, groupRequestsIdKey)
	return err
}

func (m *customGroupRequestModel) ListByReqUid(ctx context.Context, reqUid string) ([]*GroupRequest, error) {
	query := fmt.Sprintf("select %s from %s where `req_uid` = ? order by req_time desc", groupRequestRows, m.table)
	var resp []*GroupRequest
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, reqUid)

	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

// NewGroupRequestModel returns a model for the database table.
func NewGroupRequestModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) GroupRequestModel {
	return &customGroupRequestModel{
		defaultGroupRequestModel: newGroupRequestModel(conn, c, opts...),
	}
}
