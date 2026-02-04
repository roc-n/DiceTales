package logic

import (
	"context"
	"database/sql"
	"time"

	"dicetales.com/apps/social/model"
	"dicetales.com/apps/social/rpc/internal/svc"
	"dicetales.com/apps/social/rpc/social"
	"dicetales.com/pkg/errorx"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type GroupPutinLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupPutinLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupPutinLogic {
	return &GroupPutinLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GroupPutinLogic) GroupPutin(in *social.GroupPutinReq) (*social.GroupPutinResp, error) {

	//  普通用户申请： 群无验证直接进入
	//  群成员邀请：   群无验证直接进入
	//  群管理员/群创建者邀请：直接进群
	var (
		inviteGroupMember *model.GroupMember
		userGroupMember   *model.GroupMember
		groupInfo         *model.Groups

		err error
	)

	// 查询该用户是否已在群里
	userGroupMember, err = l.svcCtx.GroupMemberModel.FindByGroudIdAndUserId(l.ctx, in.ReqUid, in.GroupId)
	if err != nil && err != model.ErrNotFound {
		return nil, errors.Wrapf(errorx.NewDBErr(), "find group member by groud id and  req id err [%v], req [%v], [%v]", err, in.GroupId, in.ReqUid)
	} else if userGroupMember != nil {
		return &social.GroupPutinResp{}, nil // 已在群里，直接返回
	}

	// 查询该用户是否已申请过该群
	groupReq, err := l.svcCtx.GroupRequestModel.FindByGroupIdAndReqId(l.ctx, in.GroupId, in.ReqUid)
	if err != nil && err != model.ErrNotFound {
		return nil, errors.Wrapf(errorx.NewDBErr(), "find group req by groud id and user id err [%v], req [%v], [%v]", err,
			in.GroupId, in.ReqUid)
	}
	if groupReq != nil {
		return &social.GroupPutinResp{}, nil // 已申请过，直接返回
	}

	groupReq = &model.GroupRequest{
		ReqUid:  in.ReqUid,
		GroupId: in.GroupId,
		ReqMsg: sql.NullString{
			String: in.ReqMsg,
			Valid:  true,
		},
		ReqTime: sql.NullTime{
			Time:  time.Unix(in.ReqTime, 0),
			Valid: true,
		},
		JoinSource: sql.NullInt64{
			Int64: int64(in.JoinSource),
			Valid: true,
		},
		InviterUid: sql.NullString{
			String: in.InviterUid,
			Valid:  true,
		},
		HandleResult: sql.NullInt64{
			Int64: int64(NoHandlerResult),
			Valid: true,
		},
	}

	// 查询群是否需要验证
	groupInfo, err = l.svcCtx.GroupsModel.FindOne(l.ctx, in.GroupId)
	if err != nil {
		return nil, errors.Wrapf(errorx.NewDBErr(), "find group by groud id err [%v], req [%v]", err, in.GroupId)
	}
	if !groupInfo.IsVerify {
		defer l.createGroupMember(in)

		groupReq.HandleResult = sql.NullInt64{
			Int64: int64(PassHandlerResult),
			Valid: true,
		}
		return l.createGroupReq(groupReq, true)
	}

	// 需要验证，看看是申请进群还是邀请进群
	if GroupJoinSource(in.JoinSource) == PutInGroupJoinSource {
		return l.createGroupReq(groupReq, false)
	} else {
		// 查询邀请人相关信息
		inviteGroupMember, err = l.svcCtx.GroupMemberModel.FindByGroudIdAndUserId(l.ctx, in.InviterUid, in.GroupId)
		if err != nil {
			return nil, errors.Wrapf(errorx.NewDBErr(), "find group member by groud id and user id err [%v], req [%v]", in.InviterUid, in.GroupId)
		}

		// 邀请人是群主或管理员，直接通过
		if GroupRoleLevel(inviteGroupMember.RoleLevel) == CreatorGroupRoleLevel || GroupRoleLevel(inviteGroupMember.RoleLevel) == ManagerGroupRoleLevel {
			defer l.createGroupMember(in)

			groupReq.HandleResult = sql.NullInt64{
				Int64: int64(PassHandlerResult),
				Valid: true,
			}
			return l.createGroupReq(groupReq, true)
		}

		// 邀请人是普通成员，等待处理
		return l.createGroupReq(groupReq, false)
	}
}

// 添加群申请记录
func (l *GroupPutinLogic) createGroupReq(groupReq *model.GroupRequest, isPass bool) (*social.GroupPutinResp, error) {
	_, err := l.svcCtx.GroupRequestModel.Insert(l.ctx, groupReq)
	if err != nil {
		return nil, errors.Wrapf(errorx.NewDBErr(), "insert group req err [%v] req [%v]", err, groupReq)
	}

	if isPass {
		return &social.GroupPutinResp{GroupId: groupReq.GroupId}, nil
	}

	return &social.GroupPutinResp{}, nil
}

// 添加群成员
func (l *GroupPutinLogic) createGroupMember(in *social.GroupPutinReq) error {
	groupMember := &model.GroupMember{
		GroupId:   in.GroupId,
		Uid:       in.ReqUid,
		RoleLevel: int64(MemberGroupRoleLevel),
		InviterUid: sql.NullString{
			String: in.InviterUid,
			Valid:  true,
		},
		JoinTime:   sql.NullTime{Time: time.Unix(in.ReqTime, 0), Valid: true},
		JoinSource: sql.NullInt64{Int64: int64(in.JoinSource), Valid: true},
	}
	_, err := l.svcCtx.GroupMemberModel.Insert(l.ctx, groupMember)
	if err != nil {
		return errors.Wrapf(errorx.NewDBErr(), "insert group member err [%v], req [%v]", err, groupMember)
	}

	return nil
}
