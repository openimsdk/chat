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
	"github.com/OpenIMSDK/Open-IM-Server/pkg/a2r"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/apiresp"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/checker"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
	"github.com/OpenIMSDK/chat/pkg/common/apicall"
	"github.com/OpenIMSDK/chat/pkg/common/apistruct"
	"github.com/OpenIMSDK/chat/pkg/common/mctx"
	"github.com/OpenIMSDK/chat/pkg/proto/admin"
	"github.com/OpenIMSDK/chat/pkg/proto/chat"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

// create a new admin channel
func NewAdmin(chatConn, adminConn grpc.ClientConnInterface) *AdminApi {
	return &AdminApi{chatClient: chat.NewChatClient(chatConn), adminClient: admin.NewAdminClient(adminConn), imApiCaller: apicall.NewCallerInterface()}
}

// define a struct named adminapi
type AdminApi struct {
	chatClient  chat.ChatClient
	adminClient admin.AdminClient
	imApiCaller apicall.CallerInterface
}

// admin login
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
	imToken, err := o.imApiCaller.UserToken(c, loginResp.AdminUserID, constant.AdminPlatformID)
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
	resp.ImUserID = loginResp.AdminUserID
	log.ZInfo(c, "AdminLogin api", "resp", resp)
	apiresp.GinSuccess(c, resp)
}

// reset user password
func (o *AdminApi) ResetUserPassword(c *gin.Context) {
	a2r.Call(chat.ChatClient.ChangePassword, o.chatClient, c)
}

// update admin-infomation
func (o *AdminApi) AdminUpdateInfo(c *gin.Context) {
	a2r.Call(admin.AdminClient.AdminUpdateInfo, o.adminClient, c)
}

// admin user info
func (o *AdminApi) AdminInfo(c *gin.Context) {
	a2r.Call(admin.AdminClient.GetAdminInfo, o.adminClient, c)
}

// add default friend
func (o *AdminApi) AddDefaultFriend(c *gin.Context) {
	a2r.Call(admin.AdminClient.AddDefaultFriend, o.adminClient, c)
}

// delete default friend
func (o *AdminApi) DelDefaultFriend(c *gin.Context) {
	a2r.Call(admin.AdminClient.DelDefaultFriend, o.adminClient, c)
}

// search default friend
func (o *AdminApi) SearchDefaultFriend(c *gin.Context) {
	a2r.Call(admin.AdminClient.SearchDefaultFriend, o.adminClient, c)
}

// search default friend
func (o *AdminApi) FindDefaultFriend(c *gin.Context) {
	a2r.Call(admin.AdminClient.FindDefaultFriend, o.adminClient, c)
}

// add fefault group
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
	token, err := o.imApiCaller.AdminToken(c)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	groups, err := o.imApiCaller.FindGroupInfo(mctx.WithApiToken(c, token), req.GroupIDs)
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

// delete default group
func (o *AdminApi) DelDefaultGroup(c *gin.Context) {
	a2r.Call(admin.AdminClient.DelDefaultGroup, o.adminClient, c)
}

// find fefault group
func (o *AdminApi) FindDefaultGroup(c *gin.Context) {
	a2r.Call(admin.AdminClient.FindDefaultGroup, o.adminClient, c)
}

// search default group
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
		token, err := o.imApiCaller.AdminToken(c)
		if err != nil {
			apiresp.GinError(c, err)
			return
		}
		groups, err := o.imApiCaller.FindGroupInfo(mctx.WithApiToken(c, token), searchResp.GroupIDs)
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

// add inviate code
func (o *AdminApi) AddInvitationCode(c *gin.Context) {
	a2r.Call(admin.AdminClient.AddInvitationCode, o.adminClient, c)
}

// generate invivate code
func (o *AdminApi) GenInvitationCode(c *gin.Context) {
	a2r.Call(admin.AdminClient.GenInvitationCode, o.adminClient, c)
}

// delete inviate code
func (o *AdminApi) DelInvitationCode(c *gin.Context) {
	a2r.Call(admin.AdminClient.DelInvitationCode, o.adminClient, c)
}

// search invitate code
func (o *AdminApi) SearchInvitationCode(c *gin.Context) {
	a2r.Call(admin.AdminClient.SearchInvitationCode, o.adminClient, c)
}

// add user login limit by ip
func (o *AdminApi) AddUserIPLimitLogin(c *gin.Context) {
	a2r.Call(admin.AdminClient.AddUserIPLimitLogin, o.adminClient, c)
}

// search userlogin by ip limit
func (o *AdminApi) SearchUserIPLimitLogin(c *gin.Context) {
	a2r.Call(admin.AdminClient.SearchUserIPLimitLogin, o.adminClient, c)
}

// delete user limit login by ip limit
func (o *AdminApi) DelUserIPLimitLogin(c *gin.Context) {
	a2r.Call(admin.AdminClient.DelUserIPLimitLogin, o.adminClient, c)
}

// search ip in 403
func (o *AdminApi) SearchIPForbidden(c *gin.Context) {
	a2r.Call(admin.AdminClient.SearchIPForbidden, o.adminClient, c)
}

// add ip into 403
func (o *AdminApi) AddIPForbidden(c *gin.Context) {
	a2r.Call(admin.AdminClient.AddIPForbidden, o.adminClient, c)
}

// delete ip from 403
func (o *AdminApi) DelIPForbidden(c *gin.Context) {
	a2r.Call(admin.AdminClient.DelIPForbidden, o.adminClient, c)
}

// parse token
func (o *AdminApi) ParseToken(c *gin.Context) {
	a2r.Call(admin.AdminClient.ParseToken, o.adminClient, c)
}

// block user
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
	token, err := o.imApiCaller.UserToken(c, mctx.GetOpUserID(c), constant.AdminPlatformID)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	err = o.imApiCaller.ForceOffLine(mctx.WithApiToken(c, token), req.UserID)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiresp.GinSuccess(c, resp)
}

// unblock user
func (o *AdminApi) UnblockUser(c *gin.Context) {
	a2r.Call(admin.AdminClient.UnblockUser, o.adminClient, c)
}

// search user blocked
func (o *AdminApi) SearchBlockUser(c *gin.Context) {
	a2r.Call(admin.AdminClient.SearchBlockUser, o.adminClient, c)
}

// set client config
func (o *AdminApi) SetClientConfig(c *gin.Context) {
	a2r.Call(admin.AdminClient.SetClientConfig, o.adminClient, c)
}

func (o *AdminApi) DelClientConfig(c *gin.Context) {
	a2r.Call(admin.AdminClient.DelClientConfig, o.adminClient, c)
}

func (o *AdminApi) GetClientConfig(c *gin.Context) {
	a2r.Call(admin.AdminClient.GetClientConfig, o.adminClient, c)
}

// add applet
func (o *AdminApi) AddApplet(c *gin.Context) {
	a2r.Call(admin.AdminClient.AddApplet, o.adminClient, c)
}

// delete applet
func (o *AdminApi) DelApplet(c *gin.Context) {
	a2r.Call(admin.AdminClient.DelApplet, o.adminClient, c)
}

// update applet
func (o *AdminApi) UpdateApplet(c *gin.Context) {
	a2r.Call(admin.AdminClient.UpdateApplet, o.adminClient, c)
}

// search applet
func (o *AdminApi) SearchApplet(c *gin.Context) {
	a2r.Call(admin.AdminClient.SearchApplet, o.adminClient, c)
}
