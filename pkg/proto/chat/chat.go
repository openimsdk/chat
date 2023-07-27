package chat

import (
	"github.com/OpenIMSDK/chat/pkg/common/constant"
	constant2 "github.com/OpenIMSDK/protocol/constant"
	"github.com/OpenIMSDK/tools/errs"
	"regexp"
	"strconv"
)

func (x *UpdateUserInfoReq) Check() error {
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
