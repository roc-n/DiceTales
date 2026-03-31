# 社区广场(Post)微服务需求文档 (纯 RPC 层)

## 1. 服务概述
- **服务名称**: Post Service (post-rpc)
- **核心职责**: 处理用户产生的内容(UGC)，包括发帖、互动(点赞/评论)以及时间流(Feed)的高性能聚合与分发。
- **技术定位**: **写多读多**的高吞吐写入与低延迟读取服务，作为简历中的**Feed 流架构设计**亮点展示。
- **涉及数据库表**: 
  - `post` (帖子本体 - MySQL：强一致、低频写、多条件检索) 
  - `comment` (评论 - MongoDB：高频写、文档结构灵活)
  - `social_interaction` (关系记录表:点赞/收藏 - MongoDB：超高频写)
  - `user_feed` (Redis ZSet，用于极速拉取好友流)

## 2. RPC 接口设计 (gRPC)
> 以下定义了向内网其他微服务（如 BFF 层）提供的内部方法。

### 2.1 CreatePost (发布动态)
- **方法名**: `CreatePost`
- **Req**: 
  - `UserId` (int64, 必填): 发布者ID
  - `Content` (string, 必填): 文本内容
  - `Images` ([]string, 选填): 图片URL列表
  - `Visibility` (enum, 必填): PUBLIC(公开), FRIENDS(仅好友), PRIVATE(私密)
  - `RelatedGameId` (int64, 选填): 关联的游戏ID
- **Resp**: `PostId` (int64)
- **用途描述**: 核心发布接口，采用**纯推模式(Push)写扩散**架构。
- **技术亮点描述**: 
  - **纯写扩散 (Push)**: 当用户发布帖子后，除落库 DB 外，通过异步(MQ)将 `PostId` 直接推送到所有好友的 `Inbox` (Redis ZSet，Score为时间戳)。由于取消了复杂的粉丝机制而仅保持好友关系(数量有上限)，该模式避免了雪崩效应，实现极速读取。

### 2.2 GetFriendTimeline (获取好友动态流 - **核心技术亮点**)
- **方法名**: `GetFriendTimeline`
- **Req**: 
  - `UserId` (int64, 必填): 查看者ID
  - `LastPostId` / `Cursor` (int64, 选填): 分页游标 (如上一页最后一条的Score/时间戳)
  - `Limit` (int32, 必填): 每页数量
- **Resp**: `List<PostDetail>` (聚合了点赞数、评论预览、发布者信息)
- **用途描述**: 朋友圈页面，仅展示用户的好友动态。
- **技术亮点描述**: **即时拉取 (Inbox Pull)**:
  1. **读取收件箱**: O(1) 或 O(log n) 复杂度直接从 Redis `user:inbox:{uid}` 分页拉取数据，响应极快，无需复杂的运行时归并排序。
  2. **详情并发组装**: 针对提取出的 PostId 列表，并发读取 Redis + MySQL/Mongo 两级缓存即可完成帖子、评论、互动三类明细装配，无需额外的本地内存缓存组件，保持 V1 架构极简。

### 2.3 GetPublicTimeline (获取全站广场流 - **可演进性亮点**)
- **方法名**: `GetPublicTimeline`
- **Req**: 
  - `UserId` (int64, 必填): 查看者ID
  - `Cursor` (int64, 选填): 分页游标 (如最后一条帖子的创建时间或 ID)
  - `Limit` (int32, 必填): 每页数量
- **Resp**: `List<PostDetail>`
- **用途描述**: 全局公开大厅（System-wide Feed），首页默认的推荐广场，不局限于好友。
- **技术亮点描述**: 
  - **基于游标的深分页 (Cursor-based Pagination)**: V1 版本直接请求 MySQL，通过 `WHERE id < ? ORDER BY id DESC LIMIT ?` 或基于时间戳游标进行持续分页读取，避免了传统的 `OFFSET` 深分页性能衰减。
  - **简单倒序与预留推荐策略**: 初版无需额外经过 Redis 做聚合索引，直接打底层数据库即可满足性能与极简需求。但是我们将此接口提取为独立的 RPC 接口，为日后接入推荐系统做好了铺垫。当下游演进出（召回 -> 粗排 -> 精排）等个性化推荐模型时，客户端完全零入侵即可切换为千人千面的信息流。

