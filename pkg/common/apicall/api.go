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

// im caller
var (
	importFriend   = NewApiCaller[friend.ImportFriendReq, friend.ImportFriendResp]("/friend/import_friend", imApi)
	userToken      = NewApiCaller[auth.UserTokenReq, auth.UserTokenResp]("/auth/user_token", imApi)
	inviteToGroup  = NewApiCaller[group.InviteUserToGroupReq, group.InviteUserToGroupResp]("/group/invite_user_to_group", imApi)
	updateUserInfo = NewApiCaller[user.UpdateUserInfoReq, user.UpdateUserInfoResp]("/user/update_user_info", imApi)
	registerUser   = NewApiCaller[user.UserRegisterReq, user.UserRegisterResp]("/user/user_register", imApi)
	forceOffLine   = NewApiCaller[auth.ForceLogoutReq, auth.ForceLogoutResp]("/auth/force_logout", imApi)
	getGroupsInfo  = NewApiCaller[group.GetGroupsInfoReq, group.GetGroupsInfoResp]("/group/get_groups_info", imApi)
)
