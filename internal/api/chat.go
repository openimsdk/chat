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

package api

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/a2r"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/apiresp"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/checker"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"github.com/OpenIMSDK/chat/pkg/common/apicall"
	"github.com/OpenIMSDK/chat/pkg/common/apistruct"
	"github.com/OpenIMSDK/chat/pkg/common/config"
	constant2 "github.com/OpenIMSDK/chat/pkg/common/constant"
	"github.com/OpenIMSDK/chat/pkg/common/mctx"
	"github.com/OpenIMSDK/chat/pkg/proto/admin"
	"github.com/OpenIMSDK/chat/pkg/proto/chat"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

// create a new chat
func NewChat(chatConn, adminConn grpc.ClientConnInterface) *ChatApi {
	return &ChatApi{chatClient: chat.NewChatClient(chatConn), adminClient: admin.NewAdminClient(adminConn), imApiCaller: apicall.NewCallerInterface()}
}

// define a struct named chatapi
type ChatApi struct {
	chatClient  chat.ChatClient
	adminClient admin.AdminClient
	imApiCaller apicall.CallerInterface
}

// ################## ACCOUNT ##################

func (o *ChatApi) SendVerifyCode(c *gin.Context) {
	var req chat.SendVerifyCodeReq
	if err := c.BindJSON(&req); err != nil {
		apiresp.GinError(c, err)
		return
	}
	ip, err := o.getClientIP(c)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	req.Ip = ip
	resp, err := o.chatClient.SendVerifyCode(c, &req)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiresp.GinSuccess(c, resp)
}

// vertify code
func (o *ChatApi) VerifyCode(c *gin.Context) {
	a2r.Call(chat.ChatClient.VerifyCode, o.chatClient, c)
}

// user registe
func (o *ChatApi) RegisterUser(c *gin.Context) {
	var (
		req  chat.RegisterUserReq
		resp apistruct.UserRegisterResp
	)
	if err := c.BindJSON(&req); err != nil {
		apiresp.GinError(c, err)
		return
	}
	log.ZInfo(c, "registerUser", "req", &req)
	if err := checker.Validate(&req); err != nil {
		apiresp.GinError(c, errs.ErrArgs.Wrap(err.Error())) // 参数校验失败
		return
	}
	ip, err := o.getClientIP(c)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	req.Ip = ip
	resp1, err := o.chatClient.RegisterUser(c, &req)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	userInfo := &sdkws.UserInfo{
		UserID:     resp1.UserID,
		Nickname:   req.User.Nickname,
		FaceURL:    req.User.FaceURL,
		CreateTime: time.Now().UnixMilli(),
	}
	err = o.imApiCaller.RegisterUser(c, []*sdkws.UserInfo{userInfo})
	if err != nil {
		apiresp.GinError(c, err)
		return
	}

	imAdminID := config.GetDefaultIMAdmin()
	token, err := o.imApiCaller.UserToken(c, imAdminID, constant.AdminPlatformID)
	if err != nil {
		log.ZError(c, "GetIMAdminUserToken Failed", err, "userID", imAdminID)
		apiresp.GinError(c, err)
		return
	}
	// c.Set(constant.Token, token)

	t := mctx.WithOpUserID(c, config.Config.AdminList[0].AdminID, constant2.AdminUser)
	resp2, err := o.adminClient.FindDefaultFriend(t, &admin.FindDefaultFriendReq{})
	if err != nil {
		log.ZError(t, "FindDefaultFriend Failed", err, "userID", req.User.UserID)
		apiresp.GinError(c, err)
		return
	} else if len(resp2.UserIDs) > 0 {
		if err := o.imApiCaller.ImportFriend(c, resp1.UserID, resp2.UserIDs, token); err != nil {
			apiresp.GinError(c, err)
			return
		}
	}

	resp3, err := o.adminClient.FindDefaultGroup(t, &admin.FindDefaultGroupReq{})
	if err != nil {
		log.ZError(t, "FindDefaultGroupID Failed", err, "userID", req.User.UserID)
		apiresp.GinError(c, err)
		return
	} else if len(resp3.GroupIDs) > 0 {
		for _, groupID := range resp3.GroupIDs {
			if err := o.imApiCaller.InviteToGroup(c, resp1.UserID, groupID, token); err != nil {
				log.ZError(c, "inviteUserToGroup Failed", err, "userID", req.User.UserID, "groupID", groupID)
				apiresp.GinError(c, err)
				return
			}
		}
	}
	if req.AutoLogin {
		token, err := o.imApiCaller.UserToken(c, resp1.UserID, req.Platform)
		if err != nil {
			log.ZError(c, "GetIMAdminUserToken Failed", err, "userID", req.User.UserID)
			apiresp.GinError(c, err)
			return
		}
		resp.ImToken = token
	}
	resp.ChatToken = resp1.ChatToken
	resp.UserID = resp1.UserID
	log.ZInfo(c, "registerUser api", "resp", &resp)
	apiresp.GinSuccess(c, &resp)
	// a2r.Call(chat.ChatClient.RegisterUser, o.chatClient, c)
}

