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

package admin

import (
	"github.com/openimsdk/chat/pkg/common/constant"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/utils/datautil"
)

func (x *LoginReq) Check() error {
	if x.Account == "" {
		return errs.ErrArgs.WrapMsg("account is empty")
	}
	if x.Password == "" {
		return errs.ErrArgs.WrapMsg("password is empty")
	}
	return nil
}

func (x *ChangePasswordReq) Check() error {
	if x.Password == "" {
		return errs.ErrArgs.WrapMsg("password is empty")
	}
	return nil
}

func (x *AddDefaultFriendReq) Check() error {
	if x.UserIDs == nil {
		return errs.ErrArgs.WrapMsg("userIDs is empty")
	}
	if datautil.Duplicate(x.UserIDs) {
		return errs.ErrArgs.WrapMsg("userIDs has duplicate")
	}
	return nil
}

func (x *DelDefaultFriendReq) Check() error {
	if x.UserIDs == nil {
		return errs.ErrArgs.WrapMsg("userIDs is empty")
	}
	return nil
}

func (x *SearchDefaultFriendReq) Check() error {
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

func (x *AddDefaultGroupReq) Check() error {
	if x.GroupIDs == nil {
		return errs.ErrArgs.WrapMsg("GroupIDs is empty")
	}
	if datautil.Duplicate(x.GroupIDs) {
		return errs.ErrArgs.WrapMsg("GroupIDs has duplicate")
	}
	return nil
}

func (x *DelDefaultGroupReq) Check() error {
	if x.GroupIDs == nil {
		return errs.ErrArgs.WrapMsg("GroupIDs is empty")
	}
	return nil
}

func (x *SearchDefaultGroupReq) Check() error {
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

func (x *AddInvitationCodeReq) Check() error {
	if x.Codes == nil {
		return errs.ErrArgs.WrapMsg("codes is invalid")
	}
	return nil
}

func (x *GenInvitationCodeReq) Check() error {
	if x.Len < 1 {
		return errs.ErrArgs.WrapMsg("len is invalid")
	}
	if x.Num < 1 {
		return errs.ErrArgs.WrapMsg("num is invalid")
	}
	if x.Chars == "" {
		return errs.ErrArgs.WrapMsg("chars is in invalid")
	}
	return nil
}

func (x *FindInvitationCodeReq) Check() error {
	if x.Codes == nil {
		return errs.ErrArgs.WrapMsg("codes is empty")
	}
	return nil
}

func (x *UseInvitationCodeReq) Check() error {
	if x.Code == "" {
		return errs.ErrArgs.WrapMsg("code is empty")
	}
	if x.UserID == "" {
		return errs.ErrArgs.WrapMsg("userID is empty")
	}
	return nil
}

func (x *DelInvitationCodeReq) Check() error {
	if x.Codes == nil {
		return errs.ErrArgs.WrapMsg("codes is empty")
	}
	return nil
}

func (x *SearchInvitationCodeReq) Check() error {
	if !datautil.Contain(x.Status, constant.InvitationCodeUnused, constant.InvitationCodeUsed, constant.InvitationCodeAll) {
		return errs.ErrArgs.WrapMsg("state invalid")
	}
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

func (x *SearchUserIPLimitLoginReq) Check() error {
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

func (x *AddUserIPLimitLoginReq) Check() error {
	if x.Limits == nil {
		return errs.ErrArgs.WrapMsg("limits is empty")
	}
	return nil
}

func (x *DelUserIPLimitLoginReq) Check() error {
	if x.Limits == nil {
		return errs.ErrArgs.WrapMsg("limits is empty")
	}
	return nil
}

func (x *SearchIPForbiddenReq) Check() error {
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

func (x *AddIPForbiddenReq) Check() error {
	if x.Forbiddens == nil {
		return errs.ErrArgs.WrapMsg("forbiddens is empty")
	}
	return nil
}

func (x *DelIPForbiddenReq) Check() error {
	if x.Ips == nil {
		return errs.ErrArgs.WrapMsg("ips is empty")
	}
	return nil
}

func (x *CheckRegisterForbiddenReq) Check() error {
	if x.Ip == "" {
		return errs.ErrArgs.WrapMsg("ip is empty")
	}
	return nil
}

func (x *CheckLoginForbiddenReq) Check() error {
	if x.Ip == "" && x.UserID == "" {
		return errs.ErrArgs.WrapMsg("ip and userID is empty")
	}
	return nil
}

func (x *CancellationUserReq) Check() error {
	if x.UserID == "" {
		return errs.ErrArgs.WrapMsg("userID is empty")
	}
	return nil
}

func (x *BlockUserReq) Check() error {
	if x.UserID == "" {
		return errs.ErrArgs.WrapMsg("userID is empty")
	}
	return nil
}

func (x *UnblockUserReq) Check() error {
	if x.UserIDs == nil {
		return errs.ErrArgs.WrapMsg("userIDs is empty")
	}
	return nil
}

func (x *SearchBlockUserReq) Check() error {
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

func (x *FindUserBlockInfoReq) Check() error {
	if x.UserIDs == nil {
		return errs.ErrArgs.WrapMsg("userIDs is empty")
	}
	return nil
}

func (x *CreateTokenReq) Check() error {
	if x.UserID == "" {
		return errs.ErrArgs.WrapMsg("userID is empty")
	}
	if x.UserType > constant.AdminUser || x.UserType < constant.NormalUser {
		return errs.ErrArgs.WrapMsg("userType is invalid")
	}
	return nil
}

func (x *ParseTokenReq) Check() error {
	if x.Token == "" {
		return errs.ErrArgs.WrapMsg("token is empty")
	}
	return nil
}

func (x *AddAppletReq) Check() error {
	if x.Name == "" {
		return errs.ErrArgs.WrapMsg("name is empty")
	}
	if x.AppID == "" {
		return errs.ErrArgs.WrapMsg("appID is empty")
	}
	if x.Icon == "" {
		return errs.ErrArgs.WrapMsg("icon is empty")
	}
	if x.Url == "" {
		return errs.ErrArgs.WrapMsg("url is empty")
	}
	if x.Md5 == "" {
		return errs.ErrArgs.WrapMsg("md5 is empty")
	}
	if x.Size <= 0 {
		return errs.ErrArgs.WrapMsg("size is invalid")
	}
	if x.Version == "" {
		return errs.ErrArgs.WrapMsg("version is empty")
	}
	if x.Status < constant.StatusOnShelf || x.Status > constant.StatusUnShelf {
		return errs.ErrArgs.WrapMsg("status is invalid")
	}
	return nil
}

func (x *DelAppletReq) Check() error {
	if x.AppletIds == nil {
		return errs.ErrArgs.WrapMsg("appletIds is empty")
	}
	return nil
}

func (x *UpdateAppletReq) Check() error {
	if x.Id == "" {
		return errs.ErrArgs.WrapMsg("id is empty")
	}
	return nil
}

func (x *SearchAppletReq) Check() error {
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

func (x *SetClientConfigReq) Check() error {
	if x.Config == nil {
		return errs.ErrArgs.WrapMsg("config is empty")
	}
	return nil
}

func (x *ChangeAdminPasswordReq) Check() error {
	if x.UserID == "" {
		return errs.ErrArgs.WrapMsg("userID is empty")
	}
	if x.CurrentPassword == "" {
		return errs.ErrArgs.WrapMsg("currentPassword is empty")
	}
	if x.NewPassword == "" {
		return errs.ErrArgs.WrapMsg("newPassword is empty")
	}
	if x.CurrentPassword == x.NewPassword {
		return errs.ErrArgs.WrapMsg("currentPassword is equal to newPassword")
	}
	return nil
}

func (x *AddAdminAccountReq) Check() error {
	if x.Account == "" {
		return errs.ErrArgs.WrapMsg("account is empty")
	}
	if x.Password == "" {
		return errs.ErrArgs.WrapMsg("password is empty")
	}
	return nil
}

func (x *DelAdminAccountReq) Check() error {
	if len(x.UserIDs) == 0 {
		return errs.ErrArgs.WrapMsg("userIDs is empty")
	}
	return nil
}

func (x *SearchAdminAccountReq) Check() error {
	if x.Pagination.ShowNumber == 0 {
		return errs.ErrArgs.WrapMsg("showNumber is empty")
	}
	if x.Pagination.PageNumber == 0 {
		return errs.ErrArgs.WrapMsg("pageNumber is empty")
	}
	return nil
}
