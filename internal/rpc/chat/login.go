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

package chat

import (
	"context"
	"fmt"
	"github.com/openimsdk/tools/utils/datautil"
	"math/rand"
	"strconv"
	"strings"
	"time"

	constant2 "github.com/openimsdk/protocol/constant"

	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/mcontext"

	"github.com/openimsdk/chat/pkg/common/constant"
	"github.com/openimsdk/chat/pkg/common/db/dbutil"
	chat2 "github.com/openimsdk/chat/pkg/common/db/table/chat"
	"github.com/openimsdk/chat/pkg/eerrs"
	"github.com/openimsdk/chat/pkg/protocol/chat"
)

func (o *chatSvr) verifyCodeJoin(areaCode, phoneNumber string) string {
	return areaCode + " " + phoneNumber
}

func (o *chatSvr) SendVerifyCode(ctx context.Context, req *chat.SendVerifyCodeReq) (*chat.SendVerifyCodeResp, error) {
	switch int(req.UsedFor) {
	case constant.VerificationCodeForRegister:
		if err := o.Admin.CheckRegister(ctx, req.Ip); err != nil {
			return nil, err
		}
		if req.Email == "" {
			if req.AreaCode == "" || req.PhoneNumber == "" {
				return nil, errs.ErrArgs.WrapMsg("area code or phone number is empty")
			}
			if req.AreaCode[0] != '+' {
				req.AreaCode = "+" + req.AreaCode
			}
			if _, err := strconv.ParseUint(req.AreaCode[1:], 10, 64); err != nil {
				return nil, errs.ErrArgs.WrapMsg("area code must be number")
			}
			if _, err := strconv.ParseUint(req.PhoneNumber, 10, 64); err != nil {
				return nil, errs.ErrArgs.WrapMsg("phone number must be number")
			}
			_, err := o.Database.TakeAttributeByPhone(ctx, req.AreaCode, req.PhoneNumber)
			if err == nil {
				return nil, eerrs.ErrPhoneAlreadyRegister.WrapMsg("phone already register")
			} else if !dbutil.IsDBNotFound(err) {
				return nil, err
			}
		} else {
			if err := chat.EmailCheck(req.Email); err != nil {
				return nil, errs.ErrArgs.WrapMsg("email must be right")
			}
			_, err := o.Database.TakeAttributeByEmail(ctx, req.Email)
			if err == nil {
				return nil, eerrs.ErrEmailAlreadyRegister.WrapMsg("email already register")
			} else if !dbutil.IsDBNotFound(err) {
				return nil, err
			}
		}
		conf, err := o.Admin.GetConfig(ctx)
		if err != nil {
			return nil, err
		}
		if val := conf[constant.NeedInvitationCodeRegisterConfigKey]; datautil.Contain(strings.ToLower(val), "1", "true", "yes") {
			if req.InvitationCode == "" {
				return nil, errs.ErrArgs.WrapMsg("invitation code is empty")
			}
			if err := o.Admin.CheckInvitationCode(ctx, req.InvitationCode); err != nil {
				return nil, err
			}
		}
	case constant.VerificationCodeForLogin, constant.VerificationCodeForResetPassword:
		if req.Email == "" {
			_, err := o.Database.TakeAttributeByPhone(ctx, req.AreaCode, req.PhoneNumber)
			if dbutil.IsDBNotFound(err) {
				return nil, eerrs.ErrAccountNotFound.WrapMsg("phone unregistered")
			} else if err != nil {
				return nil, err
			}
		} else {
			_, err := o.Database.TakeAttributeByEmail(ctx, req.Email)
			if dbutil.IsDBNotFound(err) {
				return nil, eerrs.ErrAccountNotFound.WrapMsg("email unregistered")
			} else if err != nil {
				return nil, err
			}
		}

	default:
		return nil, errs.ErrArgs.WrapMsg("used unknown")
	}
	if o.SMS == nil && o.Mail == nil {
		return &chat.SendVerifyCodeResp{}, nil // super code
	}
	isEmail := req.Email != ""
	var (
		code     = o.genVerifyCode()
		account  string
		sendCode func() error
	)
	if isEmail {
		if o.Mail == nil {
			return nil, errs.ErrInternalServer.WrapMsg("email verification code is not enabled")
		}
		sendCode = func() error {
			return o.Mail.SendMail(ctx, req.Email, code)
		}
		account = req.Email
	} else {
		if o.SMS == nil {
			return nil, errs.ErrInternalServer.WrapMsg("mobile phone verification code is not enabled")
		}
		sendCode = func() error {
			return o.SMS.SendCode(ctx, req.AreaCode, req.PhoneNumber, code)
		}
		account = o.verifyCodeJoin(req.AreaCode, req.PhoneNumber)
	}
	now := time.Now()
	count, err := o.Database.CountVerifyCodeRange(ctx, account, now.Add(-o.Code.UintTime), now)
	if err != nil {
		return nil, err
	}
	if o.Code.MaxCount < int(count) {
		return nil, eerrs.ErrVerifyCodeSendFrequently.Wrap()
	}
	platformName := constant2.PlatformIDToName(int(req.Platform))
	if platformName == "" {
		platformName = fmt.Sprintf("platform:%d", req.Platform)
	}
	vc := &chat2.VerifyCode{
		Account:    account,
		Code:       code,
		Platform:   platformName,
		Duration:   uint(o.Code.ValidTime / time.Second),
		Count:      0,
		Used:       false,
		CreateTime: now,
	}
	if err := o.Database.AddVerifyCode(ctx, vc, sendCode); err != nil {
		return nil, err
	}
	log.ZDebug(ctx, "send code success", "account", account, "code", code, "platform", platformName)
	return &chat.SendVerifyCodeResp{}, nil
}

