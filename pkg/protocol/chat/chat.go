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
	"regexp"
	"strconv"
	"strings"

	"github.com/openimsdk/chat/pkg/common/constant"
	pconstant "github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/tools/errs"
)

func (x *UpdateUserInfoReq) Check() error {
	if x.UserID == "" {
		return errs.ErrArgs.WrapMsg("userID is empty")
	}
	if x.Email != nil && x.Email.Value != "" {
		if err := EmailCheck(x.Email.Value); err != nil {
			return err
		}
	}
	return nil
}

func (x *FindUserPublicInfoReq) Check() error {
	if len(x.UserIDs) == 0 {
		return errs.ErrArgs.WrapMsg("userIDs is empty")
	}
	return nil
}

func (x *SearchUserPublicInfoReq) Check() error {
	if x.Pagination == nil {
		return errs.ErrArgs.WrapMsg("pagination is empty")
	}
	if x.Pagination.PageNumber < 1 {
		return errs.ErrArgs.WrapMsg("pageNumber is invalid")
	}
	if x.Pagination.ShowNumber < 1 {
		return errs.ErrArgs.WrapMsg("showNumber is invalid")
	}
	return nil
}

func (x *FindUserFullInfoReq) Check() error {
	if len(x.UserIDs) == 0 {
		return errs.ErrArgs.WrapMsg("userIDs is empty")
	}
	return nil
}

func (x *SendVerifyCodeReq) Check() error {
	if x.UsedFor < constant.VerificationCodeForRegister || x.UsedFor > constant.VerificationCodeForLogin {
		return errs.ErrArgs.WrapMsg("usedFor flied is empty")
	}
	if x.PhoneNumber == "" && x.Email == "" {
		return errs.ErrArgs.WrapMsg("PhoneNumber and Email is empty")
	}
	if x.PhoneNumber != "" && x.Email != "" {
		return errs.ErrArgs.WrapMsg("PhoneNumber and Email is not empty")
	}
	if x.Email != "" {
		if err := EmailCheck(x.Email); err != nil {
			return err
		}
	}
	if x.PhoneNumber != "" {
		if err := AreaCodeCheck(x.AreaCode); err != nil {
			return err
		}
		if err := PhoneNumberCheck(x.PhoneNumber); err != nil {
			return err
		}
	}
	return nil
}

func (x *VerifyCodeReq) Check() error {
	if x.PhoneNumber == "" && x.Email == "" {
		return errs.ErrArgs.WrapMsg("PhoneNumber and Email is empty")
	}
	if x.PhoneNumber != "" && x.Email != "" {
		return errs.ErrArgs.WrapMsg("PhoneNumber and Email is not empty")
	}
	if x.Email != "" {
		if err := EmailCheck(x.Email); err != nil {
			return err
		}
	}
	if x.PhoneNumber != "" {
		if err := AreaCodeCheck(x.AreaCode); err != nil {
			return err
		}
		if err := PhoneNumberCheck(x.PhoneNumber); err != nil {
			return err
		}
	}
	return nil
}

func (x *RegisterUserReq) Check() error {
	if x.User.Nickname == "" {
		return errs.ErrArgs.WrapMsg("Nickname is nil")
	}
	if x.Platform < pconstant.IOSPlatformID || x.Platform > pconstant.AdminPlatformID {
		return errs.ErrArgs.WrapMsg("platform is invalid")
	}
	if x.User == nil {
		return errs.ErrArgs.WrapMsg("user is empty")
	}
	if x.User.PhoneNumber == "" && x.User.Email == "" {
		return errs.ErrArgs.WrapMsg("PhoneNumber and Email is empty")
	}
	if x.User.PhoneNumber != "" && x.User.Email != "" {
		return errs.ErrArgs.WrapMsg("PhoneNumber and Email is not empty")
	}
	if x.User.Email != "" {
		if err := EmailCheck(x.User.Email); err != nil {
			return err
		}
	}
	if x.User.PhoneNumber != "" {
		if err := AreaCodeCheck(x.User.AreaCode); err != nil {
			return err
		}
		if err := PhoneNumberCheck(x.User.PhoneNumber); err != nil {
			return err
		}
	}
	return nil
}

func (x *LoginReq) Check() error {
	if x.Platform < pconstant.IOSPlatformID || x.Platform > pconstant.AdminPlatformID {
		return errs.ErrArgs.WrapMsg("platform is invalid")
	}
	if x.Email != "" {
		if err := EmailCheck(x.Email); err != nil {
			return err
		}
	}
	if x.PhoneNumber != "" {
		if err := AreaCodeCheck(x.AreaCode); err != nil {
			return err
		}
		if err := PhoneNumberCheck(x.PhoneNumber); err != nil {
			return err
		}
	}
	return nil
}

