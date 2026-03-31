# DiceTales 双 Token 认证体系落地方案（Access 15min + Refresh 7d）

## 0. 本轮已拍板决策（2026-03-31）

- Refresh Token 形态：采用 JWT + Redis 白名单。
- 多端策略：单设备登录（新登录会顶掉旧会话）。
- 兼容策略：保留 1 个版本周期的单 token 兼容字段。

## 1. 目标与边界

目标：在当前 DiceTales 项目中落地简历亮点：

- Access Token：短效（15 分钟）、无状态校验，保障 API/WS 吞吐。
- Refresh Token：长效（7 天）、Redis 白名单管理，支持即时撤销与旋转。
- 会话策略：登录成功后客户端持有 Access + Refresh；后续 HTTP/WS 走 Access；Access 过期后调用 Refresh 接口换新。

边界：

- 首期聚焦 `apps/user`（登录、刷新、登出）与 `apps/im/ws`（Access 校验兼容）。
- 不强制改数据库表结构（可先完成无侵入方案），可选增强项后置。

---

## 2. 当前实现基线（基于现有代码）

### 2.1 单 token + 超长有效期

- `apps/user/rpc/internal/logic/loginlogic.go`：登录后调用 `auth.GetJwtToken` 返回单个 `Token`。
- `apps/user/rpc/internal/logic/registerlogic.go`：注册后同样返回单个 `Token`。
- `apps/user/rpc/etc/user.yaml` 与 `apps/api/etc/api-api.yaml`：`AccessExpire` 当前配置为 `8640000`（长效）。
- `pkg/auth/token.go`：仅有通用 JWT 生成函数，未区分 Access/Refresh。

### 2.2 Access 校验链路

- `apps/api/internal/handler/routes.go`：受保护路由通过 `rest.WithJwt` 进行 JWT 校验。
- `apps/im/ws/internal/handler/auth.go`：WS 握手通过 go-zero token parser 校验 JWT。

### 2.3 API 层尚未落地业务转发

- `apps/api/internal/logic/user/loginlogic.go` 与 `registerlogic.go` 仍为 TODO。

结论：当前属于“单 JWT 长会话”，未具备刷新令牌白名单、令牌旋转与即时撤销能力。

---

## 3. 目标设计

## 3.1 Token 模型

- Access Token（JWT）
  - 建议声明：`uid`（或沿用现有 identify claim）、`iat`、`exp`、`typ=access`。
  - TTL：`900s`（15 分钟）。
  - 校验方式：纯本地签名校验（不查 Redis）。

- Refresh Token（JWT 或随机串）
  - 推荐用 JWT（便于携带 `uid/jti/dev`）+ Redis 白名单双重约束。
  - 必要声明：`uid`、`jti`、`typ=refresh`、`iat`、`exp`。
  - TTL：`604800s`（7 天）。
  - Redis 白名单：保存“可用 refresh 会话”，不存在即视为撤销。

说明：

- 为避免 Redis 泄露后直接重放，Redis 内建议保存 Refresh Token 哈希（如 SHA256）而非明文。

## 3.2 Redis 白名单设计

建议 key 结构：

- `auth:refresh:wl:{uid}:{jti}` -> `sha256(refreshToken)`，TTL=7d
- 可选：`auth:refresh:user:{uid}`（Set）用于批量下线/审计

关键策略：

- 登录：签发 refresh，写入白名单。
- 登录：单设备策略下，先清理该用户历史 refresh 白名单，再签发新 refresh 并写入白名单。
- 刷新：必须命中白名单且哈希一致；成功后执行“旋转”
  - 删除旧 `jti`
  - 签发新 refresh（新 `jti`）并写入白名单
  - 同时签发新 access
- 登出：删除当前 refresh `jti`；需要“踢下线全部设备”时删用户集合全部 `jti`。

## 3.3 接口语义

- 登录/注册返回：`accessToken`, `accessExpire`, `refreshToken`, `refreshExpire`
- 新增刷新接口：`POST /v1/user/token/refresh`
  - 入参：`refreshToken`
  - 出参：新 access + 新 refresh（旋转）
- 新增登出接口（推荐）：`POST /v1/user/logout`
  - 入参：`refreshToken`（或从 header/cookie 取）
  - 行为：移除白名单

---

## 4. 文件级改造清单（按优先级）

## P0: 认证核心与配置

1. `apps/user/rpc/internal/config/config.go`
- `Jwt` 配置扩展：
  - `AccessSecret`, `AccessExpire`（改为 900）
  - `RefreshSecret`, `RefreshExpire`（604800）

2. `apps/user/rpc/etc/user.yaml`
- 增加 `RefreshSecret` 与 `RefreshExpire`。
- 将 `AccessExpire` 调整为 `900`。

3. `apps/api/internal/config/config.go` + `apps/api/etc/api-api.yaml`
- 同步 Access 配置（用于 `rest.WithJwt`）。
- 是否保留 `Refresh` 配置取决于刷新逻辑放置层（建议刷新逻辑在 rpc 层，api 不需要 refresh secret）。

