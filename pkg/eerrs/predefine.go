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
