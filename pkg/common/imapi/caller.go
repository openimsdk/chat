// Copyright © 2023 OpenIM open source community. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package imapi

import (
	"context"
	"sync"
	"time"

	"github.com/openimsdk/chat/pkg/eerrs"
	"github.com/openimsdk/tools/log"

	"github.com/openimsdk/protocol/auth"
	"github.com/openimsdk/protocol/constant"
	constantpb "github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/protocol/group"
	"github.com/openimsdk/protocol/relation"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/protocol/user"
)

type CallerInterface interface {
	ImAdminTokenWithDefaultAdmin(ctx context.Context) (string, error)
	ImportFriend(ctx context.Context, ownerUserID string, friendUserID []string) error
	GetUserToken(ctx context.Context, userID string, platform int32) (string, error)
	AdminToken(ctx context.Context, userID string) (string, error)
	InviteToGroup(ctx context.Context, userID string, groupIDs []string) error
	UpdateUserInfo(ctx context.Context, userID string, nickName string, faceURL string) error
	ForceOffLine(ctx context.Context, userID string) error
	RegisterUser(ctx context.Context, users []*sdkws.UserInfo) error
	FindGroupInfo(ctx context.Context, groupIDs []string) ([]*sdkws.GroupInfo, error)
	UserRegisterCount(ctx context.Context, start int64, end int64) (map[string]int64, int64, error)
	FriendUserIDs(ctx context.Context, userID string) ([]string, error)
	AccountCheckSingle(ctx context.Context, userID string) (bool, error)
}

type Caller struct {
	imApi           string
	imSecret        string
	defaultIMUserID string
	token           string
	timeout         time.Time
	lock            sync.Mutex
}

func New(imApi string, imSecret string, defaultIMUserID string) CallerInterface {
	return &Caller{
		imApi:           imApi,
		imSecret:        imSecret,
		defaultIMUserID: defaultIMUserID,
	}
}

func (c *Caller) ImportFriend(ctx context.Context, ownerUserID string, friendUserIDs []string) error {
	if len(friendUserIDs) == 0 {
		return nil
	}
	_, err := importFriend.Call(ctx, c.imApi, &relation.ImportFriendReq{
		OwnerUserID:   ownerUserID,
		FriendUserIDs: friendUserIDs,
	})
	return err
}

func (c *Caller) ImAdminTokenWithDefaultAdmin(ctx context.Context) (string, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.token == "" || c.timeout.Before(time.Now()) {
		userID := c.defaultIMUserID
		token, err := c.AdminToken(ctx, userID)
		if err != nil {
			log.ZError(ctx, "get im admin token", err, "userID", userID)
			return "", err
		}
		log.ZDebug(ctx, "get im admin token", "userID", userID)
		c.token = token
		c.timeout = time.Now().Add(time.Minute * 5)
	}
	return c.token, nil
}

func (c *Caller) AdminToken(ctx context.Context, userID string) (string, error) {
	resp, err := userToken.Call(ctx, c.imApi, &auth.UserTokenReq{
		Secret: c.imSecret,
		UserID: userID,
	})
	if err != nil {
		return "", err
	}
	return resp.Token, nil
}

func (c *Caller) GetUserToken(ctx context.Context, userID string, platformID int32) (string, error) {
	resp, err := getuserToken.Call(ctx, c.imApi, &auth.GetUserTokenReq{
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
		_, _ = inviteToGroup.Call(ctx, c.imApi, &group.InviteUserToGroupReq{
			GroupID:        groupID,
			Reason:         "",
			InvitedUserIDs: []string{userID},
		})
	}
	return nil
}

func (c *Caller) UpdateUserInfo(ctx context.Context, userID string, nickName string, faceURL string) error {
	_, err := updateUserInfo.Call(ctx, c.imApi, &user.UpdateUserInfoReq{UserInfo: &sdkws.UserInfo{
		UserID:   userID,
		Nickname: nickName,
		FaceURL:  faceURL,
	}})
	return err
}

func (c *Caller) RegisterUser(ctx context.Context, users []*sdkws.UserInfo) error {
	_, err := registerUser.Call(ctx, c.imApi, &user.UserRegisterReq{
		Users: users,
	})
	return err
}

func (c *Caller) ForceOffLine(ctx context.Context, userID string) error {
	for id := range constantpb.PlatformID2Name {
		_, _ = forceOffLine.Call(ctx, c.imApi, &auth.ForceLogoutReq{
			PlatformID: int32(id),
			UserID:     userID,
		})
	}
	return nil
}

func (c *Caller) FindGroupInfo(ctx context.Context, groupIDs []string) ([]*sdkws.GroupInfo, error) {
	resp, err := getGroupsInfo.Call(ctx, c.imApi, &group.GetGroupsInfoReq{
		GroupIDs: groupIDs,
	})
	if err != nil {
		return nil, err
	}
	return resp.GroupInfos, nil
}

func (c *Caller) UserRegisterCount(ctx context.Context, start int64, end int64) (map[string]int64, int64, error) {
	resp, err := registerUserCount.Call(ctx, c.imApi, &user.UserRegisterCountReq{
		Start: start,
		End:   end,
	})
	if err != nil {
		return nil, 0, err
	}
	return resp.Count, resp.Total, nil
}

func (c *Caller) FriendUserIDs(ctx context.Context, userID string) ([]string, error) {
	resp, err := friendUserIDs.Call(ctx, c.imApi, &relation.GetFriendIDsReq{UserID: userID})
	if err != nil {
		return nil, err
	}
	return resp.FriendIDs, nil
}

// return true when isUserNotExist.
func (c *Caller) AccountCheckSingle(ctx context.Context, userID string) (bool, error) {
	resp, err := accountCheck.Call(ctx, c.imApi, &user.AccountCheckReq{CheckUserIDs: []string{userID}})
	if err != nil {
		return false, err
	}
	if resp.Results[0].AccountStatus == constant.Registered {
		return false, eerrs.ErrAccountAlreadyRegister.Wrap()
	}
	return true, nil
}
