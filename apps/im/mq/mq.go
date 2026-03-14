package main

import (
	"flag"
	"fmt"

	"dicetales.com/apps/im/mq/internal/config"
	"dicetales.com/apps/im/mq/internal/handler"
	"dicetales.com/apps/im/mq/internal/svc"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
)

var configFile = flag.String("f", "etc/dev/mq.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	// 设置go-zero日志、监听等行为
	if err := c.SetUp(); err != nil {
		panic(err)
	}

	svc := svc.NewServiceContext(c)
	listen := handler.NewListen(svc)

	serviceGroup := service.NewServiceGroup()
	for _, s := range listen.Services() {
		serviceGroup.Add(s)
	}

	//websocket.WithServerAuthentication(handler.NewJwtAuth(ctx)),
	//websocket.WithServerAck(websocket.RigorAck),
	//websocket.WithServerMaxConnectionIdle(10*time.Second),
	// handler.RegisterHandlers(srv, ctx)
	fmt.Println("start mq server at ", c.ListenOn, " ..... ")
	serviceGroup.Start()
}
