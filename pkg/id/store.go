package id

import (
	"context"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// Store 定义持久层接口, 用于管理号码池数据
// 作用:
// 1. 作为发号器的兜底仓库, 保证数据不丢失
// 2. 记录每个号码的状态(已分配/保留等)
type Store interface {
	// GetAvailable 从数据库获取一批未分配的号码
	GetAvailable(ctx context.Context, limit int) ([]int64, error)
	// GenerateNewSegment 生成并持久化新的号段
	GenerateNewSegment(ctx context.Context, step int) error
	// GetMaxID 获取当前库中最大的ID
	GetMaxID(ctx context.Context) (int64, error)
	// UpdateStatus 批量更新号码状态
	UpdateStatus(ctx context.Context, ids []int64, status int) error
}

type mysqlStore struct {
	conn      sqlx.SqlConn
	tableName string
}

// NewMysqlStore 创建默认的 ID 存储 (针对 id_pool 表)
// 默认用于用户 ID 生成
func NewMysqlStore(conn sqlx.SqlConn) Store {
	return &mysqlStore{
		conn:      conn,
		tableName: "id_pool",
	}
}

// NewStoreWithTable 创建指定表名的 ID 存储
// 可用于创建针对 group_id_pool 等其他表的存储实例
func NewStoreWithTable(conn sqlx.SqlConn, tableName string) Store {
	return &mysqlStore{
		conn:      conn,
		tableName: tableName,
	}
}

func (s *mysqlStore) GetAvailable(ctx context.Context, limit int) ([]int64, error) {
	var ids []int64
	// status=0 表示未分配
	query := fmt.Sprintf("SELECT id FROM %s WHERE status = 0 ORDER BY id ASC LIMIT ?", s.tableName)
	err := s.conn.QueryRowsCtx(ctx, &ids, query, limit)
	if err != nil {
		return nil, err
	}
	return ids, nil
}

func (s *mysqlStore) GetMaxID(ctx context.Context) (int64, error) {
	var maxId int64
	query := fmt.Sprintf("SELECT COALESCE(MAX(id), 0) FROM %s", s.tableName)
	err := s.conn.QueryRowCtx(ctx, &maxId, query)
	return maxId, err
}

func (s *mysqlStore) GenerateNewSegment(ctx context.Context, step int) error {
	maxId, err := s.GetMaxID(ctx)
	if err != nil {
		return err
	}

	// 如果是全新初始化, 且 maxId 为 0, 可以设定一个起始值
	start := maxId + 1
	if maxId == 0 {
		start = defaultStartID + 1
	}

	// 批量插入
	// 注意: 构建一个大 Values 字符串
	values := make([]string, 0, step)
	args := make([]any, 0, step)

	for i := range step {
		values = append(values, "(?, 0)") // id, status=0
		args = append(args, start+int64(i))
	}

	query := fmt.Sprintf("INSERT INTO %s (id, status) VALUES %s", s.tableName, strings.Join(values, ","))

	// 这里假设 step 不会设置得大到超出 MySQL packet size
	_, err = s.conn.ExecCtx(ctx, query, args...)
	return err
}

func (s *mysqlStore) UpdateStatus(ctx context.Context, ids []int64, status int) error {
	if len(ids) == 0 {
		return nil
	}

	// 拼接 WHERE IN id (?, ?, ...)
	query := fmt.Sprintf("UPDATE %s SET status = ? WHERE id IN (?%s)", s.tableName, strings.Repeat(",?", len(ids)-1))

	// 构建 args
	args := make([]any, 0, len(ids)+1)
	args = append(args, status)
	for _, id := range ids {
		args = append(args, id)
	}

	_, err := s.conn.ExecCtx(ctx, query, args...)
	return err
}
