package api

import (
	"context"
	a2r "github.com/OpenIMSDK/Open-IM-Server/pkg/a2r"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/chat/pkg/common/config"
	"github.com/OpenIMSDK/chat/pkg/proto/admin"
	"github.com/OpenIMSDK/chat/pkg/proto/chat"
	"github.com/gin-gonic/gin"
)

func NewAdmin(zk discoveryregistry.SvcDiscoveryRegistry) *Admin {
	return &Admin{zk: zk}
}

type Admin struct {
	zk discoveryregistry.SvcDiscoveryRegistry
}

func (o *Admin) adminClient(ctx context.Context) (admin.AdminClient, error) {
	conn, err := o.zk.GetConn(ctx, config.Config.RpcRegisterName.OpenImAdminName)
	if err != nil {
		return nil, err
	}
	return admin.NewAdminClient(conn), nil
}

func (o *Admin) chatClient(ctx context.Context) (chat.ChatClient, error) {
	conn, err := o.zk.GetConn(ctx, config.Config.RpcRegisterName.OpenImChatName)
	if err != nil {
		return nil, err
	}
	return chat.NewChatClient(conn), nil
}

func (o *Admin) AdminLogin(c *gin.Context) {
	a2r.Call(admin.AdminClient.Login, o.adminClient, c)
}

func (o *Admin) ResetUserPassword(c *gin.Context) {
	a2r.Call(chat.ChatClient.ChangePassword, o.chatClient, c)
}

func (o *Admin) AdminUpdateInfo(c *gin.Context) {
	a2r.Call(admin.AdminClient.AdminUpdateInfo, o.adminClient, c)
}

func (o *Admin) AdminInfo(c *gin.Context) {
	a2r.Call(admin.AdminClient.GetAdminInfo, o.adminClient, c)
}

func (o *Admin) AddDefaultFriend(c *gin.Context) {
	a2r.Call(admin.AdminClient.AddDefaultFriend, o.adminClient, c)
}

func (o *Admin) DelDefaultFriend(c *gin.Context) {
	a2r.Call(admin.AdminClient.DelDefaultFriend, o.adminClient, c)
}

func (o *Admin) SearchDefaultFriend(c *gin.Context) {
	a2r.Call(admin.AdminClient.SearchDefaultFriend, o.adminClient, c)
}

func (o *Admin) FindDefaultFriend(c *gin.Context) {
	a2r.Call(admin.AdminClient.FindDefaultFriend, o.adminClient, c)
}

func (o *Admin) AddDefaultGroup(c *gin.Context) {
	a2r.Call(admin.AdminClient.AddDefaultGroup, o.adminClient, c)
}

func (o *Admin) DelDefaultGroup(c *gin.Context) {
	a2r.Call(admin.AdminClient.DelDefaultGroup, o.adminClient, c)
}

func (o *Admin) FindDefaultGroup(c *gin.Context) {
	a2r.Call(admin.AdminClient.FindDefaultGroup, o.adminClient, c)
}

func (o *Admin) SearchDefaultGroup(c *gin.Context) {
	a2r.Call(admin.AdminClient.SearchDefaultGroup, o.adminClient, c)
}

func (o *Admin) AddInvitationCode(c *gin.Context) {
	a2r.Call(admin.AdminClient.AddInvitationCode, o.adminClient, c)
}

func (o *Admin) GenInvitationCode(c *gin.Context) {
	a2r.Call(admin.AdminClient.GenInvitationCode, o.adminClient, c)
}

func (o *Admin) DelInvitationCode(c *gin.Context) {
	a2r.Call(admin.AdminClient.DelInvitationCode, o.adminClient, c)
}

func (o *Admin) SearchInvitationCode(c *gin.Context) {
	a2r.Call(admin.AdminClient.SearchInvitationCode, o.adminClient, c)
}

func (o *Admin) AddUserIPLimitLogin(c *gin.Context) {
	a2r.Call(admin.AdminClient.AddUserIPLimitLogin, o.adminClient, c)
}

func (o *Admin) SearchUserIPLimitLogin(c *gin.Context) {
	a2r.Call(admin.AdminClient.SearchUserIPLimitLogin, o.adminClient, c)
}

func (o *Admin) DelUserIPLimitLogin(c *gin.Context) {
	a2r.Call(admin.AdminClient.DelUserIPLimitLogin, o.adminClient, c)
}

func (o *Admin) SearchIPForbidden(c *gin.Context) {
	a2r.Call(admin.AdminClient.SearchIPForbidden, o.adminClient, c)
}

func (o *Admin) AddIPForbidden(c *gin.Context) {
	a2r.Call(admin.AdminClient.AddIPForbidden, o.adminClient, c)
}

func (o *Admin) DelIPForbidden(c *gin.Context) {
	a2r.Call(admin.AdminClient.DelIPForbidden, o.adminClient, c)
}

func (o *Admin) ParseToken(c *gin.Context) {
	a2r.Call(admin.AdminClient.ParseToken, o.adminClient, c)
}

func (o *Admin) BlockUser(c *gin.Context) {
	a2r.Call(admin.AdminClient.BlockUser, o.adminClient, c)
}

func (o *Admin) UnblockUser(c *gin.Context) {
	a2r.Call(admin.AdminClient.UnblockUser, o.adminClient, c)
}

func (o *Admin) SearchBlockUser(c *gin.Context) {
	a2r.Call(admin.AdminClient.SearchBlockUser, o.adminClient, c)
}

func (o *Admin) SetClientConfig(c *gin.Context) {
	a2r.Call(admin.AdminClient.SetClientConfig, o.adminClient, c)
}

func (o *Admin) GetClientConfig(c *gin.Context) {
	a2r.Call(admin.AdminClient.GetClientConfig, o.adminClient, c)
}

func (o *Admin) AddApplet(c *gin.Context) {
	a2r.Call(admin.AdminClient.AddApplet, o.adminClient, c)
}

func (o *Admin) DelApplet(c *gin.Context) {
	a2r.Call(admin.AdminClient.DelApplet, o.adminClient, c)
}

func (o *Admin) UpdateApplet(c *gin.Context) {
	a2r.Call(admin.AdminClient.UpdateApplet, o.adminClient, c)
}

func (o *Admin) SearchApplet(c *gin.Context) {
	a2r.Call(admin.AdminClient.SearchApplet, o.adminClient, c)
}
