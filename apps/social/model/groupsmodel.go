package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ GroupsModel = (*customGroupsModel)(nil)

type (
	// GroupsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customGroupsModel.
	GroupsModel interface {
		groupsModel

		Trans(ctx context.Context, fn func(context.Context, sqlx.Session) error) error
		InsertTx(ctx context.Context, session sqlx.Session, data *Groups) (sql.Result, error)
		ListByGroupIds(ctx context.Context, ids []string) ([]*Groups, error)
		UpdateNotification(ctx context.Context, groupId, notification, notificationUid string) error
	}

	customGroupsModel struct {
		*defaultGroupsModel
	}
)

func (m *customGroupsModel) Trans(ctx context.Context, fn func(context.Context, sqlx.Session) error) error {
	return m.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		return fn(ctx, session)
	})
}

func (m *customGroupsModel) InsertTx(ctx context.Context, session sqlx.Session, data *Groups) (sql.Result, error) {
	groupsIdKey := fmt.Sprintf("%s%v", cacheGroupsIdPrefix, data.Id)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, groupsRowsExpectAutoSet)
		return session.ExecCtx(ctx, query, data.Id, data.Name, data.Icon, data.Status, data.CreatorUid, data.GroupType, data.IsVerify, data.Notification, data.NotificationUid, data.DeletedAt)
	}, groupsIdKey)
	return ret, err
}

func (m *customGroupsModel) ListByGroupIds(ctx context.Context, ids []string) ([]*Groups, error) {

	var resp []*Groups
	placeholders := strings.Repeat("?,", len(ids)-1) + "?"
	query := fmt.Sprintf("select %s from %s where id in (%s)", groupsRows, m.table, placeholders)
	args := make([]any, len(ids))
	for i, id := range ids {
		args[i] = id
	}
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, args...)

	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *customGroupsModel) UpdateNotification(ctx context.Context, groupId, notification, notificationUid string) error {
	groupsIdKey := fmt.Sprintf("%s%v", cacheGroupsIdPrefix, groupId)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set `notification` = ?, `notification_uid` = ? where `id` = ? ", m.table)
		return conn.ExecCtx(ctx, query, notification, notificationUid, groupId)
	}, groupsIdKey)
	return err
}

// NewGroupsModel returns a model for the database table.
func NewGroupsModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) GroupsModel {
	return &customGroupsModel{
		defaultGroupsModel: newGroupsModel(conn, c, opts...),
	}
}
