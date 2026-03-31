package config

import (
        "github.com/zeromicro/go-queue/kq"
        "github.com/zeromicro/go-zero/core/service"
        "github.com/zeromicro/go-zero/core/stores/redis"
)

type Config struct {
        service.ServiceConf

        PostEventTransfer kq.KqConf

        Redisx redis.RedisConf
        Mysql struct {
                DataSource string
        }
}
