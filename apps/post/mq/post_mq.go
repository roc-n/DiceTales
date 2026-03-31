package main

import (
        "flag"
        "fmt"

        "dicetales.com/apps/post/mq/internal/config"
        "dicetales.com/apps/post/mq/internal/logic"
        "dicetales.com/apps/post/mq/internal/svc"

        "github.com/zeromicro/go-queue/kq"
        "github.com/zeromicro/go-zero/core/conf"
)

var configFile = flag.String("f", "etc/mq.yaml", "the config file")

func main() {
        flag.Parse()

        var c config.Config
        conf.MustLoad(*configFile, &c)

        ctx := svc.NewServiceContext(c)

        srv := kq.MustNewQueue(c.PostEventTransfer, kq.WithHandle(logic.NewPostEventHandler(ctx).Consume))
        defer srv.Stop()

        fmt.Println("post_mq is starting...")
        srv.Start()
}