### 2.4 Interact (互动操作)
- **方法名**: `Interact`
- **Req**: 
  - `UserId` (int64, 必填)
  - `TargetId` (int64, 必填): 目标ID (PostId / CommentId)
  - `TargetType` (enum, 必填): POST / COMMENT
  - `Action` (enum, 必填): LIKE_POST, UNLIKE_POST, COLLECT, UNCOLLECT 等
- **Resp**: `Success` (bool)
- **用途描述**: 处理点赞、收藏等对性能要求极高的零碎互动行为。
- **技术亮点描述**: **写合并与多级异步计数架构 (Write-Behind & Coalescing)**:
  1. **事实极速落库**: 请求到达后，直接向 MongoDB 的 `social_interaction` 集合执行高效的 upsert (依赖其极高的并发写能力)，成功后直接阻塞返回前端，保证极致流畅的用户体验。
  2. **异步事件分发**: 开启异步协程将此次变动包装为 `InteractionEvent` 发送至 Kafka，解除统计逻辑与主请求的耦合。
  3. **基于 Redis 聚合缓冲**: Consumer 从 MQ 消费消息，并将相对的增减量记录至 Redis Hash 中 (`agg:post:{id} -> {like: +1, collect: -1}`)。
  4. **局部批量落库**: 当 Redis 中的缓冲事件积攒到设定阈值 (如 100 次操作或达到 20 秒延迟)，利用原子命令提取缓冲，并执行一次合并的 SQL `UPDATE`，大幅削减 MySQL 的行锁竞争，抵御突发的高热点流量冲击。

### 2.5 CreateComment (发表评论)
- **方法名**: `CreateComment`
- **Req**:
  - `UserId` (int64, 必填)
  - `PostId` (int64, 必填)
  - `Content` (string, 必填)
  - `ParentId` (int64, 选填): 所属上级评论ID
  - `RootId` (int64, 选填): 根评论ID(用于嵌套折叠)
- **Resp**: `CommentId` (int64)
- **用途描述**: 为主推文或已有评论添加回复。
- **技术亮点描述**:
  - **隔离落库**: 结构灵活的评论本体被写入高吞吐的 MongoDB `comment` 集合中。
  - **复用计数通道**: 落库完成后，同样抛出一个特殊的 `InteractionEvent` 至 Kafka 聚合通道中，由上述同一套消费者系统在缓冲后批量累加 `post` 表的 `comment_count` 字段。
  - **敏感内容安全**: 在写入底层存储前，需强制调用接入的 `SensitiveWordFilter` 中间件完成敏感词阻断过滤。

### 2.6 GetPostDetail (获取单条详情)
- **方法名**: `GetPostDetail`
- **Req**: `PostId` (int64)
- **Resp**: `PostDetail` (包含评论预览 Top 3)

## 3. 业务逻辑与规则补充

### 3.1 评论树结构
- **二级评论**: 为了简化设计与交互，仅支持两级评论体系 (评论 -> 回复)，不无限嵌套。
- **热门评论**: 在 Redis ZSet 中维护每个帖子的 Top 5 热门评论 (按点赞数排序)，随帖子详情一同返回。

### 3.2 敏感词过滤
- **内容安全**: 在 `CreatePost` 与 `AddComment` Logic 层，必须接入 `SensitiveWordFilter` 中间件，阻断违规内容发布。

### 3.3 存储架构与选型说明
- **帖子本体 (MySQL)**: 帖子数据写入频率相对可控，且需要多维检索（按用户、时间、可见范围）；选择 MySQL 可以利用强一致事务与成熟索引，方便做分页与落地归档。聚合后计数的批量回写也能规避热点行锁问题。
- **评论 & 互动事实 (MongoDB)**: 评论、点赞、收藏属于高频写零碎数据，结构化且可变动大。利用 Mongo 高速并发写入特点保证操作响应时间天然横向扩展。
- **异步聚合缓冲与 Feed 索引 (Redis)**: 
  - `user:inbox:{uid}` 等 ZSet 负责高频的时间序列索引加载。
  - `agg:post:{id}` 等 Hash 对象负责承载来自 Kafka 的“多级缓存与写聚合(Write-Coalescing)”，极大降低对 MySQL 写计数值的 I/O 次数。

## 4. 项目工程结构期望
- **RPC Proto**: `apps/post/rpc/post.proto`
- **Model**: `apps/post/model/` (含 `PostModel`, `CommentModel`, `SocialInteractionModel`)
- **Logic**: 
  - `apps/post/rpc/internal/logic/gettimelinelogic.go` (**重点实现**)
  - `apps/post/rpc/internal/logic/interactlogic.go`
- **SQL**: `deploy/sql/post.sql`