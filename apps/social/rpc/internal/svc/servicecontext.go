package svc

import (
	"dicetales.com/apps/social/model"
	"dicetales.com/apps/social/rpc/internal/config"
	"dicetales.com/pkg/id"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config

	model.FriendModel
	model.FriendRequestModel
	model.GroupsModel
	model.GroupRequestModel
	model.GroupMemberModel

	IDGen id.Generator
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewMysql(c.Mysql.DataSource)
	rds := redis.MustNewRedis(c.Redisx)

	idGen := id.NewGeneratorWithTable(rds, sqlConn, "social", "group_id_pool")

	return &ServiceContext{
		Config: c,

		FriendModel:        model.NewFriendModel(sqlConn, c.Cache),
		FriendRequestModel: model.NewFriendRequestModel(sqlConn, c.Cache),
		GroupsModel:        model.NewGroupsModel(sqlConn, c.Cache),
		GroupRequestModel:  model.NewGroupRequestModel(sqlConn, c.Cache),
		GroupMemberModel:   model.NewGroupMemberModel(sqlConn, c.Cache),

		IDGen: idGen,
	}
}
