package logic

import (
        "context"
        "encoding/json"
        "fmt"

        "dicetales.com/apps/post/mq/internal/svc"
        "github.com/zeromicro/go-zero/core/logx"
)

type PostEventHandler struct {
        svcCtx *svc.ServiceContext
}

func NewPostEventHandler(svcCtx *svc.ServiceContext) *PostEventHandler {
        return &PostEventHandler{
                svcCtx: svcCtx,
        }
}

func (h *PostEventHandler) Consume(ctx context.Context, key, val string) error {
        logx.Infof("PostEventHandler receive msg: %s", val)

        var event map[string]interface{}
        if err := json.Unmarshal([]byte(val), &event); err != nil {
                logx.Errorf("unmarshal error: %v", err)
                return err
        }

        eventType, flag := event["type"].(string)
        if flag == false {
                return nil
        }

        switch eventType {
        case "CREATE_POST":
                logx.Infof("Processing CREATE_POST event mapping for post_id: %v", event["post_id"])
        case "INTERACT", "COMMENT":
                targetId, _ := event["target_id"].(float64)
                action, _ := event["action"].(float64) 

                aggKey := fmt.Sprintf("agg:post:%v", int64(targetId))
                
                if action == 1 {
                        h.svcCtx.RedisClient.HincrbyCtx(ctx, aggKey, "like", 1)
                } else if action == 2 {
                        h.svcCtx.RedisClient.HincrbyCtx(ctx, aggKey, "like", -1)
                } else if action == 5 { 
                        h.svcCtx.RedisClient.HincrbyCtx(ctx, aggKey, "comment", 1)
                }
        }
        return nil
}
