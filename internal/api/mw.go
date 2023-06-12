package api

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/apiresp"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/chat/pkg/common/config"
	"github.com/OpenIMSDK/chat/pkg/common/constant"
	"github.com/OpenIMSDK/chat/pkg/proto/admin"
	"github.com/gin-gonic/gin"
	"strconv"
)

func NewMW(zk discoveryregistry.SvcDiscoveryRegistry) *MW {
	return &MW{zk: zk}
}

type MW struct {
	zk discoveryregistry.SvcDiscoveryRegistry
}

func (o *MW) adminClient(ctx context.Context) (admin.AdminClient, error) {
	conn, err := o.zk.GetConn(ctx, config.Config.RpcRegisterName.OpenImAdminName)
	if err != nil {
		return nil, err
	}
	return admin.NewAdminClient(conn), nil
}

func (o *MW) parseToken(c *gin.Context) (string, int32, error) {
	token := c.GetHeader("token")
	if token == "" {
		return "", 0, errs.ErrArgs.Wrap("token is empty")
	}
	client, err := o.adminClient(c)
	if err != nil {
		return "", 0, err
	}
	resp, err := client.ParseToken(c, &admin.ParseTokenReq{Token: token})
	if err != nil {
		return "", 0, err
	}
	return resp.UserID, resp.UserType, nil
}

func (o *MW) parseTokenType(c *gin.Context, userType int32) (string, error) {
	userID, t, err := o.parseToken(c)
	if err != nil {
		return "", err
	}
	if t != userType {
		return "", errs.ErrArgs.Wrap("token type error")
	}
	return userID, nil
}

func (o *MW) setToken(c *gin.Context, userID string, userType int32) {
	c.Set(constant.RpcOpUserID, userID)
	c.Set(constant.RpcOpUserType, []string{strconv.Itoa(int(userType))})
	c.Set(constant.RpcCustomHeader, []string{constant.RpcOpUserType})
}

func (o *MW) CheckToken(c *gin.Context) {
	userID, userType, err := o.parseToken(c)
	if err != nil {
		c.Abort()
		apiresp.GinError(c, err)
		return
	}
	o.setToken(c, userID, userType)
}

func (o *MW) CheckAdmin(c *gin.Context) {
	userID, err := o.parseTokenType(c, constant.AdminUser)
	if err != nil {
		c.Abort()
		apiresp.GinError(c, err)
		return
	}
	o.setToken(c, userID, constant.AdminUser)
}

func (o *MW) CheckUser(c *gin.Context) {
	userID, err := o.parseTokenType(c, constant.NormalUser)
	if err != nil {
		c.Abort()
		apiresp.GinError(c, err)
		return
	}
	o.setToken(c, userID, constant.NormalUser)
}

func (o *MW) CheckAdminOrNil(c *gin.Context) {

}
