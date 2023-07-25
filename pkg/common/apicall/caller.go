package apicall

import (
	"context"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/auth"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/friend"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/group"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/user"
	"github.com/OpenIMSDK/chat/pkg/common/config"
)

type CallerInterface interface {
	AdminToken(ctx context.Context) (string, error)
	ImportFriend(ctx context.Context, ownerUserID string, friendUserID []string) error
	UserToken(ctx context.Context, userID string, platform int32) (string, error)
	InviteToGroup(ctx context.Context, userID string, groupIDs []string) error
	UpdateUserInfo(ctx context.Context, userID string, nickName string, faceURL string) error
	ForceOffLine(ctx context.Context, userID string) error
	RegisterUser(ctx context.Context, users []*sdkws.UserInfo) error
	FindGroupInfo(ctx context.Context, groupIDs []string) ([]*sdkws.GroupInfo, error)
}

type Caller struct{}

func NewCallerInterface() CallerInterface {
	return &Caller{}
}

// imporrt friend
func (c *Caller) ImportFriend(ctx context.Context, ownerUserID string, friendUserID []string, token string) error {
	importFriend := NewApiCaller[friend.ImportFriendReq, friend.ImportFriendResp]("/friend/import_friend")

	_, err := importFriend.Call(ctx, &friend.ImportFriendReq{
		OwnerUserID:   ownerUserID,
		FriendUserIDs: friendUserIDs,
	})
	return err
}

func (c *Caller) AdminToken(ctx context.Context) (string, error) {
	return c.UserToken(ctx, config.GetDefaultIMAdmin(), constant.AdminPlatformID)
}

// get user token
func (c *Caller) UserToken(ctx context.Context, userID string, platformID int32) (string, error) {
	resp, err := userToken.Call(ctx, &auth.UserTokenReq{
		Secret:     *config.Config.Secret,
		PlatformID: platformID,
		UserID:     userID,
	})
	if err != nil {
		return "", err
	}
	return resp.Token, nil
}

// invitate user to group
func (c *Caller) InviteToGroup(ctx context.Context, userID string, groupID string, token string) error {
	inviteToGroup := NewApiCaller[group.InviteUserToGroupReq, group.InviteUserToGroupResp]("/group/invite_user_to_group")
	_, err := inviteToGroup.Call(ctx, &group.InviteUserToGroupReq{
		GroupID:        groupID,
		Reason:         "",
		InvitedUserIDs: []string{userID},
	}, token)
	if err != nil {
		log.ZError(ctx, "inviteToGroup", err, "userID", userID, "groupID", groupID)
		return err
	}
	return nil
}

// update user info
func (c *Caller) UpdateUserInfo(ctx context.Context, userID string, nickName string, faceURL string, token string) error {
	updateUserInfo := NewApiCaller[user.UpdateUserInfoReq, user.UpdateUserInfoResp]("/user/update_user_info")

	_, err := updateUserInfo.Call(ctx, &user.UpdateUserInfoReq{UserInfo: &sdkws.UserInfo{
		UserID:   userID,
		Nickname: nickName,
		FaceURL:  faceURL,
	}})
	return err
}

// register user
func (c *Caller) RegisterUser(ctx context.Context, users []*sdkws.UserInfo) error {
	_, err := registerUser.Call(ctx, &user.UserRegisterReq{
		Secret: *config.Config.Secret,
		Users:  users,
	})
	return err
}

// force user offline
func (c *Caller) ForceOffLine(ctx context.Context, userID string, token string) error {

	for id := range constant.PlatformID2Name {
		_, _ = forceOffLine.Call(ctx, &auth.ForceLogoutReq{
			PlatformID: int32(id),
			UserID:     userID,
		})
	}
	return nil
}

func (c *Caller) FindGroupInfo(ctx context.Context, groupIDs []string) ([]*sdkws.GroupInfo, error) {
	resp, err := getGroupsInfo.Call(ctx, &group.GetGroupsInfoReq{
		GroupIDs: groupIDs,
	})
	if err != nil {
		return nil, err
	}
	return resp.GroupInfos, nil
}
