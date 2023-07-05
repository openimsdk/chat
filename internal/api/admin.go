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
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"

	"github.com/OpenIMSDK/chat/pkg/proto/admin"
	"github.com/OpenIMSDK/chat/pkg/proto/chat"
)

func NewAdmin(chatConn, adminConn grpc.ClientConnInterface) *AdminApi {
	return &AdminApi{chatClient: chat.NewChatClient(chatConn), adminClient: admin.NewAdminClient(adminConn)}
}

type AdminApi struct {
	chatClient  chat.ChatClient
	adminClient admin.AdminClient
}

func (o *AdminApi) AdminLogin(c *gin.Context) {
	a2r.Call(admin.AdminClient.Login, o.adminClient, c)
}

func (o *AdminApi) ResetUserPassword(c *gin.Context) {
	a2r.Call(chat.ChatClient.ChangePassword, o.chatClient, c)
}

func (o *AdminApi) AdminUpdateInfo(c *gin.Context) {
	a2r.Call(admin.AdminClient.AdminUpdateInfo, o.adminClient, c)
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
	a2r.Call(admin.AdminClient.AddDefaultGroup, o.adminClient, c)
}

func (o *AdminApi) DelDefaultGroup(c *gin.Context) {
	a2r.Call(admin.AdminClient.DelDefaultGroup, o.adminClient, c)
}

func (o *AdminApi) FindDefaultGroup(c *gin.Context) {
	a2r.Call(admin.AdminClient.FindDefaultGroup, o.adminClient, c)
}

func (o *AdminApi) SearchDefaultGroup(c *gin.Context) {
	a2r.Call(admin.AdminClient.SearchDefaultGroup, o.adminClient, c)
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
	a2r.Call(admin.AdminClient.BlockUser, o.adminClient, c)
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
