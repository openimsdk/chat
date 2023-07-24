package apicall

import (
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/auth"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/friend"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/group"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/user"
	"github.com/OpenIMSDK/chat/pkg/common/config"
)

// im api
var (
	importFriend   = NewApiCaller[friend.ImportFriendReq, friend.ImportFriendResp](config.Config.OpenIMUrl + "/friend/import_friend")
	userToken      = NewApiCaller[auth.UserTokenReq, auth.UserTokenResp](config.Config.OpenIMUrl + "/auth/user_token")
	inviteToGroup  = NewApiCaller[group.InviteUserToGroupReq, group.InviteUserToGroupResp](config.Config.OpenIMUrl + "/group/invite_user_to_group")
	updateUserInfo = NewApiCaller[user.UpdateUserInfoReq, user.UpdateUserInfoResp](config.Config.OpenIMUrl + "/user/update_user_info")
	registerUser   = NewApiCaller[user.UserRegisterReq, user.UserRegisterResp](config.Config.OpenIMUrl + "/user/user_register")
	forceOffLine   = NewApiCaller[auth.ForceLogoutReq, auth.ForceLogoutResp](config.Config.OpenIMUrl + "/auth/force_logout")
)
