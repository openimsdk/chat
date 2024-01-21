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
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/OpenIMSDK/protocol/msg"
	"github.com/OpenIMSDK/tools/utils"
	"io"
	"net"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/OpenIMSDK/chat/pkg/common/apicall"
	"github.com/OpenIMSDK/chat/pkg/common/apistruct"
	constant2 "github.com/OpenIMSDK/chat/pkg/common/constant"
	"github.com/OpenIMSDK/chat/pkg/common/mctx"
	"github.com/OpenIMSDK/protocol/sdkws"
	"github.com/OpenIMSDK/tools/checker"
	"github.com/OpenIMSDK/tools/log"

	"github.com/OpenIMSDK/protocol/constant"
	"github.com/OpenIMSDK/tools/a2r"
	"github.com/OpenIMSDK/tools/apiresp"
	"github.com/OpenIMSDK/tools/errs"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"

	"github.com/OpenIMSDK/chat/pkg/common/config"
	"github.com/OpenIMSDK/chat/pkg/proto/admin"
	"github.com/OpenIMSDK/chat/pkg/proto/chat"
)

func NewChat(chatConn, adminConn grpc.ClientConnInterface) *ChatApi {
	return &ChatApi{chatClient: chat.NewChatClient(chatConn), adminClient: admin.NewAdminClient(adminConn), imApiCaller: apicall.NewCallerInterface()}
}

type ChatApi struct {
	chatClient  chat.ChatClient
	adminClient admin.AdminClient
	imApiCaller apicall.CallerInterface
}

// ################## ACCOUNT ##################

func (o *ChatApi) SendVerifyCode(c *gin.Context) {
	req := chat.SendVerifyCodeReq{}

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

func (o *ChatApi) VerifyCode(c *gin.Context) {
	a2r.Call(chat.ChatClient.VerifyCode, o.chatClient, c)
}

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
		apiresp.GinError(c, err) // 参数校验失败
		return
	}
	ip, err := o.getClientIP(c)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	req.Ip = ip
	respRegisterUser, err := o.chatClient.RegisterUser(c, &req)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	userInfo := &sdkws.UserInfo{
		UserID:     respRegisterUser.UserID,
		Nickname:   req.User.Nickname,
		FaceURL:    req.User.FaceURL,
		CreateTime: time.Now().UnixMilli(),
	}
	err = o.imApiCaller.RegisterUser(c, []*sdkws.UserInfo{userInfo})
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	imToken, err := o.imApiCaller.ImAdminTokenWithDefaultAdmin(c)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiCtx := mctx.WithApiToken(c, imToken)
	rpcCtx := mctx.WithAdminUser(c)
	if resp, err := o.adminClient.FindDefaultFriend(rpcCtx, &admin.FindDefaultFriendReq{}); err == nil {
		_ = o.imApiCaller.ImportFriend(apiCtx, respRegisterUser.UserID, resp.UserIDs)
	}
	if resp, err := o.adminClient.FindDefaultGroup(rpcCtx, &admin.FindDefaultGroupReq{}); err == nil {
		_ = o.imApiCaller.InviteToGroup(apiCtx, respRegisterUser.UserID, resp.GroupIDs)
	}
	if req.AutoLogin {
		resp.ImToken, err = o.imApiCaller.UserToken(c, respRegisterUser.UserID, req.Platform)
		if err != nil {
			apiresp.GinError(c, err)
			return
		}
	}
	resp.ChatToken = respRegisterUser.ChatToken
	resp.UserID = respRegisterUser.UserID
	log.ZInfo(c, "registerUser api", "resp", &resp)
	apiresp.GinSuccess(c, &resp)
}

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
		apiresp.GinError(c, err) // 参数校验失败
		return
	}
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
	imToken, err := o.imApiCaller.UserToken(c, resp1.UserID, req.Platform)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	resp.ImToken = imToken
	resp.UserID = resp1.UserID
	resp.ChatToken = resp1.ChatToken
	apiresp.GinSuccess(c, resp)
}

func (o *ChatApi) ResetPassword(c *gin.Context) {
	a2r.Call(chat.ChatClient.ResetPassword, o.chatClient, c)
}

func (o *ChatApi) ChangePassword(c *gin.Context) {
	a2r.Call(chat.ChatClient.ChangePassword, o.chatClient, c)
}

// ################## USER ##################

