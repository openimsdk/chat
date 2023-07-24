package admin

import (
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
	"github.com/OpenIMSDK/chat/pkg/common/constant"
)

func (x *LoginReq) Check() error {
	if x.Account == "" {
		return errs.ErrArgs.Wrap("account is empty")
	}
	if x.Password == "" {
		return errs.ErrArgs.Wrap("password is empty")
	}
	return nil
}

func (x *ChangePasswordReq) Check() error {
	if x.Password == "" {
		return errs.ErrArgs.Wrap("password is empty")
	}
	return nil
}

func (x *AddDefaultFriendReq) Check() error {
	if x.UserIDs == nil {
		return errs.ErrArgs.Wrap("userIDs is empty")
	}
	if utils.Duplicate(x.UserIDs) {
		return errs.ErrArgs.Wrap("userIDs has duplicate")
	}
	return nil
}

func (x *DelDefaultFriendReq) Check() error {
	if x.UserIDs == nil {
		return errs.ErrArgs.Wrap("userIDs is empty")
	}
	return nil
}

func (x *SearchDefaultFriendReq) Check() error {
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

func (x *AddDefaultGroupReq) Check() error {
	if x.GroupIDs == nil {
		return errs.ErrArgs.Wrap("GroupIDs is empty")
	}
	if utils.Duplicate(x.GroupIDs) {
		return errs.ErrArgs.Wrap("GroupIDs has duplicate")
	}
	return nil
}

func (x *DelDefaultGroupReq) Check() error {
	if x.GroupIDs == nil {
		return errs.ErrArgs.Wrap("GroupIDs is empty")
	}
	return nil
}

func (x *SearchDefaultGroupReq) Check() error {
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

func (x *AddInvitationCodeReq) Check() error {
	if x.Codes == nil {
		return errs.ErrArgs.Wrap("codes is invalid")
	}
	return nil
}

func (x *GenInvitationCodeReq) Check() error {
	if x.Len < 1 {
		return errs.ErrArgs.Wrap("len is invalid")
	}
	if x.Num < 1 {
		return errs.ErrArgs.Wrap("num is invalid")
	}
	if x.Chars == "" {
		return errs.ErrArgs.Wrap("chars is in invalid")
	}
	return nil
}

func (x *FindInvitationCodeReq) Check() error {
	if x.Codes == nil {
		return errs.ErrArgs.Wrap("codes is empty")
	}
	return nil
}

func (x *UseInvitationCodeReq) Check() error {
	if x.Code == "" {
		return errs.ErrArgs.Wrap("code is empty")
	}
	if x.UserID == "" {
		return errs.ErrArgs.Wrap("userID is empty")
	}
	return nil
}

func (x *DelInvitationCodeReq) Check() error {
	if x.Codes == nil {
		return errs.ErrArgs.Wrap("codes is empty")
	}
	return nil
}

func (x *SearchInvitationCodeReq) Check() error {
	if x.Codes == nil {
		return errs.ErrArgs.Wrap("codes is empty")
	}
	if x.UserIDs == nil {
		return errs.ErrArgs.Wrap("userIDs is empty")
	}

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

func (x *SearchUserIPLimitLoginReq) Check() error {
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

func (x *AddUserIPLimitLoginReq) Check() error {
	if x.Limits == nil {
		return errs.ErrArgs.Wrap("limits is empty")
	}
	return nil
}

func (x *DelUserIPLimitLoginReq) Check() error {
	if x.Limits == nil {
		return errs.ErrArgs.Wrap("limits is empty")
	}
	return nil
}

func (x *SearchIPForbiddenReq) Check() error {
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

func (x *AddIPForbiddenReq) Check() error {
	if x.Forbiddens == nil {
		return errs.ErrArgs.Wrap("forbiddens is empty")
	}
	return nil
}

func (x *DelIPForbiddenReq) Check() error {
	if x.Ips == nil {
		return errs.ErrArgs.Wrap("ips is empty")
	}
	return nil
}

func (x *CheckRegisterForbiddenReq) Check() error {
	if x.Ip == "" {
		return errs.ErrArgs.Wrap("ip is empty")
	}
	return nil
}

func (x *CheckLoginForbiddenReq) Check() error {
	if x.Ip == "" && x.UserID == "" {
		return errs.ErrArgs.Wrap("ip and userID is empty")
	}
	return nil
}

func (x *CancellationUserReq) Check() error {
	if x.UserID == "" {
		return errs.ErrArgs.Wrap("userID is empty")
	}
	return nil
}

func (x *BlockUserReq) Check() error {
	if x.UserID == "" {
		return errs.ErrArgs.Wrap("userID is empty")
	}
	return nil
}

func (x *UnblockUserReq) Check() error {
	if x.UserIDs == nil {
		return errs.ErrArgs.Wrap("userIDs is empty")
	}
	return nil
}

func (x *SearchBlockUserReq) Check() error {
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

func (x *FindUserBlockInfoReq) Check() error {
	if x.UserIDs == nil {
		return errs.ErrArgs.Wrap("userIDs is empty")
	}
	return nil
}

func (x *CreateTokenReq) Check() error {
	if x.UserID == "" {
		return errs.ErrArgs.Wrap("userID is empty")
	}
	if x.UserType > constant.AdminUser || x.UserType < constant.NormalUser {
		return errs.ErrArgs.Wrap("userType is invalid")
	}
	return nil
}

func (x *ParseTokenReq) Check() error {
	if x.Token == "" {
		return errs.ErrArgs.Wrap("token is empty")
	}
	return nil
}

func (x *AddAppletReq) Check() error {
	if x.Name == "" {
		return errs.ErrArgs.Wrap("name is empty")
	}
	if x.AppID == "" {
		return errs.ErrArgs.Wrap("appID is empty")
	}
	if x.Icon == "" {
		return errs.ErrArgs.Wrap("icon is empty")
	}
	if x.Url == "" {
		return errs.ErrArgs.Wrap("url is empty")
	}
	if x.Md5 == "" {
		return errs.ErrArgs.Wrap("md5 is empty")
	}
	if x.Size <= 0 {
		return errs.ErrArgs.Wrap("size is invalid")
	}
	if x.Version == "" {
		return errs.ErrArgs.Wrap("version is empty")
	}
	if x.Status < constant.StatusOnShelf || x.Status > constant.StatusUnShelf {
		return errs.ErrArgs.Wrap("status is invalid")
	}
	return nil
}

func (x *DelAppletReq) Check() error {
	if x.AppletIds == nil {
		return errs.ErrArgs.Wrap("appletIds is empty")
	}
	return nil
}

func (x *UpdateAppletReq) Check() error {
	if x.Id == "" {
		return errs.ErrArgs.Wrap("id is empty")
	}
	return nil
}

func (x *SearchAppletReq) Check() error {
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

func (x *SetClientConfigReq) Check() error {
	if x.Config == nil {
		return errs.ErrArgs.Wrap("config is empty")
	}
	return nil
}
