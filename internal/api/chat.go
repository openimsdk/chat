package api

import (
	"context"
	"fmt"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/a2r"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/apiresp"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/chat/pkg/common/config"
	"github.com/OpenIMSDK/chat/pkg/proto/admin"
	"github.com/OpenIMSDK/chat/pkg/proto/chat"
	"github.com/gin-gonic/gin"
	"io"
	"net"
)

func NewChat(zk discoveryregistry.SvcDiscoveryRegistry) *Chat {
	return &Chat{zk: zk}
}

type Chat struct {
	zk discoveryregistry.SvcDiscoveryRegistry
}

func (o *Chat) chatClient(ctx context.Context) (chat.ChatClient, error) {
	conn, err := o.zk.GetConn(ctx, config.Config.RpcRegisterName.OpenImChatName)
	if err != nil {
		return nil, err
	}
	return chat.NewChatClient(conn), nil
}

func (o *Chat) adminClient(ctx context.Context) (admin.AdminClient, error) {
	conn, err := o.zk.GetConn(ctx, config.Config.RpcRegisterName.OpenImAdminName)
	if err != nil {
		return nil, err
	}
	return admin.NewAdminClient(conn), nil
}

// ################## ACCOUNT ##################

func (o *Chat) SendVerifyCode(c *gin.Context) {
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
	client, err := o.chatClient(c)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	resp, err := client.SendVerifyCode(c, &req)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiresp.GinSuccess(c, resp)
	//a2r.Call(chat.ChatClient.SendVerifyCode, o.chatClient, c)
}

func (o *Chat) VerifyCode(c *gin.Context) {
	a2r.Call(chat.ChatClient.VerifyCode, o.chatClient, c)
}

func (o *Chat) RegisterUser(c *gin.Context) {
	var req chat.RegisterUserReq
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
	client, err := o.chatClient(c)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	resp, err := client.RegisterUser(c, &req)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiresp.GinSuccess(c, resp)
	//a2r.Call(chat.ChatClient.RegisterUser, o.chatClient, c)
}

func (o *Chat) Login(c *gin.Context) {
	var req chat.LoginReq
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
	client, err := o.chatClient(c)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	resp, err := client.Login(c, &req)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiresp.GinSuccess(c, resp)
	//a2r.Call(chat.ChatClient.Login, o.chatClient, c)
}

func (o *Chat) ResetPassword(c *gin.Context) {
	a2r.Call(chat.ChatClient.ResetPassword, o.chatClient, c)
}

func (o *Chat) ChangePassword(c *gin.Context) {
	a2r.Call(chat.ChatClient.ResetPassword, o.chatClient, c)
}

// ################## USER ##################

func (o *Chat) UpdateUserInfo(c *gin.Context) {
	a2r.Call(chat.ChatClient.UpdateUserInfo, o.chatClient, c)
}

func (o *Chat) FindUserPublicInfo(c *gin.Context) {
	a2r.Call(chat.ChatClient.FindUserPublicInfo, o.chatClient, c)
}

func (o *Chat) FindUserFullInfo(c *gin.Context) {
	a2r.Call(chat.ChatClient.FindUserFullInfo, o.chatClient, c)
}

//func (o *Chat) GetUsersFullInfo(c *gin.Context) {
//	a2r.Call(chat.ChatClient.GetUsersFullInfo, o.chatClient, c)
//}

func (o *Chat) SearchUserFullInfo(c *gin.Context) {
	a2r.Call(chat.ChatClient.SearchUserFullInfo, o.chatClient, c)
}

func (o *Chat) SearchUserPublicInfo(c *gin.Context) {
	a2r.Call(chat.ChatClient.SearchUserPublicInfo, o.chatClient, c)
}

// ################## APPLET ##################

func (o *Chat) FindApplet(c *gin.Context) {
	a2r.Call(admin.AdminClient.FindApplet, o.adminClient, c)
}

// ################## CONFIG ##################

func (o *Chat) GetClientConfig(c *gin.Context) {
	a2r.Call(admin.AdminClient.GetClientConfig, o.adminClient, c)
}

// ################## CALLBACK ##################

func (o *Chat) OpenIMCallback(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	cli, err := o.chatClient(c)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	req := &chat.OpenIMCallbackReq{
		Command: c.Query(constant.CallbackCommand),
		Body:    string(body),
	}
	if _, err := cli.OpenIMCallback(c, req); err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiresp.GinSuccess(c, nil)
}

func (o *Chat) getClientIP(c *gin.Context) (string, error) {
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
