# 即时通讯(IM)微服务需求文档 (RPC + 存储层)

## 1. 服务概述
- **服务名称**: IM Service (im-rpc)
- **核心职责**: 即时消息的持久化与时序管理，基于会话游标的消息漫游(History Messages)，离线未读消息同步(读扩散模型)。
- **技术定位**: 持载极高频的在线 IM 消息收发场景，解决数据写风暴，设计优雅的“基于游标的红点分离+读扩散同步架构”。
- **简历亮点对齐**: 
  - **雪花算法**: 用于确保海量群消息主键(MsgId)的全局唯一与粗略有序。
  - **双Buffer发号器**: 基于 Redis Set (或 Hash/List) 提前预取分配，实现高性能的递增 ID 发放（用于用户UID、群组GroupID 等业务实体的创建）。
- **涉及数据库表(MongoDB)**: `conversation` (用户拉群/单聊的会话状态及红点游标), `chat_log` (物理消息体，同一条群消息全网仅存一份)

---

## 2. MongoDB 数据模型设计
由于 MongoDB 不需要预先通过 sql 建表，我们将在此直接定义好 Schema 的核心字段（后续将通过 `goctl model mongo` 生成并补充对应的 BSON 标签）。

### 2.1 会话表 (Conversation)
表示**某个用户**视角下的某个聊天窗口的状态，所以每一个参与者都会有一条独立的记录（便于各自管理未读数）。
```json
{
  "_id": "ObjectId",           // Mongo自增主键 (或可选自己填入ID)
  "conversationId": "string",  // 单聊如: "u1_u2" (按字典序拼接)，群聊如: "g_15234" (这里群ID是由双Buffer发号器生成)
  "chatType": "int8",          // 枚举: 1-单聊 2-群聊
  "ownerId": "string",         // 当前会话状态的归属用户ID (基于双Buffer发号器生成)
  "targetId": "string",        // 聊天对象ID / 对象群ID
  "readSeq": "int64",          // 用户已读/已同步拉取完毕的最后一条消息 Seq 游标
  "maxSeq": "int64",           // 会话当前的最新一条 Seq (通常由 Redis 缓存维护，定期刷盘到此字段)
  // 【冗余辅助展示字段】
  "lastMsgId": "string",       // 预留: 用于 UI 加载时免查 ChatLog 即可渲染最后一条消息摘要
  "lastMsgTime": "int64",      // 极重要: 用于会话列表页按时间降序排列
  "isShow": "bool",            // 是否展示在用户的会话列表页中 (例如清空了聊天记录但未退群时可设为 false)
  "createAt": "int64",
  "updateAt": "int64"
}
```
> **预估关键索引**: 
> 1. `uid_time_idx`: `{ownerId: 1, lastMsgTime: -1}` (用于首页消息列表的高效拉取)
> 2. `uid_conv_idx`: `{ownerId: 1, conversationId: 1}` (唯一索引，防并发创建重复会话)

### 2.2 消息流水表 (ChatLog)
存储绝对的消息实体，采用**读扩散**，无论群里有多少人，物理消息只存储一条。
```json
{
  "_id": "string",             // 采用雪花算法(Snowflake)生成的全局唯一 MsgId
  "conversationId": "string",  // 归属的会话 ID
  "seq": "int64",              // 该会话内部**严格单调递增**的序号 (依赖 Redis 原子 INCR 生成)
  "sendId": "string",          // 发送人
  "recvId": "string",          // 接收人 (仅单聊有，群聊此字段为空即可)
  "msgType": "int8",           // 消息格式: 1-文本 2-图片 3-语音 4-系统通知
  "msgContent": "string",      // 消息正文主体 (可考虑后续加解密包装)
  "sendTime": "int64"
}
```
> **预估关键索引**: 
> 1. `conv_seq_idx`: `{conversationId: 1, seq: 1}` (用于离线游标递增拉取: `where seq > my_read_seq`)

---

## 3. RPC 接口设计 (gRPC)

结合我们此前的技术推演方案，在 `im.proto` 中应该暴露以下核心内网接口（通常由 `im.ws` 或 `bff` 网关调用）。

