package apicall

import (
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/auth"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/friend"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/group"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/user"
	"github.com/OpenIMSDK/chat/pkg/common/config"
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
)