func (o *ChatApi) UpdateUserInfo(c *gin.Context) {
	var (
		req  chat.UpdateUserInfoReq
		resp apistruct.UpdateUserInfoResp
	)
	if err := c.BindJSON(&req); err != nil {
		apiresp.GinError(c, err)
		return
	}
	log.ZInfo(c, "updateUserInfo", "req", &req)
	if err := checker.Validate(&req); err != nil {
		apiresp.GinError(c, err) // 参数校验失败
		return
	}
	respUpdate, err := o.chatClient.UpdateUserInfo(c, &req)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	opUserType, err := mctx.GetUserType(c)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	var imToken string
	if opUserType == constant2.NormalUser {
		imToken, err = o.imApiCaller.ImAdminTokenWithDefaultAdmin(c)
	} else if opUserType == constant2.AdminUser {
		imToken, err = o.imApiCaller.UserToken(c, config.GetIMAdmin(mctx.GetOpUserID(c)), constant.AdminPlatformID)
	} else {
		apiresp.GinError(c, errs.ErrArgs.Wrap("opUserType unknown"))
		return
	}
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	var (
		nickName string
		faceURL  string
	)
	if req.Nickname != nil {
		nickName = req.Nickname.Value
	} else {
		nickName = respUpdate.NickName
	}
	if req.FaceURL != nil {
		faceURL = req.FaceURL.Value
	} else {
		faceURL = respUpdate.FaceUrl
	}
	err = o.imApiCaller.UpdateUserInfo(mctx.WithApiToken(c, imToken), req.UserID, nickName, faceURL)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiresp.GinSuccess(c, resp)
}

func (o *ChatApi) FindUserPublicInfo(c *gin.Context) {
	a2r.Call(chat.ChatClient.FindUserPublicInfo, o.chatClient, c)
}

func (o *ChatApi) FindUserFullInfo(c *gin.Context) {
	a2r.Call(chat.ChatClient.FindUserFullInfo, o.chatClient, c)
}

//func (o *ChatApi) GetUsersFullInfo(c *gin.Context) {
//	a2r.Call(chat.ChatClient.GetUsersFullInfo, o.chatClient, c)
//}

func (o *ChatApi) SearchUserFullInfo(c *gin.Context) {
	a2r.Call(chat.ChatClient.SearchUserFullInfo, o.chatClient, c)
}

func (o *ChatApi) SearchUserPublicInfo(c *gin.Context) {
	a2r.Call(chat.ChatClient.SearchUserPublicInfo, o.chatClient, c)
}

// ################## APPLET ##################

func (o *ChatApi) FindApplet(c *gin.Context) {
	a2r.Call(admin.AdminClient.FindApplet, o.adminClient, c)
}

// ################## CONFIG ##################

func (o *ChatApi) GetClientConfig(c *gin.Context) {
	a2r.Call(admin.AdminClient.GetClientConfig, o.adminClient, c)
}

// ################## CALLBACK ##################

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

func (o *ChatApi) UploadLogs(c *gin.Context) {
	a2r.Call(chat.ChatClient.UploadLogs, o.chatClient, c)
}

func (o *ChatApi) SearchFriend(c *gin.Context) {
	var req struct {
		UserID string `json:"userID"`
		chat.SearchUserInfoReq
	}
	if err := c.BindJSON(&req); err != nil {
		apiresp.GinError(c, err)
		return
	}
	if req.UserID == "" {
		req.UserID = mctx.GetOpUserID(c)
	}
	imToken, err := o.imApiCaller.ImAdminTokenWithDefaultAdmin(c)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	userIDs, err := o.imApiCaller.FriendUserIDs(mctx.WithApiToken(c, imToken), req.UserID)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	if len(userIDs) == 0 {
		apiresp.GinSuccess(c, &chat.SearchUserInfoResp{})
		return
	}
	req.SearchUserInfoReq.UserIDs = userIDs
	resp, err := o.chatClient.SearchUserInfo(c, &req.SearchUserInfoReq)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiresp.GinSuccess(c, resp)
}

