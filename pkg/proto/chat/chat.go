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
	"github.com/OpenIMSDK/tools/utils"
	"regexp"
	"strconv"

	"github.com/OpenIMSDK/chat/pkg/common/constant"
	constant2 "github.com/OpenIMSDK/protocol/constant"
	"github.com/OpenIMSDK/tools/errs"
)

func (x *UpdateUserInfoReq) Check() error {
	if x.UserID == "" {
		return errs.ErrArgs.Wrap("userID is empty")
	}
	if x.Email != nil && x.Email.Value != "" {
		if err := EmailCheck(x.Email.Value); err != nil {
			return err
		}
	}
	return nil
}

func (x *FindUserPublicInfoReq) Check() error {
	if x.UserIDs == nil {
		return errs.ErrArgs.Wrap("userIDs is empty")
	}
	return nil
}

func (x *SearchUserPublicInfoReq) Check() error {
	if x.Pagination == nil {
		return errs.ErrArgs.Wrap("pagination is empty")
	}
	if x.Pagination.PageNumber < 1 {
		return errs.ErrArgs.Wrap("pageNumber is invalid")
	}
	if x.Pagination.ShowNumber < 1 {
		return errs.ErrArgs.Wrap("showNumber is invalid")
	}
	return nil
}

func (x *FindUserFullInfoReq) Check() error {
	if x.UserIDs == nil {
		return errs.ErrArgs.Wrap("userIDs is empty")
	}
	return nil
}

func (x *SendVerifyCodeReq) Check() error {
	if x.UsedFor < constant.VerificationCodeForRegister || x.UsedFor > constant.VerificationCodeForLogin {
		return errs.ErrArgs.Wrap("usedFor flied is empty")
	}
	if x.AreaCode == "" {
		return errs.ErrArgs.Wrap("AreaCode is empty")
	} else if err := AreaCodeCheck(x.AreaCode); err != nil {
		return err
	}
	if x.PhoneNumber == "" {
		return errs.ErrArgs.Wrap("PhoneNumber is empty")
	} else if err := PhoneNumberCheck(x.PhoneNumber); err != nil {
		return err
	}
	return nil
}

func (x *VerifyCodeReq) Check() error {
	if x.AreaCode == "" {
		return errs.ErrArgs.Wrap("AreaCode is empty")
	} else if err := AreaCodeCheck(x.AreaCode); err != nil {
		return err
	}
	if x.PhoneNumber == "" {
		return errs.ErrArgs.Wrap("PhoneNumber is empty")
	} else if err := PhoneNumberCheck(x.PhoneNumber); err != nil {
		return err
	}
	if x.VerifyCode == "" {
		return errs.ErrArgs.Wrap("VerifyCode is empty")
	}
	return nil
}

func (x *RegisterUserReq) Check() error {
	if x.VerifyCode == "" {
		return errs.ErrArgs.Wrap("VerifyCode is empty")
	}
	if x.Platform < constant2.IOSPlatformID || x.Platform > constant2.AdminPlatformID {
		return errs.ErrArgs.Wrap("platform is invalid")
	}
	if x.User == nil {
		return errs.ErrArgs.Wrap("user is empty")
	}
	if x.User.AreaCode == "" {
		return errs.ErrArgs.Wrap("AreaCode is empty")
	} else if err := AreaCodeCheck(x.User.AreaCode); err != nil {
		return err
	}
	if x.User.PhoneNumber == "" {
		return errs.ErrArgs.Wrap("PhoneNumber is empty")
	} else if err := PhoneNumberCheck(x.User.PhoneNumber); err != nil {
		return err
	}
	if x.User.Email != "" {
		if err := EmailCheck(x.User.Email); err != nil {
			return err
		}
	}
	return nil
}

func (x *LoginReq) Check() error {
	if x.Platform < constant2.IOSPlatformID || x.Platform > constant2.AdminPlatformID {
		return errs.ErrArgs.Wrap("platform is invalid")
	}
	if x.PhoneNumber != "" {
		if err := PhoneNumberCheck(x.PhoneNumber); err != nil {
			return err
		}
	}
	if x.AreaCode != "" {
		if err := AreaCodeCheck(x.AreaCode); err != nil {
			return err
		}
	}
	return nil
}

