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
	"github.com/OpenIMSDK/chat/pkg/common/apicall"
	"github.com/OpenIMSDK/chat/pkg/common/apistruct"
	"github.com/OpenIMSDK/chat/pkg/common/config"
	"github.com/OpenIMSDK/chat/pkg/common/mctx"
	"github.com/OpenIMSDK/chat/pkg/proto/admin"
	"github.com/OpenIMSDK/chat/pkg/proto/chat"
	"github.com/OpenIMSDK/protocol/constant"
	"github.com/OpenIMSDK/protocol/sdkws"
	"github.com/OpenIMSDK/protocol/user"
	"github.com/OpenIMSDK/tools/a2r"
	"github.com/OpenIMSDK/tools/apiresp"
	"github.com/OpenIMSDK/tools/checker"
	"github.com/OpenIMSDK/tools/errs"
	"github.com/OpenIMSDK/tools/log"
	"github.com/OpenIMSDK/tools/utils"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func NewAdmin(chatConn, adminConn grpc.ClientConnInterface) *AdminApi {
	return &AdminApi{chatClient: chat.NewChatClient(chatConn), adminClient: admin.NewAdminClient(adminConn), imApiCaller: apicall.NewCallerInterface()}
}

type AdminApi struct {
	chatClient  chat.ChatClient
	adminClient admin.AdminClient
	imApiCaller apicall.CallerInterface
}

func (o *AdminApi) AdminLogin(c *gin.Context) {
	var (
		req  admin.LoginReq
		resp apistruct.AdminLoginResp
	)
	if err := c.BindJSON(&req); err != nil {
		apiresp.GinError(c, err)
		return
	}
	log.ZInfo(c, "AdminLogin api", "req", &req)
	if err := checker.Validate(&req); err != nil {
		apiresp.GinError(c, err) // 参数校验失败
		return
	}
	loginResp, err := o.adminClient.Login(c, &req)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	imAdminUserID := config.GetIMAdmin(loginResp.AdminUserID)
	imToken, err := o.imApiCaller.UserToken(c, imAdminUserID, constant.AdminPlatformID)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	err = utils.CopyStructFields(&resp, loginResp)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	resp.ImToken = imToken
	resp.ImUserID = imAdminUserID
	log.ZInfo(c, "AdminLogin api", "resp", resp)
	apiresp.GinSuccess(c, resp)
}

func (o *AdminApi) ResetUserPassword(c *gin.Context) {
	a2r.Call(chat.ChatClient.ChangePassword, o.chatClient, c)
}

func (o *AdminApi) AdminUpdateInfo(c *gin.Context) {
	var req admin.AdminUpdateInfoReq
	if err := c.BindJSON(&req); err != nil {
		apiresp.GinError(c, err)
		return
	}
	if err := checker.Validate(&req); err != nil {
		apiresp.GinError(c, err) // 参数校验失败
		return
	}
	resp, err := o.adminClient.AdminUpdateInfo(c, &req)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiresp.GinSuccess(c, nil)
	imAdminUserID := config.GetIMAdmin(resp.UserID)
	imToken, err := o.imApiCaller.UserToken(c, imAdminUserID, constant.AdminPlatformID)
	if err != nil {
		log.ZError(c, "AdminUpdateInfo ImAdminTokenWithDefaultAdmin", err)
		return
	}
	if err := o.imApiCaller.UpdateUserInfo(mctx.WithApiToken(c, imToken), imAdminUserID, resp.Nickname, resp.FaceURL); err != nil {
		log.ZError(c, "AdminUpdateInfo UpdateUserInfo", err, "userID", resp.UserID, "nickName", resp.Nickname, "faceURL", resp.FaceURL)
	}
}

func (o *AdminApi) AdminInfo(c *gin.Context) {
	a2r.Call(admin.AdminClient.GetAdminInfo, o.adminClient, c)
}

func (o *AdminApi) AddDefaultFriend(c *gin.Context) {
	a2r.Call(admin.AdminClient.AddDefaultFriend, o.adminClient, c)
}

func (o *AdminApi) DelDefaultFriend(c *gin.Context) {
	a2r.Call(admin.AdminClient.DelDefaultFriend, o.adminClient, c)
}

func (o *AdminApi) SearchDefaultFriend(c *gin.Context) {
	a2r.Call(admin.AdminClient.SearchDefaultFriend, o.adminClient, c)
}