// user login
func (o *ChatApi) Login(c *gin.Context) {
	var (
		req  chat.LoginReq
		resp apistruct.LoginResp
	)
	if err := c.BindJSON(&req); err != nil {
		apiresp.GinError(c, err)
		return
	}
	if err := checker.Validate(&req); err != nil {
		apiresp.GinError(c, errs.ErrArgs.Wrap(err.Error())) // 参数校验失败
		return
	}
	log.ZInfo(c, "Login", "req", &req)
	ip, err := o.getClientIP(c)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	req.Ip = ip
	resp1, err := o.chatClient.Login(c, &req)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	imToken, err := o.imApiCaller.UserToken(c, resp1.UserID, constant2.NormalUser)
	if err != nil {
		return
	}
	resp.ImToken = imToken
	resp.UserID = resp1.UserID
	resp.ChatToken = resp1.ChatToken
	log.ZInfo(c, "Login api", "resp", &resp)
	apiresp.GinSuccess(c, resp)
}

// reset user password
func (o *ChatApi) ResetPassword(c *gin.Context) {
	a2r.Call(chat.ChatClient.ResetPassword, o.chatClient, c)
}

// change user password
func (o *ChatApi) ChangePassword(c *gin.Context) {
	a2r.Call(chat.ChatClient.ChangePassword, o.chatClient, c)
}

// ################## USER ##################
// update user profile
func (o *ChatApi) UpdateUserInfo(c *gin.Context) {
	var (
		req        chat.UpdateUserInfoReq
		resp       apistruct.UpdateUserInfoResp
		imUserID   string
		platformID int32
		nickName   string
		faceURL    string
	)
	if err := c.BindJSON(&req); err != nil {
		apiresp.GinError(c, err)
		return
	}
	log.ZInfo(c, "updateUserInfo", "req", &req)
	if err := checker.Validate(&req); err != nil {
		apiresp.GinError(c, errs.ErrArgs.Wrap(err.Error())) // 参数校验失败
		return
	}
	resp1, err := o.chatClient.UpdateUserInfo(c, &req)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	opUserID := mctx.GetOpUserID(c)
	opUserType := mctx.GetUserType(c)
	if opUserType == constant2.AdminUser {
		platformID = constant.AdminPlatformID
		imUserID = config.GetIMAdmin(opUserID)
		if imUserID == "" {
			apiresp.GinError(c, errs.ErrUserIDNotFound.Wrap("chatAdminID to imAdminID error"))
			return
		}
	} else {
		platformID = constant2.DefaultPlatform
		imUserID = req.UserID
	}
	token, err := o.imApiCaller.UserToken(c, imUserID, platformID)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	// c.Set(constant.Token, token)

	if req.Nickname != nil {
		nickName = req.Nickname.Value
	} else {
		nickName = resp1.NickName
	}
	if req.FaceURL != nil {
		faceURL = req.FaceURL.Value
	} else {
		faceURL = resp1.FaceUrl
	}
	err = o.imApiCaller.UpdateUserInfo(c, req.UserID, nickName, faceURL, token)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	log.ZInfo(c, "updateUserInfo", "resp", &resp)
	apiresp.GinSuccess(c, resp)
}

// find user public user info
func (o *ChatApi) FindUserPublicInfo(c *gin.Context) {
	a2r.Call(chat.ChatClient.FindUserPublicInfo, o.chatClient, c)
}

// find user full information
func (o *ChatApi) FindUserFullInfo(c *gin.Context) {
	a2r.Call(chat.ChatClient.FindUserFullInfo, o.chatClient, c)
}

//func (o *ChatApi) GetUsersFullInfo(c *gin.Context) {
//	a2r.Call(chat.ChatClient.GetUsersFullInfo, o.chatClient, c)
//}

func (o *ChatApi) SearchUserFullInfo(c *gin.Context) {
	a2r.Call(chat.ChatClient.SearchUserFullInfo, o.chatClient, c)
}

// search user public information
func (o *ChatApi) SearchUserPublicInfo(c *gin.Context) {
	a2r.Call(chat.ChatClient.SearchUserPublicInfo, o.chatClient, c)
}

// ################## APPLET ##################
// find applet
func (o *ChatApi) FindApplet(c *gin.Context) {
	a2r.Call(admin.AdminClient.FindApplet, o.adminClient, c)
}

// ################## CONFIG ##################

func (o *ChatApi) GetClientConfig(c *gin.Context) {
	a2r.Call(admin.AdminClient.GetClientConfig, o.adminClient, c)
}

// ################## CALLBACK ##################
// openim cb
func (o *ChatApi) OpenIMCallback(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	req := &chat.OpenIMCallbackReq{
		Command: c.Query(constant.CallbackCommand),
		Body:    string(body),
	}
	if _, err := o.chatClient.OpenIMCallback(c, req); err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiresp.GinSuccess(c, nil)
}

// get user client ip
func (o *ChatApi) getClientIP(c *gin.Context) (string, error) {
	if config.Config.ProxyHeader == "" {
		ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)
		return ip, err
	}
	ip := c.Request.Header.Get(config.Config.ProxyHeader)
	if ip == "" {
		return "", errs.ErrInternalServer.Wrap()
	}
	if ip := net.ParseIP(ip); ip == nil {
		return "", errs.ErrInternalServer.Wrap(fmt.Sprintf("parse proxy ip header %s failed", ip))
	}
	return ip, nil
}