4. `pkg/auth/token.go`
- 从单函数升级为显式函数：
  - `GenerateAccessToken(...)`
  - `GenerateRefreshToken(...)`
  - `ParseRefreshToken(...)`
- 统一 claim key，解决 `pkg/auth` 与 `pkg/ctxdata` identify 常量不一致隐患。

5. `pkg/constants/redis.go`
- 新增 refresh 白名单 key 前缀常量。

## P1: RPC 契约与业务

6. `apps/user/rpc/user.proto`
- 修改 `LoginResp` / `RegisterResp` 为双 token 字段。
- 新增 `RefreshTokenReq/RefreshTokenResp` 与 `rpc RefreshToken(...)`。
- 可选新增 `LogoutReq/LogoutResp` 与 `rpc Logout(...)`。

7. 重新生成 user rpc 代码
- 更新 `apps/user/rpc/user/*.pb.go` 与 client 代码。

8. `apps/user/rpc/internal/logic/loginlogic.go`
- 登录成功后签发 access + refresh。
- refresh 写 Redis 白名单。

9. `apps/user/rpc/internal/logic/registerlogic.go`
- 注册后同样签发双 token 并写白名单。

10. 新增 `apps/user/rpc/internal/logic/refreshtokenlogic.go`
- 解析 refresh -> 校验 typ/exp/uid/jti。
- 查 Redis 白名单 + 哈希比对。
- 旋转 refresh 并签发新 access。

11. 可选新增 `apps/user/rpc/internal/logic/logoutlogic.go`
- 解析 refresh，删除对应白名单 key。

## P2: API 层对外路由

12. `apps/api/internal/types/types.go`
- 扩展 `LoginResp/RegisterResp` 字段。
- 新增 refresh/logout 请求响应结构。

13. `apps/api/internal/handler/routes.go`
- 新增 `/token/refresh` 与 `/logout` 路由。

14. `apps/api/internal/logic/user/loginlogic.go`
- 从 TODO 改为调用 user-rpc 的 Login。

15. `apps/api/internal/logic/user/registerlogic.go`
- 从 TODO 改为调用 user-rpc 的 Register。

16. 新增 `apps/api/internal/logic/user/refreshtokenlogic.go`（和 handler）
- 调 user-rpc 的 RefreshToken。

## P3: WS 与消费端协同

17. `apps/im/ws/internal/handler/auth.go`
- 明确仅接受 `typ=access` 的 token（防 refresh 被误用连接 WS）。

18. 客户端约定（文档）
- 登录成功后立即建立 WS，header/query 携带 access token。
- Access 过期后先 refresh，再重连/续连。

---

## 5. 关键安全策略

1. Refresh token 旋转（one-time or near one-time）
- 每次刷新都替换旧 refresh，降低重放窗口。

2. Redis 存哈希不存明文
- 泄露后无法直接伪造请求。

3. Token 类型隔离
- Access 与 Refresh 必须带 `typ`，服务端强校验用途。

4. 即时撤销
- 删除白名单即失效，无需等待 7 天过期。

5. 失败与风控（建议）
- 对 refresh 接口做限流（IP+uid 维度）。
- 异常刷新失败次数过高可触发临时冻结策略。

---

## 6. 性能与可用性评估

- 常规 API/WS 请求：仅验 Access JWT，不触 Redis，吞吐高。
- 仅 refresh/logout 命中 Redis：把状态成本集中在低频路径。
- Redis 故障退化策略：
  - refresh 暂时不可用（返回明确错误码），但已签发 access 在有效期内仍可工作。

---

## 7. 分阶段实施建议（推荐节奏）

阶段 A（最小可用）：
- P0 + P1 的 Login/Register 双 token 与 RefreshToken。
- 客户端先接入 refresh 流程。

阶段 B（安全增强）：
- Logout、全部设备下线、refresh 旋转严格 one-time、限流审计。

阶段 C（体验优化）：
- 统一错误码、补充 SDK/前端拦截器、完善观测指标（刷新成功率/撤销命中率）。

---

## 8. 验收清单（DoD）

- 登录后响应同时返回 access + refresh。
- access 在 15 分钟后过期并被 API/WS 拒绝。
- refresh 在 7 天内可换新，且每次换新后旧 refresh 立即失效。
- 手动登出后 refresh 立即不可用。
- 删除 Redis 白名单后会话即时失效。
- 现有受保护 API 路由仍可用，性能无明显回退。

---

## 9. 风险点与规避

- 风险：proto 字段改动导致旧客户端不兼容。
  - 规避：已确认保留原 `token/expire` 字段 1 个版本周期，新增双 token 字段并标注弃用计划。

- 风险：`pkg/auth` 与 `pkg/ctxdata` claim key 不一致造成解析失败。
  - 规避：统一 claim 常量并全链路替换。

- 风险：API 层当前 login/register 逻辑未完成，改造范围被低估。
  - 规避：先补通 api->rpc 基础链路，再叠加双 token 能力。