### 3.1 `SetUpUserConversation` (建立/唤醒会话)
- **Req**: `SenderId`, `ReceiverId`, `ChatType` (发信人、收信人/群、类型)
- **Resp**: `ConversationId`
- **用途叙述**: 发送第一条消息或者初次添加好友/进群时的前置信令。若会话已存在，重置 `isShow = true` 并直接返回对应的 `ConversationId`，保障双端通讯通道在 DB 中就绪。

### 3.2 `PutChatLog` (持久化聊天消息 - 内存缓冲核心)
- **Req**: `ConversationId`, `SendId`, `RecvId`, `ChatType`, `MsgType`, `MsgContent`
- **Resp**: `MsgId`, `Seq`, `SendTime`
- **用途叙述**: 承接底层网关/消费者投递过来的业务消息。
- **技术亮点流转**: 
  1. 通过**雪花算法** (Snowflake) 直接在运行时为当前记录生成一个全局唯一的 `MsgId` (字符串或 int64)，以规避 Mongo 主键生成的局限性。
  2. 通过 Redis 的 `INCRBY conv_seq:{conversationId} 1` 直接生成全局单调连续的 `Seq`。
  3. 把这最新的 `Seq` 以 `HSET` 旁路暂存进 `user_max_seq:{owner}` Hash 中（不直接写 Mongo，避免写风暴）。
  4. 使用前 3 步生成的变量，执行 Mongo `ChatLog` Collection 的 `InsertOne` 落盘操作。

### 3.3 `GetConversations` (拉取首页联系人列表)
- **Req**: `OwnerId` (当前发起请求的用户)
- **Resp**: `List<Conversation>` (包含核心渲染字段及 `UnreadCount`)
- **用途叙述**: 客户端打开 APP 时的高频调用请求，拉取最近的活跃会话红点。
- **技术亮点流转**: 服务端通过 `ownerId` 从 Mongo 的会话表中根据 `lastMsgTime` 倒序拉出近期的列表，并将列表中的 `readSeq` 与 暂存的 `maxSeq` 进行无运算 (`UnreadCount = maxSeq - readSeq`) 返回，全程极低开销。

### 3.4 `GetChatLog` (基于游标分页同步离线消息)
- **Req**: `ConversationId`, `CursorSeq` (客户端目前持有的最后一条 Seq 的数字), `Limit` (比如 50 条)
- **Resp**: `List<ChatLog>`
- **用途叙述**: 典型的“读扩散漫游”接口。点进某个带有红点的群聊时，客户端利用自身游标进行“向上追溯接轨”。
- **技术亮点流转**: 将请求被直接翻译为 Mongo Query: `Find({conversationId: id, seq: {$gt: CursorSeq}}).Sort(seq:1).Limit(50)`。这能极大降低拉群记录时的网络负载。

### 3.5 `AckMessageRead` (消息已读推游标)
- **Req**: `OwnerId`, `ConversationId`, `ReadSeq`
- **Resp**: `Success` (bool)
- **用途叙述**: 客户端确认自己已经在屏幕上渲染过了最新那批消息。
- **技术亮点流转**: 典型的 **Tick & Merge 异步延时刷库防写风暴**。接收到的大量 ACK 不直接 `Update` MongoDB 的 `conversation` 表。全部先写入 Redis，随后通过定时任务使用 MongoDB 的 BulkWrite 进行延迟覆盖刷盘。

---

## 4. 后续项目工程结构期望指南
等确认以上架构无误后，基于当前的 Go-Zero 开发习惯，我们的下一步行动方向将是：
1. 编辑 `apps/im/rpc/im.proto` 完成类型和接口约束。
2. 使用 `goctl rpc protoc` 跑通整个 `im/rpc` 初始化脚手架。
3. 使用 `goctl model mongo` 一键生成对应 `conversation` 与 `chatlog` 的增删改查 Mongo Driver 模板。
4. 按照上述**技术亮点**（如 Redis 旁路缓存与 Seq 生成器）逐步向 Logic 业务层填充实体代码。