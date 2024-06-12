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

import "github.com/openimsdk/tools/errs"

var (
	ErrPassword                 = errs.NewCodeError(20001, "PasswordError")
	ErrAccountNotFound          = errs.NewCodeError(20002, "AccountNotFound")
	ErrPhoneAlreadyRegister     = errs.NewCodeError(20003, "PhoneAlreadyRegister")
	ErrAccountAlreadyRegister   = errs.NewCodeError(20004, "AccountAlreadyRegister")
	ErrVerifyCodeSendFrequently = errs.NewCodeError(20005, "VerifyCodeSendFrequently")
	ErrVerifyCodeNotMatch       = errs.NewCodeError(20006, "VerifyCodeNotMatch")
	ErrVerifyCodeExpired        = errs.NewCodeError(20007, "VerifyCodeExpired")
	ErrVerifyCodeMaxCount       = errs.NewCodeError(20008, "VerifyCodeMaxCount")
	ErrVerifyCodeUsed           = errs.NewCodeError(20009, "VerifyCodeUsed")
	ErrInvitationCodeUsed       = errs.NewCodeError(20010, "InvitationCodeUsed")
	ErrInvitationNotFound       = errs.NewCodeError(20011, "InvitationNotFound")
	ErrForbidden                = errs.NewCodeError(20012, "Forbidden")
	ErrRefuseFriend             = errs.NewCodeError(20013, "RefuseFriend")
	ErrEmailAlreadyRegister     = errs.NewCodeError(20014, "EmailAlreadyRegister")

	ErrTokenNotExist = errs.NewCodeError(20101, "ErrTokenNotExist")
)
