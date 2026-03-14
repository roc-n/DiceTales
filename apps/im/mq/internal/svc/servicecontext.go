package svc

import (
	"net/http"

	"dicetales.com/apps/im/model"
	"dicetales.com/apps/im/mq/internal/config"
	"dicetales.com/apps/im/ws/core/client"
	"dicetales.com/apps/social/rpc/socialclient"
	"dicetales.com/pkg/constants"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	config.Config

	WsClient client.Client
	*redis.Redis
	model.ChatLogModel
	model.ConversationModel

	socialclient.Social
}

func NewServiceContext(c config.Config) *ServiceContext {
	svc := &ServiceContext{
		Config:            c,
		Redis:             redis.MustNewRedis(c.Redisx),
		ChatLogModel:      model.MustChatLogModel(c.Mongo.Url, c.Mongo.Db),
		ConversationModel: model.MustConversationModel(c.Mongo.Url, c.Mongo.Db),

		Social: socialclient.NewSocial(zrpc.MustNewClient(c.SocialRpc)),
	}

	token, err := svc.GetSystemToken()
	if err != nil {
		panic(err)
	}

	header := http.Header{}
	header.Set("Authorization", token)
	svc.WsClient = client.NewClient(c.Ws.Host, client.WithClientHeader(header))

	return svc
}

func (svc *ServiceContext) GetSystemToken() (string, error) {
	return svc.Redis.Get(constants.REDIS_SYSTEM_ROOT_TOKEN)
}
