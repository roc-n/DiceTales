package config

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Mongo struct {
		Url string
		Db  string
	}
	Redis redis.RedisConf

	Snowflake struct {
		Node int64 `json:",default=1"`
	}
}
