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
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/openimsdk/chat/pkg/common/apicall"
	"github.com/openimsdk/chat/pkg/common/apistruct"
	"github.com/openimsdk/chat/pkg/common/config"
	constant2 "github.com/openimsdk/chat/pkg/common/constant"
	"github.com/openimsdk/chat/pkg/common/mctx"
	"github.com/openimsdk/chat/pkg/common/xlsx"
	"github.com/openimsdk/chat/pkg/common/xlsx/model"
	"github.com/openimsdk/chat/pkg/proto/admin"
	"github.com/openimsdk/chat/pkg/proto/chat"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/protocol/user"
	"github.com/openimsdk/tools/a2r"
	"github.com/openimsdk/tools/apiresp"
	"github.com/openimsdk/tools/checker"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/utils"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
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

	imAdminUserID := config.GetIMAdmin(resp.UserID)
	imToken, err := o.imApiCaller.UserToken(c, imAdminUserID, constant.AdminPlatformID)
	if err != nil {
		log.ZError(c, "AdminUpdateInfo ImAdminTokenWithDefaultAdmin", err, "imAdminUserID", imAdminUserID)
		return
	}
	if err := o.imApiCaller.UpdateUserInfo(mctx.WithApiToken(c, imToken), imAdminUserID, resp.Nickname, resp.FaceURL); err != nil {
		log.ZError(c, "AdminUpdateInfo UpdateUserInfo", err, "userID", resp.UserID, "nickName", resp.Nickname, "faceURL", resp.FaceURL)
	}
	apiresp.GinSuccess(c, nil)
}

func (o *AdminApi) AdminInfo(c *gin.Context) {
	a2r.Call(admin.AdminClient.GetAdminInfo, o.adminClient, c)
}

func (o *AdminApi) ChangeAdminPassword(c *gin.Context) {
	a2r.Call(admin.AdminClient.ChangeAdminPassword, o.adminClient, c)
}

func (o *AdminApi) AddAdminAccount(c *gin.Context) {
	a2r.Call(admin.AdminClient.AddAdminAccount, o.adminClient, c)
}

func (o *AdminApi) AddUserAccount(c *gin.Context) {
	var req chat.AddUserAccountReq
	if err := c.BindJSON(&req); err != nil {
		apiresp.GinError(c, err)
		return
	}
	if err := checker.Validate(&req); err != nil {
		apiresp.GinError(c, err) // 参数校验失败
		return
	}

	_, err := o.chatClient.AddUserAccount(c, &req)

	userInfo := &sdkws.UserInfo{
		UserID:     req.User.UserID,
		Nickname:   req.User.Nickname,
		FaceURL:    req.User.FaceURL,
		CreateTime: time.Now().UnixMilli(),
	}
	err = o.imApiCaller.RegisterUser(c, []*sdkws.UserInfo{userInfo})
	if err != nil {
		apiresp.GinError(c, err)
		return
	}

	apiresp.GinSuccess(c, nil)

}

func (o *AdminApi) DelAdminAccount(c *gin.Context) {
	a2r.Call(admin.AdminClient.DelAdminAccount, o.adminClient, c)
}

