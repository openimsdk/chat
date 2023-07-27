package apicall

import (
	"context"
	"fmt"
	"github.com/OpenIMSDK/chat/pkg/common/config"
	"github.com/OpenIMSDK/protocol/auth"
	"github.com/OpenIMSDK/protocol/constant"
	"github.com/OpenIMSDK/protocol/friend"
	"github.com/OpenIMSDK/protocol/group"
	"github.com/OpenIMSDK/protocol/sdkws"
	"github.com/OpenIMSDK/protocol/user"
)

type CallerInterface interface {
	ImAdminTokenWithDefaultAdmin(ctx context.Context) (string, error)
	ImportFriend(ctx context.Context, ownerUserID string, friendUserID []string) error
	UserToken(ctx context.Context, userID string, platform int32) (string, error)
	InviteToGroup(ctx context.Context, userID string, groupIDs []string) error
	UpdateUserInfo(ctx context.Context, userID string, nickName string, faceURL string) error
	ForceOffLine(ctx context.Context, userID string) error
	RegisterUser(ctx context.Context, users []*sdkws.UserInfo) error
	FindGroupInfo(ctx context.Context, groupIDs []string) ([]*sdkws.GroupInfo, error)
}

type Caller struct {
}

func NewCallerInterface() CallerInterface {
	return &Caller{}
}

func (c *Caller) ImportFriend(ctx context.Context, ownerUserID string, friendUserIDs []string) error {
	if len(friendUserIDs) == 0 {
		return nil
	}
	_, err := importFriend.Call(ctx, &friend.ImportFriendReq{
		OwnerUserID:   ownerUserID,
		FriendUserIDs: friendUserIDs,
	})
	return err
}

func (c *Caller) ImAdminTokenWithDefaultAdmin(ctx context.Context) (string, error) {
	return c.UserToken(ctx, config.GetDefaultIMAdmin(), constant.AdminPlatformID)
}

func (c *Caller) UserToken(ctx context.Context, userID string, platformID int32) (string, error) {
	fmt.Println(*config.Config.Secret)
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

func (c *Caller) InviteToGroup(ctx context.Context, userID string, groupIDs []string) error {
	for _, groupID := range groupIDs {
		_, _ = inviteToGroup.Call(ctx, &group.InviteUserToGroupReq{
			GroupID:        groupID,
			Reason:         "",
			InvitedUserIDs: []string{userID},
		})
	}
	return nil
}

func (c *Caller) UpdateUserInfo(ctx context.Context, userID string, nickName string, faceURL string) error {
	_, err := updateUserInfo.Call(ctx, &user.UpdateUserInfoReq{UserInfo: &sdkws.UserInfo{
		UserID:   userID,
		Nickname: nickName,
		FaceURL:  faceURL,
	}})
	return err
}

func (c *Caller) RegisterUser(ctx context.Context, users []*sdkws.UserInfo) error {
	_, err := registerUser.Call(ctx, &user.UserRegisterReq{
		Secret: *config.Config.Secret,
		Users:  users,
	})
	return err
}

func (c *Caller) ForceOffLine(ctx context.Context, userID string) error {
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
