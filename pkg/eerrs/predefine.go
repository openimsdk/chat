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

package eerrs

import "github.com/OpenIMSDK/Open-IM-Server/pkg/errs"

var (
	ErrPassword                 = errs.NewCodeError(10001, "PasswordError")
	ErrAccountNotFound          = errs.NewCodeError(10002, "AccountNotFound")
	ErrPhoneAlreadyRegister     = errs.NewCodeError(10003, "PhoneAlreadyRegister")
	ErrAccountAlreadyRegister   = errs.NewCodeError(10004, "AccountAlreadyRegister")
	ErrVerifyCodeSendFrequently = errs.NewCodeError(10005, "VerifyCodeSendFrequently") // 频繁获取验证码
	ErrVerifyCodeNotMatch       = errs.NewCodeError(10006, "NotMatch")                 // 验证码错误
	ErrVerifyCodeExpired        = errs.NewCodeError(10007, "expired")                  // 验证码过期
	ErrVerifyCodeMaxCount       = errs.NewCodeError(10008, "Attempts limit")           // 验证码失败次数过多
	ErrVerifyCodeUsed           = errs.NewCodeError(10009, "used")                     // 已经使用

	ErrInvitationCodeUsed = errs.NewCodeError(10010, "InvitationCodeUsed") // 邀请码已经使用
	ErrInvitationNotFound = errs.NewCodeError(10011, "InvitationNotFound") // 邀请码不存在

	ErrForbidden = errs.NewCodeError(10012, "Forbidden")

	ErrRefuseFriend = errs.NewCodeError(10013, "user refused to add as friend") // 拒绝添加好友
)
