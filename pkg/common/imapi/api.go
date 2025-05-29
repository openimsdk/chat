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
	"github.com/openimsdk/protocol/auth"
	"github.com/openimsdk/protocol/group"
	"github.com/openimsdk/protocol/relation"
	"github.com/openimsdk/protocol/user"
)

// im caller.
var (
	getAdminToken = NewApiCaller[auth.GetAdminTokenReq, auth.GetAdminTokenResp]("/auth/get_admin_token")
	getuserToken  = NewApiCaller[auth.GetUserTokenReq, auth.GetUserTokenResp]("/auth/get_user_token")
	forceOffLine  = NewApiCaller[auth.ForceLogoutReq, auth.ForceLogoutResp]("/auth/force_logout")

	updateUserInfo            = NewApiCaller[user.UpdateUserInfoReq, user.UpdateUserInfoResp]("/user/update_user_info")
	registerUser              = NewApiCaller[user.UserRegisterReq, user.UserRegisterResp]("/user/user_register")
	getUserInfo               = NewApiCaller[user.GetDesignateUsersReq, user.GetDesignateUsersResp]("/user/get_users_info")
	accountCheck              = NewApiCaller[user.AccountCheckReq, user.AccountCheckResp]("/user/account_check")
	addNotificationAccount    = NewApiCaller[user.AddNotificationAccountReq, user.AddNotificationAccountResp]("/user/add_notification_account")
	updateNotificationAccount = NewApiCaller[user.UpdateNotificationAccountInfoReq, user.UpdateNotificationAccountInfoResp]("/user/update_notification_account")

	getGroupsInfo = NewApiCaller[group.GetGroupsInfoReq, group.GetGroupsInfoResp]("/group/get_groups_info")
	inviteToGroup = NewApiCaller[group.InviteUserToGroupReq, group.InviteUserToGroupResp]("/group/invite_user_to_group")

	registerUserCount = NewApiCaller[user.UserRegisterCountReq, user.UserRegisterCountResp]("/statistics/user/register")

	friendUserIDs = NewApiCaller[relation.GetFriendIDsReq, relation.GetFriendIDsResp]("/friend/get_friend_id")
	importFriend  = NewApiCaller[relation.ImportFriendReq, relation.ImportFriendResp]("/friend/import_friend")

	sendSimpleMsg = NewApiCaller[SendSingleMsgReq, SendSingleMsgResp]("/msg/send_simple_msg")
)
