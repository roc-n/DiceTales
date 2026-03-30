# 社区广场(Post)微服务需求文档 (纯 RPC 层)

## 1. 服务概述
- **服务名称**: Post Service (post-rpc)
- **核心职责**: 处理用户产生的内容(UGC)，包括发帖、互动(点赞/评论)以及时间流(Feed)的高性能聚合与分发。
- **技术定位**: **写多读多**的高吞吐写入与低延迟读取服务，作为简历中的**Feed 流架构设计**亮点展示。
- **涉及数据库表**: `post` (帖子本体), `comment` (评论), `social_interaction` (关系表:点赞/收藏), `user_feed` (Redis Only)

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
- **用途描述**: 核心发布接口，涉及**推模式(Push)写扩散**。
- **技术亮点描述**: 
  - **写扩散 (Push)**: 当用户发布帖子后，除写入 DB 外，需异步(通过 MQ)将 `PostId` 推送到其所有在线好友的 `Inbox` (Redis List/ZSet)。
  - **大V限流**: 加入逻辑判断，若该用户粉丝数 > 5000，则降级为只写入此用户的 `Outbox`，避免写放大导致系统雪崩。

### 2.2 GetTimeline (获取动态流 - **核心技术亮点**)
- **方法名**: `GetTimeline`
- **Req**: 
  - `UserId` (int64, 必填): 查看者ID
  - `LastPostId` (int64, 选填): 分页游标 (上一页最后一条ID)
  - `Limit` (int32, 必填): 每页数量
- **Resp**: `List<PostDetail>` (聚合了点赞数、评论预览、发布者信息)
- **用途描述**: 首页推荐流、好友关注流。
- **技术亮点描述**: **推拉结合 (Hybrid) 模式**:
  1. **拉取 Inbox**: 从 Redis `user:inbox:{uid}` 读取推送到该用户的帖子 ID (普通好友动态)。
  2. **拉取 Outbox**: 并发检查用户关注的“大V”列表，拉取他们的 `user:outbox:{vid}` (官方/KOL动态)。
  3. **归并排序**: 将上述两部分 ID 集合按时间戳做多路归并排序 (Merge Sort)，提取 Top N。
  4. **详情组装**: 批量查询 Post 详情及相关 User 信息 (利用 LocalCache 优化)。

### 2.3 Interact (互动操作)
- **方法名**: `Interact`
- **Req**: 
  - `UserId` (int64, 必填)
  - `TargetId` (int64, 必填): 目标ID (PostId / CommentId)
  - `TargetType` (enum, 必填): POST / COMMENT
  - `Action` (enum, 必填): LIKE / UNLIKE / SHARE / REPORT
- **Resp**: `Success` (bool)
- **用途描述**: 点赞、取消点赞等高频操作。
- **技术亮点描述**: **最终一致性与异步计数**:
  - 先写 Redis Hash `post:stats:{id}` (field: `likes` +1/-1)。
  - 异步发送 MQ 消息 `InteractionEvent`，由消费者负责写入 DB 落库，允许 DB 与 Redis 短时数据不一致，优先保证前端响应速度。

### 2.4 GetPostDetail (获取单条详情)
- **方法名**: `GetPostDetail`
- **Req**: `PostId` (int64)
- **Resp**: `PostDetail` (包含评论预览 Top 3)

## 3. 业务逻辑与规则补充

### 3.1 评论树结构
- **二级评论**: 为了简化设计与交互，仅支持两级评论体系 (评论 -> 回复)，不无限嵌套。
- **热门评论**: 在 Redis ZSet 中维护每个帖子的 Top 5 热门评论 (按点赞数排序)，随帖子详情一同返回。

### 3.2 敏感词过滤
- **内容安全**: 在 `CreatePost` 与 `AddComment` Logic 层，必须接入 `SensitiveWordFilter` 中间件，阻断违规内容发布。

## 4. 项目工程结构期望
- **RPC Proto**: `apps/post/rpc/post.proto`
- **Model**: `apps/post/model/` (含 `PostModel`, `CommentModel`, `SocialInteractionModel`)
- **Logic**: 
  - `apps/post/rpc/internal/logic/gettimelinelogic.go` (**重点实现**)
  - `apps/post/rpc/internal/logic/interactlogic.go`
- **SQL**: `deploy/sql/post.sql`