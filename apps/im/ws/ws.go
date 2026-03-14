package main

import (
	"flag"
	"fmt"

	"dicetales.com/apps/im/ws/core"
	"dicetales.com/apps/im/ws/internal/config"
	"dicetales.com/apps/im/ws/internal/handler"
	"dicetales.com/apps/im/ws/internal/svc"
	"github.com/zeromicro/go-zero/core/conf"
)

var configFile = flag.String("f", "etc/im.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	// 设置go-zero日志、监听等行为
	if err := c.SetUp(); err != nil {
		panic(err)
	}

	svc := svc.NewServiceContext(c)
	srv := core.NewServer(c.ListenOn,
		core.WithServerAuthentication(handler.NewJwtAuth(svc)),
		core.WithServerAck(core.NoAck),
		// websocket.WithServerMaxConnectionIdle(10*time.Second),
		// websocket.WithServerSensitive(),
		core.WithServerMsgLimiter(),
	)
	defer srv.Stop()

	handler.RegisterHandlers(srv, svc)
	fmt.Println("start websocket server at ", c.ListenOn, " ..... ")
	srv.Start()
}