func (o *chatSvr) verifyCode(ctx context.Context, account string, verifyCode string) (string, error) {
	if verifyCode == "" {
		return "", errs.ErrArgs.WrapMsg("verify code is empty")
	}
	if o.SMS == nil && o.Mail == nil {
		if o.Code.SuperCode != verifyCode {
			return "", eerrs.ErrVerifyCodeNotMatch.Wrap()
		}
		return "", nil
	}
	last, err := o.Database.TakeLastVerifyCode(ctx, account)
	if err != nil {
		if dbutil.IsDBNotFound(err) {
			return "", eerrs.ErrVerifyCodeExpired.Wrap()
		}
		return "", err
	}
	if last.CreateTime.Unix()+int64(last.Duration) < time.Now().Unix() {
		return last.ID, eerrs.ErrVerifyCodeExpired.Wrap()
	}
	if last.Used {
		return last.ID, eerrs.ErrVerifyCodeUsed.Wrap()
	}
	if n := o.Code.ValidCount; n > 0 {
		if last.Count >= n {
			return last.ID, eerrs.ErrVerifyCodeMaxCount.Wrap()
		}
		if last.Code != verifyCode {
			if err := o.Database.UpdateVerifyCodeIncrCount(ctx, last.ID); err != nil {
				return last.ID, err
			}
		}
	}
	if last.Code != verifyCode {
		return last.ID, eerrs.ErrVerifyCodeNotMatch.Wrap()
	}
	return last.ID, nil
}

func (o *chatSvr) VerifyCode(ctx context.Context, req *chat.VerifyCodeReq) (*chat.VerifyCodeResp, error) {
	var account string
	if req.PhoneNumber != "" {
		account = o.verifyCodeJoin(req.AreaCode, req.PhoneNumber)
	} else {
		account = req.Email
	}
	if _, err := o.verifyCode(ctx, account, req.VerifyCode); err != nil {
		return nil, err
	}
	return &chat.VerifyCodeResp{}, nil
}

func (o *chatSvr) genUserID() string {
	const l = 10
	data := make([]byte, l)
	rand.Read(data)
	chars := []byte("0123456789")
	for i := 0; i < len(data); i++ {
		if i == 0 {
			data[i] = chars[1:][data[i]%9]
		} else {
			data[i] = chars[data[i]%10]
		}
	}
	return string(data)
}

func (o *chatSvr) genVerifyCode() string {
	data := make([]byte, o.Code.Len)
	rand.Read(data)
	chars := []byte("0123456789")
	for i := 0; i < len(data); i++ {
		data[i] = chars[data[i]%10]
	}
	return string(data)
}

