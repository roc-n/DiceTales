# DiceTales 集中式 API 层 (BFF) 架构与实现规划

## 1. 架构定位：BFF (Backend For Frontend)

在微服务架构中，如果每个 RPC 都配套一个独立的 API 服务，会导致前端需要面对多个 API 入口，且跨域数据聚合非常困难。
因此，在 DiceTales 项目中，我们引入一个**集中式的 API 服务（BFF层）**。

- **职责**：
  - **请求路由与参数校验**：接收前端 HTTP 请求，校验基础参数。
  - **统一鉴权**：在这一层进行双 Token 的校验和拦截。
  - **数据聚合（裁剪/组装）**：比如查询“某个组局（meetup）”的详情时，需要同时调用 `meetup-rpc` 获取组局信息，调用 `game-rpc` 获取桌游信息，调用 `user-rpc` 获取发起人信息。这些组装逻辑全部在 BFF 层完成，保持各 RPC 的纯粹性。
- **与真·API网关的区别**：你后续独立研发的网关负责限流、熔断、黑白名单、全链路追踪接入等网络层共性；而这个 BFF 层负责与**业务强相关**的聚合和适配。

---

## 2. 目录结构设计

我们将在 `apps` 目录下新建一个统一的 `api` 目录（即 `apps/api`）。

推荐结构如下：
```text
apps/api/
  ├── api.api            // 总入口定义文件
  ├── desc/              // 按模块拆分的 api 定义文件
  │   ├── user.api
  │   ├── game.api
  │   ├── meetup.api
  │   ├── social.api
  │   └── post.api
  ├── etc/
  │   └── api-api.yaml   // 统一服务的配置文件，包含所有 RPC 客户端配置
  ├── internal/
  │   ├── config/        // 配置信息 (包含所有 RPC client 的配置结构)
  │   ├── handler/       // 路由与 handler (由 goctl 根据 group 分组自动生成)
  │   │   ├── user/
  │   │   ├── game/
  │   │   └── meetup/
  │   ├── logic/         // 业务逻辑 (如数据聚合、调用多个 RPC)
  │   │   ├── user/
  │   │   ├── game/
  │   │   └── meetup/
  │   ├── svc/           // ServiceContext (注入所有的 rpc clients)
  │   └── types/         // 生成的请求响应结构体
  └── api.go             // main 函数
```

---

## 3. Go-zero API 切分与聚合实践

为了避免单个 `.api` 文件过长，我们使用 go-zero 的 `import` 语法和 `group` 属性来实现类似于 Gin Router Group 的效果。

### 3.1 拆分定义 (以 `user.api` 和 `game.api` 为例)

`apps/api/desc/user.api`:
```go
syntax = "v1"

type (
    LoginReq { ... }
    LoginResp { ... }
)

@server(
    prefix: /v1/user
    group: user
)
service api-api {
    @handler Login
    post /login (LoginReq) returns (LoginResp)
}
```

`apps/api/desc/game.api`:
```go
syntax = "v1"

type (
    GameInfoReq { ... }
    GameInfoResp { ... }
)

@server(
    jwt: Auth // 统一应用 jwt
    prefix: /v1/game
    group: game
)
service api-api {
    @handler GetGameInfo
    get /info (GameInfoReq) returns (GameInfoResp)
}
```

### 3.2 总入口集成

`apps/api/api.api`:
```go
syntax = "v1"

import "desc/user.api"
import "desc/game.api"
import "desc/meetup.api"
import "desc/social.api"
```
使用统一命令生成：`goctl api go -api api.api -dir .`。这样 goctl 就会把处理函数按不同的业务域存放在对应的 `internal/logic/<group>/` 目录下。

---

## 4. 依赖注入与上下文配置

在 `ServiceContext` 中统一持有所有微服务 RPC 客户端引用：

`apps/api/internal/svc/servicecontext.go`:
```go
type ServiceContext struct {
    Config      config.Config
    UserRpc     userclient.User
    GameRpc     gameclient.Game
    MeetupRpc   meetupclient.Meetup
    SocialRpc   socialclient.Social
    // ImRpc ...
}

func NewServiceContext(c config.Config) *ServiceContext {
    return &ServiceContext{
        Config:    c,
        UserRpc:   userclient.NewUser(zrpc.MustNewClient(c.UserRpcConf)),
        GameRpc:   gameclient.NewGame(zrpc.MustNewClient(c.GameRpcConf)),
        // ...
    }
}
```
此时，BFF 中的任何一个 Logic（例如获取动态详情，需要附带发帖人信息时），都能非常方便地同时调用 `l.svcCtx.PostRpc` 和 `l.svcCtx.UserRpc` 进行数据拼装。

---

## 5. 迁移与实施步骤

1. **第一步（脚手架搭建）**：在 `apps/api` 中创建总控结构，将现有的 `apps/user/api/` 的内容提取为 `apps/api/desc/user.api`。
2. **第二步（废弃旧层）**：删除原来的 `apps/user/api/`（移除历史债务）。
3. **第三步（打通双 Token）**：基于新的 `apps/api/desc/user.api` 和前面规划好的 `etc/insights/dual-token-auth-design.md`，在这个统一的 BFF 层实现用户登录和 `/v1/user/token/refresh`。
4. **第四步（预留与逐步扩展模块）**：由于目前 `game`、`meetup`、`post` 模块均未真正实现其 RPC，在 `apps/api/desc/` 下我们先为其创建占位性质的 `.api` 语法存根，并在主模块中 `import`，以便后续有了 RPC 服务后直接在 API 层填充对应的请求响应头和转发逻辑。
5. **第五步（清理代码）**：迁移完成后，统一移除旧的分散 API 服务目录。的