func (x *ResetPasswordReq) Check() error {
	if x.Password == "" {
		return errs.ErrArgs.WrapMsg("password is empty")
	}
	if x.PhoneNumber == "" && x.Email == "" {
		return errs.ErrArgs.WrapMsg("PhoneNumber and Email is empty")
	}
	if x.PhoneNumber != "" && x.Email != "" {
		return errs.ErrArgs.WrapMsg("PhoneNumber and Email is not empty")
	}
	if x.Email != "" {
		if err := EmailCheck(x.Email); err != nil {
			return err
		}
	}
	if x.PhoneNumber != "" {
		if err := AreaCodeCheck(x.AreaCode); err != nil {
			return err
		}
		if err := PhoneNumberCheck(x.PhoneNumber); err != nil {
			return err
		}
	}
	return nil
}

func (x *ChangePasswordReq) Check() error {
	if x.UserID == "" {
		return errs.ErrArgs.WrapMsg("userID is empty")
	}

	if x.NewPassword == "" {
		return errs.ErrArgs.WrapMsg("newPassword is empty")
	}

	return nil
}

func (x *FindUserAccountReq) Check() error {
	if x.UserIDs == nil {
		return errs.ErrArgs.WrapMsg("userIDs is empty")
	}
	return nil
}

func (x *FindAccountUserReq) Check() error {
	if x.Accounts == nil {
		return errs.ErrArgs.WrapMsg("Accounts is empty")
	}
	return nil
}

func (x *SearchUserFullInfoReq) Check() error {
	if x.Pagination == nil {
		return errs.ErrArgs.WrapMsg("pagination is empty")
	}
	if x.Pagination.PageNumber < 1 {
		return errs.ErrArgs.WrapMsg("pageNumber is invalid")
	}
	if x.Pagination.ShowNumber < 1 {
		return errs.ErrArgs.WrapMsg("showNumber is invalid")
	}
	if x.Normal < constant.FinDAllUser || x.Normal > constant.FindNormalUser {
		return errs.ErrArgs.WrapMsg("normal flied is invalid")
	}
	return nil
}

func (x *GetTokenForVideoMeetingReq) Check() error {
	if x.Room == "" {
		errs.ErrArgs.WrapMsg("Room is empty")
	}
	if x.Identity == "" {
		errs.ErrArgs.WrapMsg("User Identity is empty")
	}
	return nil
}

func EmailCheck(email string) error {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	if err := regexMatch(pattern, email); err != nil {
		return errs.WrapMsg(err, "Email is invalid")
	}
	return nil
}

func AreaCodeCheck(areaCode string) error {
	if areaCode == "" {
		return errs.ErrArgs.WrapMsg("AreaCode is empty")
	}
	if strings.HasPrefix(areaCode, "+") {
		areaCode = areaCode[1:]
	}
	if areaCode == "" {
		return errs.ErrArgs.WrapMsg("invalid AreaCode")
	}
	val, err := strconv.ParseUint(areaCode, 10, 32)
	if err != nil {
		return errs.ErrArgs.WrapMsg("invalid AreaCode")
	}
	if val > 1000 {
		return errs.ErrArgs.WrapMsg("invalid AreaCode")
	}
	return nil
}

func PhoneNumberCheck(phoneNumber string) error {
	if phoneNumber == "" {
		return errs.ErrArgs.WrapMsg("phoneNumber is empty")
	}
	_, err := strconv.ParseUint(phoneNumber, 10, 64)
	if err != nil {
		return errs.ErrArgs.WrapMsg("phoneNumber is invalid")
	}
	return nil
}

func regexMatch(pattern string, target string) error {
	reg := regexp.MustCompile(pattern)
	ok := reg.MatchString(target)
	if !ok {
		return errs.ErrArgs
	}
	return nil
}

func (x *SearchUserInfoReq) Check() error {
	if x.Pagination == nil {
		return errs.ErrArgs.WrapMsg("Pagination is nil")
	}
	if x.Pagination.PageNumber < 1 {
		return errs.ErrArgs.WrapMsg("pageNumber is invalid")
	}
	if x.Pagination.ShowNumber < 1 {
		return errs.ErrArgs.WrapMsg("showNumber is invalid")
	}
	return nil
}

func (x *AddUserAccountReq) Check() error {
	if x.User == nil {
		return errs.ErrArgs.WrapMsg("user is empty")
	}

	if x.User.PhoneNumber == "" && x.User.Email == "" {
		return errs.ErrArgs.WrapMsg("PhoneNumber and Email is empty")
	}
	if x.User.PhoneNumber != "" && x.User.Email != "" {
		return errs.ErrArgs.WrapMsg("PhoneNumber and Email is not empty")
	}
	if x.User.Email != "" {
		if err := EmailCheck(x.User.Email); err != nil {
			return err
		}
	}
	if x.User.PhoneNumber != "" {
		if err := AreaCodeCheck(x.User.AreaCode); err != nil {
			return err
		}
		if err := PhoneNumberCheck(x.User.PhoneNumber); err != nil {
			return err
		}
	}

	return nil
}
