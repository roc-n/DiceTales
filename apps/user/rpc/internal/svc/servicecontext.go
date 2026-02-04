package svc

import (
	"dicetales.com/apps/user/model"
	"dicetales.com/apps/user/rpc/internal/config"
	"dicetales.com/pkg/id"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config
	redis.Redis
	model.UserModel
	IDGen id.Generator
}

func NewServiceContext(c config.Config) *ServiceContext {

	sqlConn := sqlx.NewMysql(c.Mysql.DataSource)
	rds := redis.MustNewRedis(c.Redisx)

	idStore := id.NewMysqlStore(sqlConn)
	idGen := id.NewGenerator(rds, idStore, "user")

	return &ServiceContext{
		Config:    c,
		Redis:     *rds,
		UserModel: model.NewUserModel(sqlConn, c.Cache),
		IDGen:     idGen,
	}
}
