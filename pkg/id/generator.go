package id

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

const (
	defaultBatchSize = 10000 // 单次扩容生成量
	defaultWaterMark = 200   // 低水位阈值, 低于此值触发扩容
	defaultStartID   = 9999  // 初始ID起始值
	lockExpire       = 10    // 分布式锁过期时间(秒)
)

type Generator interface {
	Get(ctx context.Context) (string, error)
	Confirm(ctx context.Context, idStr string) error
}

type RedisGenerator struct {
	rds       *redis.Redis
	db        Store // 依赖数据库存储接口
	keyPool   string
	keyLock   string
	batchSize int
	waterMark int
}

// NewGenerator 创建基于Redis+DB的号码池生成器
// serviceName: 业务名称, 用于区分Redis Key (如 id:pool:user, id:pool:group)
// store: 数据库存储实现, 用于持久化号码池数据
func NewGenerator(rds *redis.Redis, store Store, serviceName string) *RedisGenerator {
	return &RedisGenerator{
		rds:       rds,
		db:        store,
		keyPool:   fmt.Sprintf("id:pool:%s", serviceName),
		keyLock:   fmt.Sprintf("id:lock:%s", serviceName),
		batchSize: defaultBatchSize,
		waterMark: defaultWaterMark,
	}
}

// NewGeneratorWithTable 快捷创建指定表名和业务名的生成器
func NewGeneratorWithTable(rds *redis.Redis, conn sqlx.SqlConn, serviceName string, tableName string) *RedisGenerator {
	return NewGenerator(rds, NewStoreWithTable(conn, tableName), serviceName)
}

// SetBatchSize 设置单次扩容数量
func (g *RedisGenerator) SetBatchSize(size int) *RedisGenerator {
	if size > 0 {
		g.batchSize = size
	}
	return g
}

// Get 获取一个唯一ID
func (g *RedisGenerator) Get(ctx context.Context) (string, error) {
	// 尝试从Redis取号
	val, err := g.rds.Spop(g.keyPool)
	if err == nil {
		g.asyncCheck() // 取号成功后, 异步检查水位
		return val, nil
	}
	if err != redis.Nil {
		return "", errors.Wrap(err, "redis spop failed")
	}

	// Redis为空, 同步触发补货
	if err := g.replenish(); err != nil {
		return "", errors.Wrap(err, "replenish failed")
	}

	// 补货后重试
	val, err = g.rds.Spop(g.keyPool)
	if err != nil {
		return "", errors.Wrap(err, "failed to get id after replenish")
	}
	return val, nil
}

func (g *RedisGenerator) asyncCheck() {
	go func() {
		_ = g.checkAndReplenish()
	}()
}

func (g *RedisGenerator) checkAndReplenish() error {
	count, err := g.rds.Scard(g.keyPool)
	if err != nil || int(count) > g.waterMark {
		return err
	}
	return g.replenish()
}

// replenish 补货核心逻辑: Redis缺货 -> 从DB取 -> DB缺货 -> 生成新段入库 -> 从DB取 -> 入Redis
func (g *RedisGenerator) replenish() error {
	lock := redis.NewRedisLock(g.rds, g.keyLock)
	lock.SetExpire(lockExpire)
	if ok, _ := lock.Acquire(); !ok {
		return nil // 拿不到锁说明正在补货, 直接返回
	}
	defer func() { _, _ = lock.Release() }()

	// Double check
	if count, _ := g.rds.Scard(g.keyPool); int(count) > g.waterMark {
		return nil
	}

	// 1. 尝试从 DB 获取一批可用号码 (status=0)
	ids, err := g.db.GetAvailable(context.Background(), g.batchSize)
	if err != nil {
		return err
	}

	// 2. 如果 DB 库存不足, 生成新号段
	if len(ids) < g.batchSize/2 { // DB扩展策略: 少于一半扩容
		if err := g.db.GenerateNewSegment(context.Background(), g.batchSize); err != nil {
			logx.Errorf("failed to generate new segment: %v", err)
			return err
		}
		// 重新拉取一次
		ids, err = g.db.GetAvailable(context.Background(), g.batchSize)
		if err != nil {
			return err
		}
	}

	if len(ids) == 0 {
		return fmt.Errorf("no available ids found")
	}

	// 3. 将 ids 打散后存入 Redis
	args := make([]any, len(ids))
	shuffle(ids)
	for i, id := range ids {
		args[i] = id
	}

	if _, err := g.rds.Sadd(g.keyPool, args...); err != nil {
		return err
	}

	logx.Infof("replenish success, added %d ids", len(ids))
	return nil
}

func shuffle(a []int64) {
	rand.Shuffle(len(a), func(i, j int) { a[i], a[j] = a[j], a[i] })
}

// Confirm 确认该ID已被正式使用, 将数据库状态置为1
func (g *RedisGenerator) Confirm(ctx context.Context, idStr string) error {
	// 简单的 string -> int64 转换
	var idVal int64
	_, err := fmt.Sscanf(idStr, "%d", &idVal)
	if err != nil {
		return errors.Wrapf(err, "invalid id format: %s", idStr)
	}

	return g.db.UpdateStatus(ctx, []int64{idVal}, 1)
}
