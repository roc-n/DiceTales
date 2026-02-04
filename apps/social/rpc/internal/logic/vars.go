package logic

// 处理结果
// 1 未处理，2-3 处理， 4 取消
type HandlerResult int

const (
	NoHandlerResult     HandlerResult = iota + 1 // 未处理
	PassHandlerResult                            // 通过
	RefuseHandlerResult                          // 拒绝
	CancelHandlerResult                          // 取消
)

// 好友申请来源
type FriendAddSource int

const (
	RequestJoinSource FriendAddSource = iota + 1
	RecommendJoinSource
)

// 群成员级别
// 1 群主，2 管理者，3 普通成员
type GroupRoleLevel int

const (
	CreatorGroupRoleLevel GroupRoleLevel = iota + 1
	ManagerGroupRoleLevel
	MemberGroupRoleLevel
)

// 进群申请的方式
// 1 邀请， 2 申请
type GroupJoinSource int

const (
	InviteGroupJoinSource GroupJoinSource = iota + 1
	PutInGroupJoinSource
)

type GroupType int

const (
	NormalGroupType GroupType = iota + 1
	MeetupGroupType
)