func (o *chatSvr) RegisterUser(ctx context.Context, req *chat.RegisterUserReq) (*chat.RegisterUserResp, error) {
	resp := &chat.RegisterUserResp{}

	isAdmin, err := o.Admin.CheckNilOrAdmin(ctx)
	ctx = o.WithAdminUser(ctx)
	if err != nil {
		return nil, err
	}
	if req.User == nil {
		return nil, errs.ErrArgs.WrapMsg("user is nil")
	}
	if req.User.Email == "" {
		if (req.User.AreaCode == "" && req.User.PhoneNumber != "") || (req.User.AreaCode != "" && req.User.PhoneNumber == "") {
			return nil, errs.ErrArgs.WrapMsg("area code or phone number error, no email provide")
		}
	}
	var usedInvitationCode bool
	if !isAdmin {
		if req.User.UserID != "" {
			return nil, errs.ErrNoPermission.WrapMsg("only admin can set user id")
		}
		if err := o.Admin.CheckRegister(ctx, req.Ip); err != nil {
			return nil, err
		}
		conf, err := o.Admin.GetConfig(ctx)
		if err != nil {
			return nil, err
		}
		if val := conf[constant.NeedInvitationCodeRegisterConfigKey]; datautil.Contain(strings.ToLower(val), "1", "true", "yes") {
			usedInvitationCode = true
			if req.InvitationCode == "" {
				return nil, errs.ErrArgs.WrapMsg("invitation code is empty")
			}
			if err := o.Admin.CheckInvitationCode(ctx, req.InvitationCode); err != nil {
				return nil, err
			}
		}
		if req.User.Email == "" {
			if _, err := o.verifyCode(ctx, o.verifyCodeJoin(req.User.AreaCode, req.User.PhoneNumber), req.VerifyCode); err != nil {
				return nil, err
			}
		} else {
			if _, err := o.verifyCode(ctx, req.User.Email, req.VerifyCode); err != nil {
				return nil, err
			}
		}

	}
	if req.User.UserID == "" {
		for i := 0; i < 20; i++ {
			userID := o.genUserID()
			_, err := o.Database.GetUser(ctx, userID)
			if err == nil {
				continue
			} else if dbutil.IsDBNotFound(err) {
				req.User.UserID = userID
				break
			} else {
				return nil, err
			}
		}
		if req.User.UserID == "" {
			return nil, errs.ErrInternalServer.WrapMsg("gen user id failed")
		}
	} else {
		_, err := o.Database.GetUser(ctx, req.User.UserID)
		if err == nil {
			return nil, errs.ErrArgs.WrapMsg("appoint user id already register")
		} else if !dbutil.IsDBNotFound(err) {
			return nil, err
		}
	}
	var registerType int32
	if req.User.PhoneNumber != "" {
		if req.User.AreaCode[0] != '+' {
			req.User.AreaCode = "+" + req.User.AreaCode
		}
		if _, err := strconv.ParseUint(req.User.AreaCode[1:], 10, 64); err != nil {
			return nil, errs.ErrArgs.WrapMsg("area code must be number")
		}
		if _, err := strconv.ParseUint(req.User.PhoneNumber, 10, 64); err != nil {
			return nil, errs.ErrArgs.WrapMsg("phone number must be number")
		}
		_, err := o.Database.TakeAttributeByPhone(ctx, req.User.AreaCode, req.User.PhoneNumber)
		if err == nil {
			return nil, eerrs.ErrPhoneAlreadyRegister.Wrap()
		} else if !dbutil.IsDBNotFound(err) {
			return nil, err
		}
		registerType = constant.PhoneRegister
	}

	if req.User.Account != "" {
		_, err := o.Database.TakeAttributeByAccount(ctx, req.User.Account)
		if err == nil {
			return nil, eerrs.ErrAccountAlreadyRegister.Wrap()
		} else if !dbutil.IsDBNotFound(err) {
			return nil, err
		}
	}

	if req.User.Email != "" {
		_, err := o.Database.TakeAttributeByEmail(ctx, req.User.Email)
		registerType = constant.EmailRegister
		if err == nil {
			return nil, eerrs.ErrEmailAlreadyRegister.Wrap()
		} else if !dbutil.IsDBNotFound(err) {
			return nil, err
		}
	}
	register := &chat2.Register{
		UserID:      req.User.UserID,
		DeviceID:    req.DeviceID,
		IP:          req.Ip,
		Platform:    constant2.PlatformID2Name[int(req.Platform)],
		AccountType: "",
		Mode:        constant.UserMode,
		CreateTime:  time.Now(),
	}
	account := &chat2.Account{
		UserID:         req.User.UserID,
		Password:       req.User.Password,
		OperatorUserID: mcontext.GetOpUserID(ctx),
		ChangeTime:     register.CreateTime,
		CreateTime:     register.CreateTime,
	}
	attribute := &chat2.Attribute{
		UserID:         req.User.UserID,
		Account:        req.User.Account,
		PhoneNumber:    req.User.PhoneNumber,
		AreaCode:       req.User.AreaCode,
		Email:          req.User.Email,
		Nickname:       req.User.Nickname,
		FaceURL:        req.User.FaceURL,
		Gender:         req.User.Gender,
		BirthTime:      time.UnixMilli(req.User.Birth),
		ChangeTime:     register.CreateTime,
		CreateTime:     register.CreateTime,
		AllowVibration: constant.DefaultAllowVibration,
		AllowBeep:      constant.DefaultAllowBeep,
		AllowAddFriend: constant.DefaultAllowAddFriend,
		RegisterType:   registerType,
	}
	if err := o.Database.RegisterUser(ctx, register, account, attribute); err != nil {
		return nil, err
	}
	if usedInvitationCode {
		if err := o.Admin.UseInvitationCode(ctx, req.User.UserID, req.InvitationCode); err != nil {
			log.ZError(ctx, "UseInvitationCode", err, "userID", req.User.UserID, "invitationCode", req.InvitationCode)
		}
	}
	if req.AutoLogin {
		chatToken, err := o.Admin.CreateToken(ctx, req.User.UserID, constant.NormalUser)
		if err == nil {
			resp.ChatToken = chatToken.Token
		} else {
			log.ZError(ctx, "Admin CreateToken Failed", err, "userID", req.User.UserID, "platform", req.Platform)
		}
	}
	resp.UserID = req.User.UserID
	return resp, nil
}