func (o *AdminApi) FindDefaultFriend(c *gin.Context) {
	a2r.Call(admin.AdminClient.FindDefaultFriend, o.adminClient, c)
}

func (o *AdminApi) AddDefaultGroup(c *gin.Context) {
	var req admin.AddDefaultGroupReq
	if err := c.BindJSON(&req); err != nil {
		apiresp.GinError(c, err)
		return
	}
	if err := checker.Validate(&req); err != nil {
		apiresp.GinError(c, err)
		return
	}
	imToken, err := o.imApiCaller.UserToken(c, config.GetIMAdmin(mctx.GetOpUserID(c)), constant.AdminPlatformID)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	groups, err := o.imApiCaller.FindGroupInfo(mctx.WithApiToken(c, imToken), req.GroupIDs)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	if len(req.GroupIDs) != len(groups) {
		apiresp.GinError(c, errs.ErrArgs.Wrap("group id not found"))
		return
	}
	resp, err := o.adminClient.AddDefaultGroup(c, &admin.AddDefaultGroupReq{
		GroupIDs: req.GroupIDs,
	})
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiresp.GinSuccess(c, resp)
}

func (o *AdminApi) DelDefaultGroup(c *gin.Context) {
	a2r.Call(admin.AdminClient.DelDefaultGroup, o.adminClient, c)
}

func (o *AdminApi) FindDefaultGroup(c *gin.Context) {
	a2r.Call(admin.AdminClient.FindDefaultGroup, o.adminClient, c)
}

func (o *AdminApi) SearchDefaultGroup(c *gin.Context) {
	var req admin.SearchDefaultGroupReq
	if err := c.BindJSON(&req); err != nil {
		apiresp.GinError(c, err)
		return
	}
	if err := checker.Validate(&req); err != nil {
		apiresp.GinError(c, err)
		return
	}
	searchResp, err := o.adminClient.SearchDefaultGroup(c, &req)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	resp := apistruct.SearchDefaultGroupResp{
		Total:  searchResp.Total,
		Groups: make([]*sdkws.GroupInfo, 0, len(searchResp.GroupIDs)),
	}
	if len(searchResp.GroupIDs) > 0 {
		imToken, err := o.imApiCaller.UserToken(c, config.GetIMAdmin(mctx.GetOpUserID(c)), constant.AdminPlatformID)
		if err != nil {
			apiresp.GinError(c, err)
			return
		}
		groups, err := o.imApiCaller.FindGroupInfo(mctx.WithApiToken(c, imToken), searchResp.GroupIDs)
		if err != nil {
			apiresp.GinError(c, err)
			return
		}
		groupMap := make(map[string]*sdkws.GroupInfo)
		for _, group := range groups {
			groupMap[group.GroupID] = group
		}
		for _, groupID := range searchResp.GroupIDs {
			if group, ok := groupMap[groupID]; ok {
				resp.Groups = append(resp.Groups, group)
			} else {
				resp.Groups = append(resp.Groups, &sdkws.GroupInfo{
					GroupID: groupID,
				})
			}
		}
	}
	apiresp.GinSuccess(c, resp)
}

func (o *AdminApi) AddInvitationCode(c *gin.Context) {
	a2r.Call(admin.AdminClient.AddInvitationCode, o.adminClient, c)
}

func (o *AdminApi) GenInvitationCode(c *gin.Context) {
	a2r.Call(admin.AdminClient.GenInvitationCode, o.adminClient, c)
}

func (o *AdminApi) DelInvitationCode(c *gin.Context) {
	a2r.Call(admin.AdminClient.DelInvitationCode, o.adminClient, c)
}

func (o *AdminApi) SearchInvitationCode(c *gin.Context) {
	a2r.Call(admin.AdminClient.SearchInvitationCode, o.adminClient, c)
}

func (o *AdminApi) AddUserIPLimitLogin(c *gin.Context) {
	a2r.Call(admin.AdminClient.AddUserIPLimitLogin, o.adminClient, c)
}

func (o *AdminApi) SearchUserIPLimitLogin(c *gin.Context) {
	a2r.Call(admin.AdminClient.SearchUserIPLimitLogin, o.adminClient, c)
}

