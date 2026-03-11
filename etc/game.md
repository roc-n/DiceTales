# 桌游(Game)微服务需求文档 (纯 RPC 层)

## 1. 服务概述
- **服务名称**: Game Service (game-rpc)
- **核心职责**: 作为底层领域服务，负责桌游核心数据的 CRUD 以及扩展分类信息的管理，不直接对外暴露 HTTP 接口。
- **涉及数据库表**: `game`, `game_category_info` (基于 deploy/sql/game.sql)

## 2. RPC 接口设计 (gRPC)
> 以下定义了向内网其他微服务（如 BFF 层、social-rpc 等）提供的内部方法。

### 2.1 GetGameDetail (获取桌游详细信息)
- **方法名**: `GetGameDetail`
- **Req**: 
  - `GameId` (int64, 必填): 游戏ID
- **Resp**: 
  - 游戏基础信息: `Id`, `Name`, `NameEn`, `CoverImg`, `Score`, `ScoreCount`, `MinPlayers`, `MaxPlayers` 等
  - 游戏分类信息: `Category`, `Mode`, `Theme`, `Mechanic`, `SuitableAge` 等
- **用途描述**: 供 BFF 层组装桌游详情页数据。

### 2.2 ListGames (分页查询桌游列表)
- **方法名**: `ListGames`
- **Req**: 
  - `Page` (int32, 必填): 当前页码
  - `Size` (int32, 必填): 每页数量
  - `Keyword` (string, 选填): 搜索关键词 (匹配中文或英文名)
  - `Category` (string, 选填): 游戏类别筛选
- **Resp**: 
  - `Total` (int64): 总记录数
  - `List` (array): 包含简化版桌游信息的列表 (例如不含长文本 description)
- **用途描述**: 供 BFF 层提供给前端的桌游瀑布流、搜索页使用。

### 2.3 GetGameBasicInfo (获取桌游基础摘要)
- **方法名**: `GetGameBasicInfo`
- **Req**: 
  - `GameId` (int64, 必填)
- **Resp**: 
  - `Id` (int64), `Name` (string), `CoverImg` (string), `Score` (double)
- **用途描述**: 轻量级接口，主要供社交服务 (social-rpc) 在发帖子、组局或评论时，快速冗余/反查桌游基础信息。

## 3. 业务逻辑与规则补充
> 编写 Logic 代码时的核心指南。

- **异常处理**: 查询不到游戏时，返回 `xerr` 规范的统一错误码 (如 `ErrGameNotFound`)。
- **事务处理**: 未来如果增加 `CreateGame` 接口，涉及往 `game` 和 `game_category_info` 同时插入数据时，必须在同一 DB 事务中。
- **缓存策略**: 对于 `GetGameBasicInfo` 这种高频且变更极少的调用，可以在 logic 层结合 Redis 缓存，过期时间设定为 1 小时。
- **测试要求**: 逻辑层 (Logic) 的核心方法需提供单元测试，不连接真实 DB，使用 sqlmock / gomock 进行拦截。

## 4. 项目工程结构期望
- RPC Proto 位置: `apps/game/rpc/game.proto`
- Model 生成位置: `apps/game/model/`
- RPC 配置文件位置: `apps/game/rpc/etc/game.yaml` (监听固定端口或使用 etcd 注册)