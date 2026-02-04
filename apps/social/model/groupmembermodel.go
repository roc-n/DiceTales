package model

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ GroupMemberModel = (*customGroupMemberModel)(nil)

type (
	// GroupMemberModel is an interface to be customized, add more methods here,
	// and implement the added methods in customGroupMemberModel.
	GroupMemberModel interface {
		groupMemberModel

		InsertTx(ctx context.Context, session sqlx.Session, data *GroupMember) (sql.Result, error)
		ListByUserId(ctx context.Context, userId string) ([]*GroupMember, error)
		ListByGroupId(ctx context.Context, groupId string) ([]*GroupMember, error)
		FindByGroudIdAndUserId(ctx context.Context, userId, groupId string) (*GroupMember, error)
	}

	customGroupMemberModel struct {
		*defaultGroupMemberModel
	}
)

func (m *customGroupMemberModel) InsertTx(ctx context.Context, session sqlx.Session, data *GroupMember) (sql.Result, error) {
	groupMembersIdKey := fmt.Sprintf("%s%v", cacheGroupMemberIdPrefix, data.Id)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?)", m.table, groupMemberRowsExpectAutoSet)

		return session.ExecCtx(ctx, query, data.GroupId, data.Uid, data.RoleLevel, data.JoinTime, data.JoinSource, data.InviterUid, data.OperatorUid)
	}, groupMembersIdKey)
	return ret, err
}

func (m *customGroupMemberModel) ListByUserId(ctx context.Context, userId string) ([]*GroupMember, error) {
	query := fmt.Sprintf("select %s from %s where `uid` = ?", groupMemberRows, m.table)
	var resp []*GroupMember
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, userId)

	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *customGroupMemberModel) ListByGroupId(ctx context.Context, groupId string) ([]*GroupMember, error) {
	query := fmt.Sprintf("select %s from %s where `group_id` = ?", groupMemberRows, m.table)
	var resp []*GroupMember
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, groupId)

	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *customGroupMemberModel) FindByGroudIdAndUserId(ctx context.Context, uid, groupId string) (*GroupMember, error) {
	query := fmt.Sprintf("select %s from %s where `uid` = ? and `group_id` = ?", groupMemberRows, m.table)
	var resp GroupMember
	err := m.QueryRowNoCacheCtx(ctx, &resp, query, uid, groupId)

	switch err {
	case nil:
		return &resp, nil
	default:
		return nil, err
	}
}

// NewGroupMemberModel returns a model for the database table.
func NewGroupMemberModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) GroupMemberModel {
	return &customGroupMemberModel{
		defaultGroupMemberModel: newGroupMemberModel(conn, c, opts...),
	}
}