func (o *AdminApi) DelUserIPLimitLogin(c *gin.Context) {
	a2r.Call(admin.AdminClient.DelUserIPLimitLogin, o.adminClient, c)
}

func (o *AdminApi) SearchIPForbidden(c *gin.Context) {
	a2r.Call(admin.AdminClient.SearchIPForbidden, o.adminClient, c)
}

func (o *AdminApi) AddIPForbidden(c *gin.Context) {
	a2r.Call(admin.AdminClient.AddIPForbidden, o.adminClient, c)
}

func (o *AdminApi) DelIPForbidden(c *gin.Context) {
	a2r.Call(admin.AdminClient.DelIPForbidden, o.adminClient, c)
}

func (o *AdminApi) ParseToken(c *gin.Context) {
	a2r.Call(admin.AdminClient.ParseToken, o.adminClient, c)
}

func (o *AdminApi) BlockUser(c *gin.Context) {
	var req admin.BlockUserReq
	if err := c.BindJSON(&req); err != nil {
		apiresp.GinError(c, err)
		return
	}
	log.ZInfo(c, "BlockUser api", "req", &req)
	if err := checker.Validate(&req); err != nil {
		apiresp.GinError(c, err) // 参数校验失败
		return
	}
	resp, err := o.adminClient.BlockUser(c, &req)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	imToken, err := o.imApiCaller.UserToken(c, config.GetIMAdmin(mctx.GetOpUserID(c)), constant.AdminPlatformID)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	err = o.imApiCaller.ForceOffLine(mctx.WithApiToken(c, imToken), req.UserID)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiresp.GinSuccess(c, resp)
}

func (o *AdminApi) UnblockUser(c *gin.Context) {
	a2r.Call(admin.AdminClient.UnblockUser, o.adminClient, c)
}

func (o *AdminApi) SearchBlockUser(c *gin.Context) {
	a2r.Call(admin.AdminClient.SearchBlockUser, o.adminClient, c)
}

func (o *AdminApi) SetClientConfig(c *gin.Context) {
	a2r.Call(admin.AdminClient.SetClientConfig, o.adminClient, c)
}

func (o *AdminApi) DelClientConfig(c *gin.Context) {
	a2r.Call(admin.AdminClient.DelClientConfig, o.adminClient, c)
}

func (o *AdminApi) GetClientConfig(c *gin.Context) {
	a2r.Call(admin.AdminClient.GetClientConfig, o.adminClient, c)
}

func (o *AdminApi) AddApplet(c *gin.Context) {
	a2r.Call(admin.AdminClient.AddApplet, o.adminClient, c)
}

func (o *AdminApi) DelApplet(c *gin.Context) {
	a2r.Call(admin.AdminClient.DelApplet, o.adminClient, c)
}

func (o *AdminApi) UpdateApplet(c *gin.Context) {
	a2r.Call(admin.AdminClient.UpdateApplet, o.adminClient, c)
}

func (o *AdminApi) SearchApplet(c *gin.Context) {
	a2r.Call(admin.AdminClient.SearchApplet, o.adminClient, c)
}

func (o *AdminApi) LoginUserCount(c *gin.Context) {
	a2r.Call(chat.ChatClient.UserLoginCount, o.chatClient, c)
}

func (o *AdminApi) NewUserCount(c *gin.Context) {
	var req user.UserRegisterCountReq
	var resp apistruct.NewUserCountResp
	if err := c.BindJSON(&req); err != nil {
		apiresp.GinError(c, err)
		return
	}
	log.ZInfo(c, "NewUserCount api", "req", &req)
	if err := checker.Validate(&req); err != nil {
		apiresp.GinError(c, err) // 参数校验失败
		return
	}
	imToken, err := o.imApiCaller.UserToken(c, config.GetIMAdmin(mctx.GetOpUserID(c)), constant.AdminPlatformID)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	dateCount, total, err := o.imApiCaller.UserRegisterCount(mctx.WithApiToken(c, imToken), req.Start, req.End)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	resp.DateCount = dateCount
	resp.Total = total
	apiresp.GinSuccess(c, resp)
}

func (o *AdminApi) SearchLogs(c *gin.Context) {
	a2r.Call(chat.ChatClient.SearchLogs, o.chatClient, c)
}

func (o *AdminApi) DeleteLogs(c *gin.Context) {
	a2r.Call(chat.ChatClient.DeleteLogs, o.chatClient, c)
}
