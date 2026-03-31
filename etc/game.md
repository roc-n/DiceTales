# 桌游(Game)微服务需求文档 (纯 RPC 层)

## 1. 服务概述
- **服务名称**: Game Service (game-rpc)
- **核心职责**: 作为全平台的“数据底座”与“流量入口”，负责桌游元数据的权威管理、高性能多维度检索与智能推荐。
- **技术定位**: **读多写少**的高并发查询服务，作为简历中的**复杂检索与缓存架构**亮点展示。
- **涉及数据库表**: 收敛对齐了最新数据模型，仅包含 `game` (主表), `game_resource` (资源表), `tag` (标签元数据), `game_tag_relation` (游戏与标签关联)。

## 2. RPC 接口设计 (gRPC)
> 以下定义了向内网其他微服务（如 BFF 层、social-rpc、meetup-rpc 等）提供的内部方法。

### 2.1 GetGameDetail (获取桌游详细信息)
- **方法名**: `GetGameDetail`
- **Req**: 
  - `GameId` (int64, 必填): 游戏ID
- **Resp**: 
  - `BaseInfo`: `Id`, `Name` (中), `NameEn`, `CoverImg`, `Score`
  - `Params`: `MinPlayers`, `MaxPlayers`, `DurationMin`, `DurationMax`, `Complexity`, `PublishYear`
  - `Metadata`: `Tags` (由机制、主题等合并)
  - `Resources`: 规则书、教学视频等关联表数据
- **用途描述**: 供 BFF 层组装桌游详情页数据。需支持高并发读取。

### 2.2 FilterGames (高级筛选 - **核心技术亮点**)
- **方法名**: `FilterGames`
- **Req**: 
  - `PlayerCount` (int32, 选填): 例如 4 (筛选支持4人的游戏)
  - `DurationRange` (enum, 选填): 例如 SHORT(<30min), MEDIUM(30-90min), LONG(>90min)
  - `ComplexityLevel` (enum, 选填): 例如 LIGHT(1-2), MEDIUM(2-3), HEAVY(3-5)
  - `TagIds` ([]int64, 选填): 机制或主题ID列表 (AND 关系)
  - `SortBy`: 排序字段 (如 "score", "hot")
  - `Page`, `Size`
- **Resp**: 
  - `Total` (int64)
  - `GameIds` ([]int64): 仅返回核心ID列表，供批量获取摘要获取渲染
- **技术亮点描述**: 
  1. **数据离散化(分桶)写入**: 将连续数值如玩家人数、游戏时长、重度划分为多个离散特征。游戏落库时将其 ID 同步标记到对应的 Redis Bitmap 分桶中。
  2. **基于 Redis Bitmap 的位图索引**: 将上述筛选条件直接映射为对应的 Redis Bitmap 键(如 `idx:game:player:4`, `idx:game:duration:short`, `idx:game:tag:101`)。通过 `BITOP AND` 直接在 Redis 中进行 O(1) 复杂度的按位与运算，得出所有满足条件的目标集合。
  3. **内存级联合排序与分页**: 将第一步计算得到的 Game ID 集合通过应用层拉回 Go 程序中。利用 `HMGET` 从 Redis 缓存批量查出对应的可算分值 (预计算 Score 或 Hot 热度)，利用 Go 极高的 CPU 优势**在内存中完成交集结果的最终排布与截断(Pagination)**。替代并秒杀传统的 SQL 多表 `Join` + `ORDER BY` 操作。

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

### 3.1 极简 L2 缓存架构与防穿透
- **抛弃 L1 本机缓存**: 摒除复杂且易带来数据脑裂的单机进程缓存，全面下沉到 L2 Redis 缓存层，兼顾极致简单与分布式强一致性。
- **缓存应用**: `game` 基本信息与详情使用 `cache:game:id:{id}` 进行存储 (基于 ProtoBuf 高效序列化)。
- **防穿透/防雪崩**: 
  - 查询不存在的 ID 时，需缓存空对象 (Empty Object) 并设置短 TTL (如 5 分钟)。
  - 对合法数据增加过期时间“随机抖动 (Random Jitter)”，防缓存雪崩引发 DB 宕机。

### 3.2 筛选引擎维护与数据预热
- **全量重建 (Cold Start)**: 提供 `/internal/rebuild_bitmap` 特权接口/脚本。支持全量遍历 MySQL `game` 与 `game_tag_relation` 表数据，将所有维度映射关系全部重构写入 Redis Bitmap，用于新系统上线或灾备恢复。
- **增量更新**: 当运营后台新增或修改桌游属性时，通过发送 MQ 消息 `GameAttributeChanged` 触发事件。Game 服务消费此消息后仅仅局部操作某个 Game ID 在特定 Bitmap 下的 `SETBIT 0/1`，在保证主业务低延迟的同时实现检索数据的最终一致。

### 3.3 评分计算策略
- **异步计算**: 
  - 用户的评分操作写入 MySQL 后，不立即重新计算平均分。
  - 通过定时任务 (Cron) 或 累计 N 次评分后触发一次重算，统一更新 `game` 表的 `score` 字段以及 Redis 哈希表中的排序依据值，避免热点行锁的高频竞争。

## 4. 项目工程结构期望
- **RPC Proto**: `apps/game/rpc/game.proto`
- **Model**: `apps/game/model/` (基于目前的 `deploy/sql/game.sql` 极简重构生成)
- **Logic**: 
  - `apps/game/rpc/internal/logic/filtergameslogic.go` (**核心重难点：位运算联合及内存排序分页**)
  - `apps/game/rpc/internal/logic/getgamedetaillogic.go`
- **SQL**: `deploy/sql/game.sql`。
- RPC 配置文件位置: `apps/game/rpc/etc/game.yaml` (监听固定端口或使用 etcd 注册)