func (m *ChatApi) CallbackExample(c *gin.Context) {

	// 1. Callback after sending a single chat message
	var req apistruct.CallbackAfterSendSingleMsgReq

	if err := c.BindJSON(&req); err != nil {
		log.ZError(c, "CallbackExample BindJSON failed", err)
		apiresp.GinError(c, errs.ErrArgs.WithDetail(err.Error()).Wrap())
		return
	}

	resp := apistruct.CallbackAfterSendSingleMsgResp{
		CommonCallbackResp: apistruct.CommonCallbackResp{
			ActionCode: 0,
			ErrCode:    200,
			ErrMsg:     "success",
			ErrDlt:     "successful",
			NextCode:   0,
		},
	}
	c.JSON(http.StatusOK, resp)

	// 2. If the user receiving the message is a customer service bot, return the message.

	// UserID of the robot account

	if req.SendID == "robotics" || req.RecvID != "robotics" {
		return
	}

	if req.ContentType != constant.Picture && req.ContentType != constant.Text {
		return
	}

	// Administrator token
	url := "http://127.0.0.1:10009/account/login"
	adminID := config.Config.ChatAdmin[0].AdminID
	paswd := md5.Sum([]byte(adminID))

	admin_input := admin.LoginReq{
		Account:  config.Config.ChatAdmin[0].AdminID,
		Password: hex.EncodeToString(paswd[:]),
	}

	header := make(map[string]string, 2)
	header["operationID"] = "111"

	b, err := Post(c, url, header, admin_input, 10)
	if err != nil {
		log.ZError(c, "CallbackExample send message failed", err)
		apiresp.GinError(c, errs.ErrInternalServer.WithDetail(err.Error()).Wrap())
		return
	}

	type TokenInfo struct {
		ErrCode int                      `json:"errCode"`
		ErrMsg  string                   `json:"errMsg"`
		ErrDlt  string                   `json:"errDlt"`
		Data    apistruct.AdminLoginResp `json:"data,omitempty"`
	}

	admin_output := &TokenInfo{}

	if err = json.Unmarshal(b, admin_output); err != nil {
		log.ZError(c, "CallbackExample unmarshal failed", err)
		apiresp.GinError(c, errs.ErrInternalServer.WithDetail(err.Error()).Wrap())
		return
	}

	header["token"] = admin_output.Data.AdminToken

	url = "http://127.0.0.1:10008/user/find/public"

	search_input := chat.FindUserFullInfoReq{
		UserIDs: []string{"robotics"},
	}

	b, err = Post(c, url, header, search_input, 10)
	if err != nil {
		log.ZError(c, "CallbackExample unmarshal failed", err)
		apiresp.GinError(c, errs.ErrInternalServer.WithDetail(err.Error()).Wrap())
		return
	}

	type UserInfo struct {
		ErrCode int                       `json:"errCode"`
		ErrMsg  string                    `json:"errMsg"`
		ErrDlt  string                    `json:"errDlt"`
		Data    chat.FindUserFullInfoResp `json:"data,omitempty"`
	}

	search_output := &UserInfo{}

	if err = json.Unmarshal(b, search_output); err != nil {
		log.ZError(c, "search_output unmarshal failed", err)
		apiresp.GinError(c, errs.ErrInternalServer.WithDetail(err.Error()).Wrap())
		return
	}

	if len(search_output.Data.Users) == 0 {
		apiresp.GinError(c, errs.ErrRecordNotFound.Wrap("the robotics not found"))
		return
	}

	log.ZDebug(c, "callback", "searchUserAccount", search_output)

	text := apistruct.TextElem{}
	picture := apistruct.PictureElem{}
	mapStruct := make(map[string]any)
	// Processing text messages

	if err != nil {
		log.ZError(c, "CallbackExample get Sender failed", err)
		apiresp.GinError(c, errs.ErrInternalServer.WithDetail(err.Error()).Wrap())
		return
	}

	// Handle message structures
	if req.ContentType == constant.Text {
		err = json.Unmarshal([]byte(req.Content), &text)
		if err != nil {
			log.ZError(c, "CallbackExample unmarshal failed", err)
			apiresp.GinError(c, errs.ErrInternalServer.WithDetail(err.Error()).Wrap())
			return
		}
		log.ZDebug(c, "callback", "text", text)
		mapStruct["content"] = text.Content
	} else {
		err = json.Unmarshal([]byte(req.Content), &picture)
		if err != nil {
			log.ZError(c, "CallbackExample unmarshal failed", err)
			apiresp.GinError(c, errs.ErrInternalServer.WithDetail(err.Error()).Wrap())
			return
		}
		log.ZDebug(c, "callback", "text", picture)
		if strings.Contains(picture.SourcePicture.Type, "/") {
			arr := strings.Split(picture.SourcePicture.Type, "/")
			picture.SourcePicture.Type = arr[1]
		}

		if strings.Contains(picture.BigPicture.Type, "/") {
			arr := strings.Split(picture.BigPicture.Type, "/")
			picture.BigPicture.Type = arr[1]
		}

		if len(picture.SnapshotPicture.Type) == 0 {
			picture.SnapshotPicture.Type = picture.SourcePicture.Type
		}

		mapStructSnap := make(map[string]interface{})
		if mapStructSnap, err = convertStructToMap(picture.SnapshotPicture); err != nil {
			log.ZError(c, "CallbackExample struct to map failed", err)
			apiresp.GinError(c, err)
			return
		}
		mapStruct["snapshotPicture"] = mapStructSnap

		mapStructBig := make(map[string]interface{})
		if mapStructBig, err = convertStructToMap(picture.BigPicture); err != nil {
			log.ZError(c, "CallbackExample struct to map failed", err)
			apiresp.GinError(c, err)
			return
		}
		mapStruct["bigPicture"] = mapStructBig

		mapStructSource := make(map[string]interface{})
		if mapStructSource, err = convertStructToMap(picture.SourcePicture); err != nil {
			log.ZError(c, "CallbackExample struct to map failed", err)
			apiresp.GinError(c, err)
			return
		}
		mapStruct["sourcePicture"] = mapStructSource
		mapStruct["sourcePath"] = picture.SourcePath
	}

	log.ZDebug(c, "callback", "mapStruct", mapStruct, "mapStructSnap")
	header["token"] = admin_output.Data.ImToken

	input := &apistruct.SendMsgReq{
		RecvID: req.SendID,
		SendMsg: apistruct.SendMsg{
			SendID:           search_output.Data.Users[0].UserID,
			SenderNickname:   search_output.Data.Users[0].Nickname,
			SenderFaceURL:    search_output.Data.Users[0].FaceURL,
			SenderPlatformID: req.SenderPlatformID,
			Content:          mapStruct,
			ContentType:      req.ContentType,
			SessionType:      req.SessionType,
			SendTime:         utils.GetCurrentTimestampByMill(), // millisecond
		},
	}

	url = "http://127.0.0.1:10002/msg/send_msg"

	type sendResp struct {
		ErrCode int             `json:"errCode"`
		ErrMsg  string          `json:"errMsg"`
		ErrDlt  string          `json:"errDlt"`
		Data    msg.SendMsgResp `json:"data,omitempty"`
	}

	output := &sendResp{}

	// Initiate a post request that calls the interface that sends the message (the bot sends a message to user)
	b, err = Post(c, url, header, input, 10)
	if err != nil {
		log.ZError(c, "CallbackExample send message failed", err)
		apiresp.GinError(c, errs.ErrInternalServer.WithDetail(err.Error()).Wrap())
		return
	}
	if err = json.Unmarshal(b, output); err != nil {
		log.ZError(c, "CallbackExample unmarshal failed", err)
		apiresp.GinError(c, errs.ErrInternalServer.WithDetail(err.Error()).Wrap())
		return
	}
	res := &msg.SendMsgResp{
		ServerMsgID: output.Data.ServerMsgID,
		ClientMsgID: output.Data.ClientMsgID,
		SendTime:    output.Data.SendTime,
	}

	apiresp.GinSuccess(c, res)
}

