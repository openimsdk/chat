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

import "github.com/OpenIMSDK/tools/errs"

var (
	ErrPassword                 = errs.NewCodeError(20001, "PasswordError")            // 密码错误
	ErrAccountNotFound          = errs.NewCodeError(20002, "AccountNotFound")          // 账号不存在
	ErrPhoneAlreadyRegister     = errs.NewCodeError(20003, "PhoneAlreadyRegister")     // 手机号已经注册
	ErrAccountAlreadyRegister   = errs.NewCodeError(20004, "AccountAlreadyRegister")   // 账号已经注册
	ErrVerifyCodeSendFrequently = errs.NewCodeError(20005, "VerifyCodeSendFrequently") // 频繁获取验证码
	ErrVerifyCodeNotMatch       = errs.NewCodeError(20006, "VerifyCodeNotMatch")       // 验证码错误
	ErrVerifyCodeExpired        = errs.NewCodeError(20007, "VerifyCodeExpired")        // 验证码过期
	ErrVerifyCodeMaxCount       = errs.NewCodeError(20008, "VerifyCodeMaxCount")       // 验证码失败次数过多
	ErrVerifyCodeUsed           = errs.NewCodeError(20009, "VerifyCodeUsed")           // 已经使用
	ErrInvitationCodeUsed       = errs.NewCodeError(20010, "InvitationCodeUsed")       // 邀请码已经使用
	ErrInvitationNotFound       = errs.NewCodeError(20011, "InvitationNotFound")       // 邀请码不存在
	ErrForbidden                = errs.NewCodeError(20012, "Forbidden")                // 限制登录注册
	ErrRefuseFriend             = errs.NewCodeError(20013, "RefuseFriend")             // 拒绝添加好友
	ErrEmailAlreadyRegister     = errs.NewCodeError(20014, "EmailAlreadyRegister")     // 邮箱已经注册
)
