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
	"github.com/OpenIMSDK/chat/pkg/common/config"
	constant2 "github.com/OpenIMSDK/protocol/constant"
	"github.com/OpenIMSDK/tools/log"

	"github.com/OpenIMSDK/tools/errs"

	"github.com/OpenIMSDK/chat/pkg/common/constant"
	"github.com/OpenIMSDK/chat/pkg/common/mctx"
	"github.com/OpenIMSDK/chat/pkg/proto/chat"
)

func (o *chatSvr) ResetPassword(ctx context.Context, req *chat.ResetPasswordReq) (*chat.ResetPasswordResp, error) {
	defer log.ZDebug(ctx, "return")
	if req.Password == "" {
		return nil, errs.ErrArgs.Wrap("password must be set")
	}
	var verifyCodeID uint
	var err error
	if req.Email == "" {
		verifyCodeID, err = o.verifyCode(ctx, o.verifyCodeJoin(req.AreaCode, req.PhoneNumber), req.VerifyCode)
	} else {
		verifyCodeID, err = o.verifyCode(ctx, req.Email, req.VerifyCode)
	}

	if err != nil {
		return nil, err
	}

	if req.Email == "" {
		attribute, err := o.Database.GetAttributeByPhone(ctx, req.AreaCode, req.PhoneNumber)
		if err != nil {
			return nil, err
		}
		err = o.Database.UpdatePasswordAndDeleteVerifyCode(ctx, attribute.UserID, req.Password, verifyCodeID)
	} else {
		attribute, err := o.Database.GetAttributeByEmail(ctx, req.Email)
		if err != nil {
			return nil, err
		}
		err = o.Database.UpdatePasswordAndDeleteVerifyCode(ctx, attribute.UserID, req.Password, verifyCodeID)
	}

	if err != nil {
		return nil, err
	}
	return &chat.ResetPasswordResp{}, nil
}

func (o *chatSvr) ChangePassword(ctx context.Context, req *chat.ChangePasswordReq) (*chat.ChangePasswordResp, error) {
	defer log.ZDebug(ctx, "return")
	if req.NewPassword == "" {
		return nil, errs.ErrArgs.Wrap("new password must be set")
	}
	if req.NewPassword == req.CurrentPassword {
		return nil, errs.ErrArgs.Wrap("new password == current password")
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
			return nil, errs.ErrNoPermission.Wrap("no permission change other user password")
		}
	case constant.AdminUser:
		if req.UserID == "" {
			return nil, errs.ErrArgs.Wrap("user id must be set")
		}
	default:
		return nil, errs.ErrInternalServer.Wrap("invalid user type")
	}
	user, err := o.Database.GetUser(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	if userType != constant.AdminUser {
		if user.Password != req.CurrentPassword {
			return nil, errs.ErrNoPermission.Wrap("current password is wrong")
		}
	}
	if user.Password != req.NewPassword {
		if err := o.Database.UpdatePassword(ctx, req.UserID, req.NewPassword); err != nil {
			return nil, err
		}
	}

	imToken, err := o.imApiCaller.UserToken(ctx, config.GetIMAdmin(mctx.GetOpUserID(ctx)), constant2.AdminPlatformID)
	if err != nil {
		return nil, err
	}

	err = o.imApiCaller.ForceOffLine(mctx.WithApiToken(ctx, imToken), req.UserID)
	if err != nil {
		return nil, err
	}

	return &chat.ChangePasswordResp{}, nil
}