func (o *chatSvr) Login(ctx context.Context, req *chat.LoginReq) (*chat.LoginResp, error) {
	resp := &chat.LoginResp{}
	if req.Password == "" && req.VerifyCode == "" {
		return nil, errs.ErrArgs.WrapMsg("password or code must be set")
	}
	var err error
	var attribute *chat2.Attribute
	if req.Account != "" {
		attribute, err = o.Database.GetAttributeByAccount(ctx, req.Account)
	} else if req.PhoneNumber != "" {
		if req.AreaCode == "" {
			return nil, errs.ErrArgs.WrapMsg("area code must")
		}
		if req.AreaCode[0] != '+' {
			req.AreaCode = "+" + req.AreaCode
		}
		if _, err := strconv.ParseUint(req.AreaCode[1:], 10, 64); err != nil {
			return nil, errs.ErrArgs.WrapMsg("area code must be number")
		}
		attribute, err = o.Database.GetAttributeByPhone(ctx, req.AreaCode, req.PhoneNumber)
	} else if req.Email != "" {
		attribute, err = o.Database.GetAttributeByEmail(ctx, req.Email)
	} else {
		err = errs.ErrArgs.WrapMsg("account or phone number or email must be set")
	}
	if err != nil {
		if dbutil.IsDBNotFound(err) {
			return nil, eerrs.ErrAccountNotFound.WrapMsg("user unregistered")
		}
		return nil, err
	}
	if err := o.Admin.CheckLogin(ctx, attribute.UserID, req.Ip); err != nil {
		return nil, err
	}
	var verifyCodeID *string
	if req.Password == "" {
		var account string
		if req.Email == "" {
			account = o.verifyCodeJoin(req.AreaCode, req.PhoneNumber)
		} else {
			account = req.Email
		}
		id, err := o.verifyCode(ctx, account, req.VerifyCode)
		if err != nil {
			return nil, err
		}
		if id != "" {
			verifyCodeID = &id
		}
	} else {
		account, err := o.Database.GetAccount(ctx, attribute.UserID)
		if err != nil {
			return nil, err
		}
		if account.Password != req.Password {
			return nil, eerrs.ErrPassword.Wrap()
		}
	}
	chatToken, err := o.Admin.CreateToken(ctx, attribute.UserID, constant.NormalUser)
	if err != nil {
		return nil, err
	}
	record := &chat2.UserLoginRecord{
		UserID:    attribute.UserID,
		LoginTime: time.Now(),
		IP:        req.Ip,
		DeviceID:  req.DeviceID,
		Platform:  constant2.PlatformIDToName(int(req.Platform)),
	}
	if err := o.Database.LoginRecord(ctx, record, verifyCodeID); err != nil {
		return nil, err
	}
	if verifyCodeID != nil {
		if err := o.Database.DelVerifyCode(ctx, *verifyCodeID); err != nil {
			return nil, err
		}
	}
	resp.UserID = attribute.UserID
	resp.ChatToken = chatToken.Token
	return resp, nil
}
