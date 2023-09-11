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

package apicall

import (
	"github.com/OpenIMSDK/chat/pkg/common/config"
	"github.com/OpenIMSDK/protocol/auth"
	"github.com/OpenIMSDK/protocol/friend"
	"github.com/OpenIMSDK/protocol/group"
	"github.com/OpenIMSDK/protocol/user"
)

func imApi() string {
	return config.Config.OpenIMUrl
}

// im caller.
var (
	importFriend      = NewApiCaller[friend.ImportFriendReq, friend.ImportFriendResp]("/friend/import_friend", imApi)
	userToken         = NewApiCaller[auth.UserTokenReq, auth.UserTokenResp]("/auth/user_token", imApi)
	inviteToGroup     = NewApiCaller[group.InviteUserToGroupReq, group.InviteUserToGroupResp]("/group/invite_user_to_group", imApi)
	updateUserInfo    = NewApiCaller[user.UpdateUserInfoReq, user.UpdateUserInfoResp]("/user/update_user_info", imApi)
	registerUser      = NewApiCaller[user.UserRegisterReq, user.UserRegisterResp]("/user/user_register", imApi)
	forceOffLine      = NewApiCaller[auth.ForceLogoutReq, auth.ForceLogoutResp]("/auth/force_logout", imApi)
	getGroupsInfo     = NewApiCaller[group.GetGroupsInfoReq, group.GetGroupsInfoResp]("/group/get_groups_info", imApi)
	registerUserCount = NewApiCaller[user.UserRegisterCountReq, user.UserRegisterCountResp]("/statistics/user/register", imApi)
	friendUserIDs     = NewApiCaller[friend.GetFriendIDsReq, friend.GetFriendIDsResp]("/friend/get_friend_id", imApi)
)