func (x *ResetPasswordReq) Check() error {
	if x.Password == "" {
		return errs.ErrArgs.Wrap("password is empty")
	}
	if x.AreaCode == "" {
		return errs.ErrArgs.Wrap("AreaCode is empty")
	} else if err := AreaCodeCheck(x.AreaCode); err != nil {
		return err
	}
	if x.PhoneNumber == "" {
		return errs.ErrArgs.Wrap("PhoneNumber is empty")
	} else if err := PhoneNumberCheck(x.PhoneNumber); err != nil {
		return err
	}
	if x.VerifyCode == "" {
		return errs.ErrArgs.Wrap("VerifyCode is empty")
	}
	return nil
}

func (x *ChangePasswordReq) Check() error {
	if x.UserID == "" {
		return errs.ErrArgs.Wrap("userID is empty")
	}
	if x.NewPassword == "" {
		return errs.ErrArgs.Wrap("newPassword is empty")
	}
	return nil
}

func (x *FindUserAccountReq) Check() error {
	if x.UserIDs == nil {
		return errs.ErrArgs.Wrap("userIDs is empty")
	}
	return nil
}

func (x *FindAccountUserReq) Check() error {
	if x.Accounts == nil {
		return errs.ErrArgs.Wrap("Accounts is empty")
	}
	return nil
}

func (x *SearchUserFullInfoReq) Check() error {
	if x.Pagination == nil {
		return errs.ErrArgs.Wrap("pagination is empty")
	}
	if x.Pagination.PageNumber < 1 {
		return errs.ErrArgs.Wrap("pageNumber is invalid")
	}
	if x.Pagination.ShowNumber < 1 {
		return errs.ErrArgs.Wrap("showNumber is invalid")
	}
	if x.Normal < constant.FinDAllUser || x.Normal > constant.FindNormalUser {
		return errs.ErrArgs.Wrap("normal flied is invalid")
	}
	return nil
}

func (x *DeleteLogsReq) Check() error {
	if x.LogIDs == nil {
		return errs.ErrArgs.Wrap("LogIDs is empty")
	}
	if utils.Duplicate(x.LogIDs) {
		return errs.ErrArgs.Wrap("Logs has duplicate")
	}
	return nil
}

func (x *UploadLogsReq) Check() error {
	if x.FileURLs == nil {
		return errs.ErrArgs.Wrap("FileUrls is empty")
	}
	if x.Platform < constant2.IOSPlatformID || x.Platform > constant2.AdminPlatformID {
		return errs.ErrArgs.Wrap("Platform is invalid")
	}
	return nil
}
func (x *SearchLogsReq) Check() error {
	if x.Pagination == nil {
		return errs.ErrArgs.Wrap("Pagination is empty")
	}
	if x.Pagination.PageNumber < 1 {
		return errs.ErrArgs.Wrap("pageNumber is invalid")
	}
	if x.Pagination.ShowNumber < 1 {
		return errs.ErrArgs.Wrap("showNumber is invalid")
	}
	return nil
}

func EmailCheck(email string) error {
	pattern := `^[0-9a-z][_.0-9a-z-]{0,31}@([0-9a-z][0-9a-z-]{0,30}[0-9a-z]\.){1,4}[a-z]{2,4}$`
	if err := regexMatch(pattern, email); err != nil {
		return errs.Wrap(err, "Email is invalid")
	}
	return nil
}

func AreaCodeCheck(areaCode string) error {
	pattern := `\+[1-9][0-9]{1,2}`
	if err := regexMatch(pattern, areaCode); err != nil {
		return errs.Wrap(err, "AreaCode is invalid")
	}
	return nil
}

func PhoneNumberCheck(phoneNumber string) error {
	if phoneNumber == "" {
		return errs.ErrArgs.Wrap("phoneNumber is empty")
	}
	_, err := strconv.ParseUint(phoneNumber, 10, 64)
	if err != nil {
		return errs.ErrArgs.Wrap("phoneNumber is invalid")
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
		return errs.ErrArgs.Wrap("Pagination is nil")
	}
	if x.Pagination.PageNumber < 1 {
		return errs.ErrArgs.Wrap("pageNumber is invalid")
	}
	if x.Pagination.ShowNumber < 1 {
		return errs.ErrArgs.Wrap("showNumber is invalid")
	}
	return nil
}
