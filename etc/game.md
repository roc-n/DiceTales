# 桌游(Game)微服务需求文档 (纯 RPC 层)

## 1. 服务概述
- **服务名称**: Game Service (game-rpc)
- **核心职责**: 作为全平台的“数据底座”与“流量入口”，负责桌游元数据的权威管理、高性能多维度检索与智能推荐。
- **技术定位**: **读多写少**的高并发查询服务，作为简历中的**复杂检索与缓存架构**亮点展示。
- **涉及数据库表**: `game`, `game_params` (扩展参数), `game_tag_relation` (标签关联), `game_category_info` (基于 deploy/sql/game.sql)

## 2. RPC 接口设计 (gRPC)
> 以下定义了向内网其他微服务（如 BFF 层、social-rpc、meetup-rpc 等）提供的内部方法。

### 2.1 GetGameDetail (获取桌游详细信息)
- **方法名**: `GetGameDetail`
- **Req**: 
  - `GameId` (int64, 必填): 游戏ID
- **Resp**: 
  - `BaseInfo`: `Id`, `Name` (中), `NameEn`, `CoverImg`, `Score`
  - `Params`: `MinPlayers`, `MaxPlayers`, `DurationLower` (最少时长), `DurationUpper`, `Complexity` (重度), `PublicationYear`
  - `Metadata`: `Mechanisms` (机制列表), `Themes` (主题列表), `Designers`
  - `Resources`: `RulebookUrl` (规则书), `GuideVideoUrl` (教学视频)
- **用途描述**: 供 BFF 层组装桌游详情页数据。需支持高并发读取。

### 2.2 FilterGames (高级筛选 - **核心技术亮点**)
- **方法名**: `FilterGames`
- **Req**: 
  - `PlayerCount` (int32, 选填): 例如 4 (筛选支持4人的游戏)
  - `DurationRange` (enum, 选填): 例如 SHORT(<30min), MEDIUM(30-90min), LONG(>90min)
  - `Complexity` (enum, 选填): 例如 LIGHT, MEDIUM, HEAVY
  - `TagIds` ([]int64, 选填): 机制或主题ID列表 (AND 关系)
  - `Page`, `Size`
- **Resp**: 
  - `Total` (int64)
  - `GameIds` ([]int64): 仅返回ID列表，详情需二次查询或批量获取
- **技术亮点描述**: 
  - **基于 Redis Bitmap 的位图索引**: 
    - 传统 SQL `WHERE` 多条件组合查询在数据量大时性能极差。
    - 这里的实现是将筛选条件映射为 Redis Bitmap (如 `idx:game:players:4`, `idx:game:time:short`, `idx:game:tag:101`)。
    - 通过 `BITOP AND` 运算直接得出满足所有条件的 GameID 集合，实现 O(1) 复杂度的极速筛选。

### 2.3 SearchGames (模糊搜索)
- **方法名**: `SearchGames`
- **Req**: 
  - `Keyword` (string, 必填): 搜索关键词 (支持中文、英文、拼音首字母)
  - `Page`, `Size`
- **Resp**: `List` (包含高亮信息的桌游列表)
- **用途描述**: 解决用户“只记得发音”或“只记得部分名字”的搜索需求。

### 2.4 GetGameBasicInfo (轻量级摘要)
- **方法名**: `GetGameBasicInfo`
- **Req**: `GameIds` ([]int64)
- **Resp**: Map<int64, GameBasic>
- **用途描述**: 供 Social/Meetup 服务批量反查游戏名称和封面 (如“用户A发布了关于《卡卡颂》的帖子”)。

## 3. 业务逻辑与规则补充 (Logic 层规范)

### 3.1 缓存一致性与防穿透
- **多级缓存**: 
  - L1: 本地缓存 (BigCache) 存储 Top 100 热门游戏详情。
  - L2: Redis 缓存 `cache:game:id:{id}` (ProtoBuf 序列化)。
- **防穿透**: 
  - 查询不存在的 ID 时，需缓存空对象 (Empty Object) 并设置较短 TTL (如 5 分钟)。
- **Bitmap 更新机制**:
  - 当运营后台新增或修改桌游属性时，需发送 MQ 消息 `GameAttributeChanged`。
  - Game 服务消费消息，异步更新对应的 Redis Bitmap 索引位，保证最终一致性。

### 3.2 评分计算策略
- **异步计算**: 
  - 用户的评分操作写入 MySQL 后，不立即重新计算平均分。
  - 通过定时任务 (Cron) 或 累计 N 次评分后触发一次重算，更新 `game` 表的 `score` 字段，避免热点行锁竞争。

## 4. 项目工程结构期望
- **RPC Proto**: `apps/game/rpc/game.proto`
- **Model**: `apps/game/model/` (包含 `GameModel`, `GameBitmapModel`)
- **Logic**: 
  - `apps/game/rpc/internal/logic/filtergameslogic.go` (**重点实现**)
  - `apps/game/rpc/internal/logic/getgamedetaillogic.go`
- **SQL**: `deploy/sql/game.sql` 需补充索引字段。
- RPC 配置文件位置: `apps/game/rpc/etc/game.yaml` (监听固定端口或使用 etcd 注册)