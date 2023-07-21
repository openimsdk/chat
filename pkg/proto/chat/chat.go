package chat

import (
	constant2 "github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/chat/pkg/common/constant"
	"github.com/OpenIMSDK/chat/pkg/proto/common"
)

func (x *UpdateUserInfoReq) Check() error {
	if x.Email != nil && x.Email.Value != "" {
		if err := common.EmailCheck(x.Email.Value); err != nil {
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
	if x.Genders == nil {
		return errs.ErrArgs.Wrap("genders is empty")
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
	} else if err := common.AreaCodeCheck(x.AreaCode); err != nil {
		return err
	}
	if x.PhoneNumber == "" {
		return errs.ErrArgs.Wrap("PhoneNumber is empty")
	} else if err := common.PhoneNumberCheck(x.PhoneNumber); err != nil {
		return err
	}
	return nil
}

func (x *VerifyCodeReq) Check() error {
	if x.AreaCode == "" {
		return errs.ErrArgs.Wrap("AreaCode is empty")
	} else if err := common.AreaCodeCheck(x.AreaCode); err != nil {
		return err
	}
	if x.PhoneNumber == "" {
		return errs.ErrArgs.Wrap("PhoneNumber is empty")
	} else if err := common.PhoneNumberCheck(x.PhoneNumber); err != nil {
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
	if x.DeviceID == "" {
		return errs.ErrArgs.Wrap("DeviceID is empty")
	}
	if x.Platform < constant2.IOSPlatformID || x.Platform > constant2.AdminPlatformID {
		return errs.ErrArgs.Wrap("platform is invalid")
	}
	if x.User == nil {
		return errs.ErrArgs.Wrap("user is empty")
	}
	if x.User.AreaCode == "" {
		return errs.ErrArgs.Wrap("AreaCode is empty")
	} else if err := common.AreaCodeCheck(x.User.AreaCode); err != nil {
		return err
	}
	if x.User.PhoneNumber == "" {
		return errs.ErrArgs.Wrap("PhoneNumber is empty")
	} else if err := common.PhoneNumberCheck(x.User.PhoneNumber); err != nil {
		return err
	}
	if x.User.Email != "" {
		if err := common.EmailCheck(x.User.Email); err != nil {
			return err
		}
	}
	return nil
}

func (x *LoginReq) Check() error {
	if x.DeviceID == "" {
		return errs.ErrArgs.Wrap("DeviceID is empty")
	}
	if x.Platform < constant2.IOSPlatformID || x.Platform > constant2.AdminPlatformID {
		return errs.ErrArgs.Wrap("platform is invalid")
	}
	if x.PhoneNumber != "" {
		if err := common.PhoneNumberCheck(x.PhoneNumber); err != nil {
			return err
		}
	}
	if x.AreaCode != "" {
		if err := common.AreaCodeCheck(x.AreaCode); err != nil {
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
	} else if err := common.AreaCodeCheck(x.AreaCode); err != nil {
		return err
	}
	if x.PhoneNumber == "" {
		return errs.ErrArgs.Wrap("PhoneNumber is empty")
	} else if err := common.PhoneNumberCheck(x.PhoneNumber); err != nil {
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
	if x.CurrentPassword == "" {
		return errs.ErrArgs.Wrap("currentPassword is empty")
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
	if x.Genders == nil {
		return errs.ErrArgs.Wrap("genders is empty")
	}
	if x.Normal > constant.FinDAllUser || x.Normal < constant.FindNormalUser {
		return errs.ErrArgs.Wrap("normal flied is invalid")
	}
	if x.Genders == nil {
		return errs.ErrArgs.Wrap("Genders is empty")
	}
	return nil
}
