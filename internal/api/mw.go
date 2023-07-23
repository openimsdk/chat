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
	"strconv"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/apiresp"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/chat/pkg/common/constant"
	"github.com/OpenIMSDK/chat/pkg/proto/admin"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func NewMW(adminConn grpc.ClientConnInterface) *MW {
	return &MW{client: admin.NewAdminClient(adminConn)}
}

// define a mv struct
type MW struct {
	client admin.AdminClient
}

// parse token
func (o *MW) parseToken(c *gin.Context) (string, int32, error) {
	token := c.GetHeader("token")
	if token == "" {
		return "", 0, errs.ErrArgs.Wrap("token is empty")
	}
	resp, err := o.client.ParseToken(c, &admin.ParseTokenReq{Token: token})
	if err != nil {
		return "", 0, err
	}
	return resp.UserID, resp.UserType, nil
}

// parse token by type
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

// set token
func (o *MW) setToken(c *gin.Context, userID string, userType int32) {
	c.Set(constant.RpcOpUserID, userID)
	c.Set(constant.RpcOpUserType, []string{strconv.Itoa(int(userType))})
	c.Set(constant.RpcCustomHeader, []string{constant.RpcOpUserType})
}

// check token
func (o *MW) CheckToken(c *gin.Context) {
	userID, userType, err := o.parseToken(c)
	if err != nil {
		c.Abort()
		apiresp.GinError(c, err)
		return
	}
	o.setToken(c, userID, userType)
}

// check admin info
func (o *MW) CheckAdmin(c *gin.Context) {
	userID, err := o.parseTokenType(c, constant.AdminUser)
	if err != nil {
		c.Abort()
		apiresp.GinError(c, err)
		return
	}
	o.setToken(c, userID, constant.AdminUser)
}

// checck user
func (o *MW) CheckUser(c *gin.Context) {
	userID, err := o.parseTokenType(c, constant.NormalUser)
	if err != nil {
		c.Abort()
		apiresp.GinError(c, err)
		return
	}
	o.setToken(c, userID, constant.NormalUser)
}
