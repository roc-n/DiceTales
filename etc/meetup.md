# 组局(Meetup)微服务需求文档 (纯 RPC 层)

## 1. 服务概述
- **服务名称**: Meetup Service (meetup-rpc)
- **核心职责**: 处理线下组局的全生命周期管理，包括发布、报名、审核、成局与结算。
- **技术定位**: **强一致性**与**LBS 位置服务**结合，作为简历中的**高并发库存扣减(秒杀)**与**地理位置检索**亮点展示。
- **涉及数据库表**: `meetup` (局信息), `meetup_member` (成员状态), `meetup_location` (位置信息)

## 2. RPC 接口设计 (gRPC)
> 以下定义了向内网其他微服务（如 BFF 层）提供的内部方法。

### 2.1 CreateMeetup (发起组局)
- **方法名**: `CreateMeetup`
- **Req**: 
  - `OrganizerId` (int64)
  - `GameId` (int64)
  - `Location`: `Address`, `Latitude` (double), `Longitude` (double)
  - `Time`: `StartTime`, `EndTime`
  - `Config`: `MaxPlayers`, `EnableAudit` (是否审核)
- **Resp**: `MeetupId`
- **业务逻辑**:
  - 写入 `meetup` 表。
  - **LBS 索引**: 将 `MeetupId` 写入 Redis GeoHash Set (`geo:meetups`)。
  - **库存预热**: 初始化 Redis 库存 Key `meetup:slots:{id}` = `MaxPlayers - 1` (发起人自动占位)。

### 2.2 JoinMeetup (报名/加入 - **核心技术亮点**)
- **方法名**: `JoinMeetup`
- **Req**: `MeetupId`, `UserId`
- **Resp**: `Status` (JOINED / PENDING / FULL)
- **技术亮点描述**: **Redis Lua 脚本实现无锁原子扣减**:
  - **场景痛点**: 热门车队（如仅剩1个位置）在并发下极易出现超卖 (4/3人)。
  - **核心方案**: 不使用 DB 悲观锁。直接执行 Lua 脚本：
    ```lua
    if redis.call("get", KEYS[1]) > 0 then
        redis.call("decr", KEYS[1])
        return 1
    else
        return 0
    end
    ```
  - **流程**: Lua 扣减成功 -> 发送 MQ `JoinSuccess` -> 消费者异步写入 MySQL `meetup_member`。若 Lua 返回 0，直接拒绝。

### 2.3 ListNearbyMeetups (LBS 发现 - **核心技术亮点**)
- **方法名**: `ListNearbyMeetups`
- **Req**: `MyLat`, `MyLon`, `RadiusKM` (半径)
- **Resp**: `List<MeetupBasic>`
- **技术亮点描述**:
  - 利用 Redis `GEORADIUS` 指令快速召回附近的 MeetupId。
  - 结合 `MeetupId` 批量查询详情 (若数据量大，可结合 GeoHash 编码前缀匹配优化)。

### 2.4 ManageMeetup (管理与审核)
- **方法名**: `ManageMeetup`
- **Req**: `MeetupId`, `OperatorId`, `TargetUserId`, `Action` (APPROVE/REJECT/KICK/CANCEL)
- **业务逻辑**:
  - **状态机**: 严格控制状态流转 (e.g., 只有 PENDING 状态可被 APPROVE)。
  - **库存回补**: 若踢出成员或成员退出，需通过 Lua 脚本原子增加 Redis 库存 `INCR`。

## 3. 业务逻辑与规则补充
- **过期处理**: 
  - 通过 Redis Key 过期通知 (Keyspace Notification) 或延迟队列 (RabbitMQ/Redisson)，在活动开始时间后自动标记为 `EXPIRED` (如未成局)。
- **签到逻辑移除**: 
  - 原定“签到”功能已移除。改为由发起人手动“结束活动”或系统自动根据时间推移标记为“已结束”。

## 4. 项目工程结构期望
- **RPC Proto**: `apps/meetup/rpc/meetup.proto`
- **Logic**: 
  - `apps/meetup/rpc/internal/logic/joinmeetuplogic.go` (**重点实现 Lua 调用**)
  - `apps/meetup/rpc/internal/logic/listnearbymeetuplogic.go`
- **SQL**: `deploy/sql/meetup.sql`
  - `OPEN` (报名中) -> `FULL` (满员, 自动触发)
  - `OPEN/FULL` -> `EXPIRED` (过期, 定时任务触发)
  - `OPEN` -> `CANCELLED` (发起人取消)
- **说明**: 移除了“签到”环节。成局后，若活动时间已过，系统自动标记为 `FINISHED` (结束)，或由发起人手动结束。

## 3. 数据库设计建议
- `meetup`: 包含 `status` (tinyint), `geo_hash` (varchar 索引)。
- `meetup_member`: 包含 `role` (Organizer/Member), `status` (Pending/Joined/Rejected)。