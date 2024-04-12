package util

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/openimsdk/chat/internal/api/mw"
	constant2 "github.com/openimsdk/chat/pkg/common/constant"
	"github.com/openimsdk/chat/pkg/common/mctx"
	"github.com/openimsdk/tools/errs"
	"net"
)

type Api struct {
	ImUserID        string
	ChatSecret      string
	ProxyHeader     string
	ChatAdminUserID string
}

func (o *Api) WithAdminUser(ctx context.Context) context.Context {
	return mctx.WithAdminUser(ctx, o.ChatAdminUserID)
}

func (o *Api) GetClientIP(c *gin.Context) (string, error) {
	if o.ProxyHeader == "" {
		ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)
		return ip, err
	}
	ip := c.Request.Header.Get(o.ProxyHeader)
	if ip == "" {
		return "", errs.ErrInternalServer.Wrap()
	}
	if ip := net.ParseIP(ip); ip == nil {
		return "", errs.ErrInternalServer.WrapMsg(fmt.Sprintf("parse proxy ip header %s failed", ip))
	}
	return ip, nil
}

func (o *Api) CheckSecretAdmin(c *gin.Context, secret string) error {
	if o.ChatSecret == "" {
		return errs.ErrNoPermission.WrapMsg("not config chat secret")
	}
	if _, ok := c.Get(constant2.RpcOpUserID); ok {
		return nil
	}
	if o.ChatSecret != secret {
		return errs.ErrNoPermission.WrapMsg("secret error")
	}
	mw.SetToken(c, o.GetDefaultIMAdminUserID(), constant2.AdminUser)
	return nil
}

func (o *Api) GetDefaultIMAdminUserID() string {
	return o.ImUserID
}
