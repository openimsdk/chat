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