func (o *AdminApi) SearchAdminAccount(c *gin.Context) {
	a2r.Call(admin.AdminClient.SearchAdminAccount, o.adminClient, c)
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
		apiresp.GinError(c, errs.ErrArgs.WrapMsg("group id not found"))
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

func (o *AdminApi) getClientIP(c *gin.Context) (string, error) {
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

func (o *AdminApi) checkSecretAdmin(c *gin.Context, secret string) error {
	if _, ok := c.Get(constant2.RpcOpUserID); ok {
		return nil
	}
	if config.Config.ChatSecret == "" {
		return errs.ErrNoPermission.WrapMsg("not config chat secret")
	}
	if config.Config.ChatSecret != secret {
		return errs.ErrNoPermission.WrapMsg("secret error")
	}
	SetToken(c, config.GetDefaultIMAdmin(), constant2.AdminUser)
	return nil
}

func (o *AdminApi) ImportUserByXlsx(c *gin.Context) {
	formFile, err := c.FormFile("data")
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	ip, err := o.getClientIP(c)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	secret := c.PostForm("secret")
	if err := o.checkSecretAdmin(c, secret); err != nil {
		apiresp.GinError(c, err)
		return
	}
	file, err := formFile.Open()
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	defer file.Close()
	var users []model.User
	if err := xlsx.ParseAll(file, &users); err != nil {
		apiresp.GinError(c, errs.ErrArgs.WrapMsg("xlsx file parse error "+err.Error()))
		return
	}
	us, err := o.xlsx2user(users)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	imToken, err := o.imApiCaller.ImAdminTokenWithDefaultAdmin(c)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}

	ctx := mctx.WithAdminUser(mctx.WithApiToken(c, imToken))
	apiresp.GinError(c, o.registerChatUser(ctx, ip, us))
}

func (o *AdminApi) ImportUserByJson(c *gin.Context) {
	var req struct {
		Secret string                   `json:"secret"`
		Users  []*chat.RegisterUserInfo `json:"users"`
	}
	if err := c.BindJSON(&req); err != nil {
		apiresp.GinError(c, err)
		return
	}
	ip, err := o.getClientIP(c)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	if err := o.checkSecretAdmin(c, req.Secret); err != nil {
		apiresp.GinError(c, err)
		return
	}
	imToken, err := o.imApiCaller.ImAdminTokenWithDefaultAdmin(c)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	ctx := mctx.WithAdminUser(mctx.WithApiToken(c, imToken))
	apiresp.GinError(c, o.registerChatUser(ctx, ip, req.Users))
}

func (o *AdminApi) xlsx2user(users []model.User) ([]*chat.RegisterUserInfo, error) {
	chatUsers := make([]*chat.RegisterUserInfo, len(users))
	for i, info := range users {
		if info.Nickname == "" {
			return nil, errs.ErrArgs.WrapMsg("nickname is empty")
		}
		if info.AreaCode == "" || info.PhoneNumber == "" {
			return nil, errs.ErrArgs.WrapMsg("areaCode or phoneNumber is empty")
		}
		if info.Password == "" {
			return nil, errs.ErrArgs.WrapMsg("password is empty")
		}
		if !strings.HasPrefix(info.AreaCode, "+") {
			return nil, errs.ErrArgs.WrapMsg("areaCode format error")
		}
		if _, err := strconv.ParseUint(info.AreaCode[1:], 10, 16); err != nil {
			return nil, errs.ErrArgs.WrapMsg("areaCode format error")
		}
		gender, _ := strconv.Atoi(info.Gender)
		chatUsers[i] = &chat.RegisterUserInfo{
			UserID:      info.UserID,
			Nickname:    info.Nickname,
			FaceURL:     info.FaceURL,
			Birth:       o.xlsxBirth(info.Birth).UnixMilli(),
			Gender:      int32(gender),
			AreaCode:    info.AreaCode,
			PhoneNumber: info.PhoneNumber,
			Email:       info.Email,
			Account:     info.Account,
			Password:    utils.Md5(info.Password),
		}
	}
	return chatUsers, nil
}

func (o *AdminApi) xlsxBirth(s string) time.Time {
	if s == "" {
		return time.Now()
	}
	var separator byte
	for _, b := range []byte(s) {
		if b < '0' || b > '9' {
			separator = b
		}
	}
	arr := strings.Split(s, string([]byte{separator}))
	if len(arr) != 3 {
		return time.Now()
	}
	year, _ := strconv.Atoi(arr[0])
	month, _ := strconv.Atoi(arr[1])
	day, _ := strconv.Atoi(arr[2])
	t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
	if t.Before(time.Date(1900, 0, 0, 0, 0, 0, 0, time.Local)) {
		return time.Now()
	}
	return t
}

func (o *AdminApi) registerChatUser(ctx context.Context, ip string, users []*chat.RegisterUserInfo) error {
	if len(users) == 0 {
		return errs.ErrArgs.WrapMsg("users is empty")
	}
	for _, info := range users {
		respRegisterUser, err := o.chatClient.RegisterUser(ctx, &chat.RegisterUserReq{Ip: ip, User: info, Platform: constant.AdminPlatformID})
		if err != nil {
			return err
		}
		userInfo := &sdkws.UserInfo{
			UserID:   respRegisterUser.UserID,
			Nickname: info.Nickname,
			FaceURL:  info.FaceURL,
		}
		if err = o.imApiCaller.RegisterUser(ctx, []*sdkws.UserInfo{userInfo}); err != nil {
			return err
		}
		if resp, err := o.adminClient.FindDefaultFriend(ctx, &admin.FindDefaultFriendReq{}); err == nil {
			_ = o.imApiCaller.ImportFriend(ctx, respRegisterUser.UserID, resp.UserIDs)
		}
		if resp, err := o.adminClient.FindDefaultGroup(ctx, &admin.FindDefaultGroupReq{}); err == nil {
			_ = o.imApiCaller.InviteToGroup(ctx, respRegisterUser.UserID, resp.GroupIDs)
		}
	}
	return nil
}

func (o *AdminApi) BatchImportTemplate(c *gin.Context) {
	md5Sum := md5.Sum(config.ImportTemplate)
	md5Val := hex.EncodeToString(md5Sum[:])
	if c.GetHeader("If-None-Match") == md5Val {
		c.Status(http.StatusNotModified)
		return
	}
	c.Header("Content-Disposition", "attachment; filename=template.xlsx")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Length", strconv.Itoa(len(config.ImportTemplate)))
	c.Header("ETag", md5Val)
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", config.ImportTemplate)
}
