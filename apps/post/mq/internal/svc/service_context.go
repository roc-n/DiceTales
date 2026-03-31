package svc

import (
        "dicetales.com/apps/post/mq/internal/config"
        "github.com/zeromicro/go-zero/core/stores/redis"
)

type ServiceContext struct {
        Config config.Config
        RedisClient *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
        return &ServiceContext{
                Config: c,
                RedisClient: redis.MustNewRedis(c.Redisx),
        }
}
