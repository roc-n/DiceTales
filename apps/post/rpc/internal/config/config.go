package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf

	Mysql struct {
		DataSource string
	}
	Mongo struct {
		Url string
		Db  string
	}
	CacheRedis cache.CacheConf
	Redisx     redis.RedisConf

	PostEventPusher struct {
		Brokers []string
		Topic   string
	}
}