// struct to map
func convertStructToMap(input interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	inputType := reflect.TypeOf(input)

	inputValue := reflect.ValueOf(input)

	if inputType.Kind() != reflect.Struct {
		return nil, errs.ErrArgs.Wrap("input is not a struct")
	}

	for i := 0; i < inputType.NumField(); i++ {
		field := inputType.Field(i)
		fieldValue := inputValue.Field(i)

		mapKey := field.Tag.Get("mapstructure")
		fmt.Println(mapKey)

		if mapKey == "" {
			mapKey = field.Name
		}

		mapKey = strings.ToLower(mapKey)

		result[mapKey] = fieldValue.Interface()
	}

	return result, nil
}

func Post(ctx context.Context, url string, header map[string]string, data any, timeout int) (content []byte, err error) {
	var (
		// define http client.
		client = &http.Client{
			Timeout: 15 * time.Second, // max timeout is 15s
		}
	)

	if timeout > 0 {
		var cancel func()
		ctx, cancel = context.WithTimeout(ctx, time.Second*time.Duration(timeout))
		defer cancel()
	}

	jsonStr, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, err
	}

	if operationID, _ := ctx.Value(constant.OperationID).(string); operationID != "" {
		req.Header.Set(constant.OperationID, operationID)
	}
	for k, v := range header {
		req.Header.Set(k, v)
	}
	req.Header.Add("content-type", "application/json; charset=utf-8")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return result, nil
}
