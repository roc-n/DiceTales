package logic

import (
	"context"
	"fmt"

	"strconv"

	"dicetales.com/apps/game/rpc/game"
	"dicetales.com/apps/game/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FilterGamesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFilterGamesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FilterGamesLogic {
	return &FilterGamesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 高级组合筛选 (核心亮点)
func (l *FilterGamesLogic) FilterGames(in *game.FilterGamesReq) (*game.FilterGamesResp, error) {
	// 1. 组装参与 BITOP 运算的 Redis Key切片
	var bitKeys []string

	if in.PlayerCount > 0 {
		bitKeys = append(bitKeys, fmt.Sprintf("idx:game:player:%d", in.PlayerCount))
	}
	if in.DurationRange != "" {
		bitKeys = append(bitKeys, fmt.Sprintf("idx:game:duration:%s", in.DurationRange))
	}
	if in.ComplexityLevel != "" {
		bitKeys = append(bitKeys, fmt.Sprintf("idx:game:complexity:%s", in.ComplexityLevel))
	}
	for _, tagId := range in.TagIds {
		bitKeys = append(bitKeys, fmt.Sprintf("idx:game:tag:%d", tagId))
	}

	// 初始化一个临时 Key 用作结果集
	resultKey := "idx:game:temp:result"

	redisClient := l.svcCtx.Redis

	// 如果没有任何条件，直接报错或返回空（业务决定，暂时返回空）
	if len(bitKeys) == 0 {
		return &game.FilterGamesResp{}, nil
	}

	// 2. 将所有筛选条件的 Bitmap 执行 BITOP AND 交集计算
	// BITOP 命令: BITOP AND destkey srckey1 srckey2 ...
	args := []interface{}{"AND", resultKey}
	for _, key := range bitKeys {
		args = append(args, key)
	}

	// 注意这里因为 go-zero 原生 Redis 没有暴露直接的 BitOp 方法片段，我们通过原生 Send 发送
	_, err := redisClient.BitOpAndCtx(l.ctx, resultKey, bitKeys...)
	if err != nil {
		l.Errorf("redis BitOp error: %v", err)
		return nil, err
	}
	// 设置极短失效时间防脏数据
	_ = redisClient.ExpireCtx(l.ctx, resultKey, 10)

	// 3. 读取计算出来的交集 Bitmap 中的有效位（由于数据包较大，在实际中应搭配 LUA 或 Redisson，或者在Redis端获取位为1的偏移量列表）
	// 由于 go-zero Redis 未在此处提供现成的获取所有 Set 位的方法，可以通过获取整个字节串自己处理，
	// 面试代码中一般说明我们用应用层强算或者 Redis 的 BITPOS 循环（这里简化表示为全量抓取出所有 Game Ids）。

	// 为了演示这段“核心重难点：位运算联合及内存排序分页”，假设我们从 resultKey 解析出几十到几百个 gameID:
	// 这里用伪实现获取匹配到的 ID，后续可用 lua 或真实遍历来做
	gameIds := parseBitmapToIDs(redisClient, resultKey)

	if len(gameIds) == 0 {
		return &game.FilterGamesResp{Total: 0, GameIds: nil}, nil
	}

	// 4. 内存级联合排序与分页
	// 从 Hash 表中拿排序指标: 假设有一个 Hash 'idx:game:scores'，存储 ID -> Score 映射
	scoreMap := make(map[int64]float64)
	for _, id := range gameIds {
		// 单条 Get，生产中换为 HMGET
		scoreStr, _ := redisClient.HgetCtx(l.ctx, "idx:game:scores", strconv.FormatInt(id, 10))
		score, _ := strconv.ParseFloat(scoreStr, 64)
		scoreMap[id] = score
	}

	// 这里可以利用 Go 堆排/快速排序进行内存排序（略过实现，代表你掌握这个技术点即可）
	// sort.Slice() ...

	// 5. 分页截断 Pagination
	page := int(in.Page)
	if page < 1 {
		page = 1
	}
	size := int(in.Size)
	if size < 1 {
		size = 10
	}

	start := (page - 1) * size
	end := start + size

	var finalIds []int64
	if start < len(gameIds) {
		if end > len(gameIds) {
			end = len(gameIds)
		}
		finalIds = gameIds[start:end]
	}

	return &game.FilterGamesResp{
		Total:   int64(len(gameIds)),
		GameIds: finalIds,
	}, nil
}

// 模拟解析，实战中可查 Redis 原生方法或者直接读 []byte 解析每一位
func parseBitmapToIDs(rc interface{}, key string) []int64 {
	// Dummy, 假设通过 Bitmap 返回了 1，2，3
	return []int64{1, 2, 3}
}
