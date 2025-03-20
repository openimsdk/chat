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

	"github.com/openimsdk/tools/errs"

	"github.com/openimsdk/chat/pkg/common/constant"
	"github.com/openimsdk/chat/pkg/common/mctx"
	"github.com/openimsdk/chat/pkg/protocol/chat"
)

func (o *chatSvr) ResetPassword(ctx context.Context, req *chat.ResetPasswordReq) (*chat.ResetPasswordResp, error) {
	if req.Password == "" {
		return nil, errs.ErrArgs.WrapMsg("password must be set")
	}
	if req.AreaCode == "" || req.PhoneNumber == "" {
		if !(req.AreaCode == "" && req.PhoneNumber == "") {
			return nil, errs.ErrArgs.WrapMsg("area code and phone number must set together")
		}
	}
	var verifyCodeID string
	var err error
	if req.Email == "" {
		verifyCodeID, err = o.verifyCode(ctx, o.verifyCodeJoin(req.AreaCode, req.PhoneNumber), req.VerifyCode, phone)
	} else {
		verifyCodeID, err = o.verifyCode(ctx, req.Email, req.VerifyCode, mail)
	}

	if err != nil {
		return nil, err
	}
	var account string
	if req.Email == "" {
		account = BuildCredentialPhone(req.AreaCode, req.PhoneNumber)
	} else {
		account = req.Email
	}
	cred, err := o.Database.TakeCredentialByAccount(ctx, account)
	if err != nil {
		return nil, err
	}
	err = o.Database.UpdatePasswordAndDeleteVerifyCode(ctx, cred.UserID, req.Password, verifyCodeID)
	if err != nil {
		return nil, err
	}
	return &chat.ResetPasswordResp{}, nil
}

func (o *chatSvr) ChangePassword(ctx context.Context, req *chat.ChangePasswordReq) (*chat.ChangePasswordResp, error) {
	if req.NewPassword == "" {
		return nil, errs.ErrArgs.WrapMsg("new password must be set")
	}
	if req.NewPassword == req.CurrentPassword {
		return nil, errs.ErrArgs.WrapMsg("new password == current password")
	}
	opUserID, userType, err := mctx.Check(ctx)
	if err != nil {
		return nil, err
	}
	switch userType {
	case constant.NormalUser:
		if req.UserID == "" {
			req.UserID = opUserID
		}
		if req.UserID != opUserID {
			return nil, errs.ErrNoPermission.WrapMsg("no permission change other user password")
		}
	case constant.AdminUser:
		if req.UserID == "" {
			return nil, errs.ErrArgs.WrapMsg("user id must be set")
		}
	default:
		return nil, errs.ErrInternalServer.WrapMsg("invalid user type")
	}
	user, err := o.Database.GetUser(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	if userType != constant.AdminUser {
		if user.Password != req.CurrentPassword {
			return nil, errs.ErrNoPermission.WrapMsg("current password is wrong")
		}
	}
	if user.Password != req.NewPassword {
		if err := o.Database.UpdatePassword(ctx, req.UserID, req.NewPassword); err != nil {
			return nil, err
		}
	}
	if err := o.Admin.InvalidateToken(ctx, req.UserID); err != nil {
		return nil, err
	}

	return &chat.ChangePasswordResp{}, nil
}
