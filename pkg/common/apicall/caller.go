package apicall

import (
	"context"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/auth"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/friend"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/group"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/user"
	"github.com/OpenIMSDK/chat/pkg/common/config"
)

type CallerInterface interface {
	ImportFriend(ctx context.Context, ownerUserID string, friendUserID []string) error
	UserToken(ctx context.Context, userID string, platformID int32) (string, error)
	InviteToGroup(ctx context.Context, userID string, groupID string) error
	UpdateUserInfo(ctx context.Context, userID string, nickName string, faceURL string) error
	ForceOffLine(ctx context.Context, userID string) error
	RegisterUser(ctx context.Context, users []*sdkws.UserInfo) error
	FindGroup(ctx context.Context, groupIDs []string) ([]*sdkws.GroupInfo, error)
	MapGroup(ctx context.Context, groupIDss []string) (map[string]*sdkws.GroupInfo, error)
}

type Caller struct {
}

func (c *Caller) ForceOffLine(ctx context.Context, userID string) error {
	for id := range constant.PlatformID2Name {
		forceOffLine := NewApiCaller[auth.ForceLogoutReq, auth.ForceLogoutResp]("/auth/force_logout")
		_, err := forceOffLine.Call(ctx, &auth.ForceLogoutReq{
			PlatformID: int32(id),
			UserID:     userID,
		})
		if err != nil {
			log.ZError(ctx, "ForceOffline", err, "userID", userID, "platformID", id)
			return err
		}
	}
	return nil
}

func NewCallerInterface() CallerInterface {
	return &Caller{}
}

func (c *Caller) ImportFriend(ctx context.Context, ownerUserID string, friendUserID []string) error {
	importFriend := NewApiCaller[friend.ImportFriendReq, friend.ImportFriendResp]("/friend/import_friend")
	_, err := importFriend.Call(ctx, &friend.ImportFriendReq{
		OwnerUserID:   ownerUserID,
		FriendUserIDs: friendUserID,
	})
	if err != nil {
		log.ZError(ctx, "ImportFriend", err, "ownerUserID", ownerUserID)
		return err
	}
	return nil
}

func (c *Caller) UserToken(ctx context.Context, userID string, platformID int32) (string, error) {
	userToken := NewApiCaller[auth.UserTokenReq, auth.UserTokenResp]("/auth/user_token")
	resp, err := userToken.Call(ctx, &auth.UserTokenReq{
		Secret:     *config.Config.Secret,
		PlatformID: platformID,
		UserID:     userID,
	})
	log.ZDebug(ctx, "userToken", "userID", userID, "platform", platformID, "resp", resp, "err", err)
	if err != nil {
		log.ZError(ctx, "userToken", err, "userID", userID, "platform", platformID)
		return "", err
	}
	return resp.Token, nil
}

func (c *Caller) FindGroup(ctx context.Context, groupIDs []string) ([]*sdkws.GroupInfo, error) {
	findGroup := NewApiCaller[group.GetGroupsInfoReq, group.GetGroupsInfoResp]("/group/find_group")
	resp, err := findGroup.Call(ctx, &group.GetGroupsInfoReq{
		GroupIDs: groupIDs,
	})
	if err != nil {
		log.ZError(ctx, "FindGroup", err, "groupIDs", groupIDs)
		return nil, err
	}
	return resp.GroupInfos, nil
}

func (c *Caller) MapGroup(ctx context.Context, groupIDs []string) (map[string]*sdkws.GroupInfo, error) {
	mapGroup := NewApiCaller[group.GetGroupsInfoReq, group.GetGroupsInfoResp]("/group/map_group")
	resp, err := mapGroup.Call(ctx, &group.GetGroupsInfoReq{
		GroupIDs: groupIDs,
	})
	if err != nil {
		log.ZError(ctx, "MapGroup", err, "groupID", groupIDs)
		return nil, err
	}
	if len(resp.GroupInfos) == 0 {
		return nil, nil
	}
	return make(map[string]*sdkws.GroupInfo), nil
}

func (c *Caller) InviteToGroup(ctx context.Context, userID string, groupID string) error {
	inviteToGroup := NewApiCaller[group.InviteUserToGroupReq, group.InviteUserToGroupResp]("/group/invite_user_to_group")
	_, err := inviteToGroup.Call(ctx, &group.InviteUserToGroupReq{
		GroupID:        groupID,
		Reason:         "",
		InvitedUserIDs: []string{userID},
	})
	if err != nil {
		log.ZError(ctx, "inviteToGroup", err, "userID", userID, "groupID", groupID)
		return err
	}
	return nil
}

func (c *Caller) UpdateUserInfo(ctx context.Context, userID string, nickName string, faceURL string) error {
	updateUserInfo := NewApiCaller[user.UpdateUserInfoReq, user.UpdateUserInfoResp]("/user/update_user_info")
	_, err := updateUserInfo.Call(ctx, &user.UpdateUserInfoReq{UserInfo: &sdkws.UserInfo{
		UserID:   userID,
		Nickname: nickName,
		FaceURL:  faceURL,
	}})
	if err != nil {
		log.ZError(ctx, "updateUserInfo", err, "userID", userID)
		return err
	}
	return nil
}

func (c *Caller) RegisterUser(ctx context.Context, users []*sdkws.UserInfo) error {
	registerUser := NewApiCaller[user.UserRegisterReq, user.UserRegisterResp]("/user/user_register")
	_, err := registerUser.Call(ctx, &user.UserRegisterReq{
		Secret: *config.Config.Secret,
		Users:  users,
	})
	if err != nil {
		log.ZError(ctx, "RegisterUser", err)
		return err
	}
	return nil
}
