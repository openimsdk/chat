package constant

import "errors"

type ErrInfo struct {
	ErrCode int32
	ErrMsg  string
}

func (e ErrInfo) Error() string {
	return e.ErrMsg
}

func (e ErrInfo) Code() int32 {
	return e.ErrCode
}

//var (
//	OK        = ErrInfo{0, ""}
//	ErrServer = ErrInfo{500, "server error"}
//
//	ErrParseToken = ErrInfo{700, ParseTokenMsg.Error()}
//
//	ErrTencentCredential = ErrInfo{400, ThirdPartyMsg.Error()}
//
//	ErrTokenKicked              = ErrInfo{706, TokenUserKickedMsg.Error()}
//	ErrTokenDifferentPlatformID = ErrInfo{707, TokenDifferentPlatformIDMsg.Error()}
//	ErrTokenDifferentUserID     = ErrInfo{708, TokenDifferentUserIDMsg.Error()}
//
//	ErrStatus                = ErrInfo{ErrCode: 804, ErrMsg: StatusMsg.Error()}
//	ErrCallback              = ErrInfo{ErrCode: 809, ErrMsg: CallBackMsg.Error()}
//	ErrSendLimit             = ErrInfo{ErrCode: 810, ErrMsg: "send msg limit, to many request, try again later"}
//	ErrMessageHasReadDisable = ErrInfo{ErrCode: 811, ErrMsg: "message has read disable"}
//
//	//通用错误
//	ErrArgs           = ErrInfo{ErrCode: FormattingError, ErrMsg: "Parameter failed"}
//	ErrDB             = ErrInfo{ErrCode: DatabaseError, ErrMsg: " operation database  failed "}
//	ErrRpcProcess     = ErrInfo{ErrCode: ServerError, ErrMsg: " rpc process failed "}
//	ErrGetRpcConn     = ErrInfo{ErrCode: ServerError, ErrMsg: " get grpc node failed "}
//	ErrRecordNotFound = ErrInfo{ErrCode: RecordNotFound, ErrMsg: " record not found "}
//	//登录相关
//	ErrNotAdmin         = ErrInfo{ErrCode: NotAdmin, ErrMsg: "not admin"}
//	ErrHasNotRegistered = ErrInfo{ErrCode: HasNotRegistered, ErrMsg: "account is not registered"}
//	ErrPasswordError    = ErrInfo{ErrCode: PasswordError, ErrMsg: "password error"}
//
//	//token相关
//	ErrTokenExpired     = ErrInfo{TokenExpired, "token expired"}
//	ErrTokenMalformed   = ErrInfo{TokenMalformed, "token format error"}
//	ErrTokenNotValidYet = ErrInfo{TokenNotValidYet, "token not valid yet"}
//	ErrTokenUnknown     = ErrInfo{TokenUnknown, "token unknown error"}
//	ErrCreateToken      = ErrInfo{CreateToken, "create token error"}
//
//	//没有权限
//	ErrNoPermission = ErrInfo{ErrCode: NoPermission, ErrMsg: "no permission"}
//)
//
//var (
//	ParseTokenMsg               = errors.New("parse token failed")
//	TokenUserKickedMsg          = errors.New("user has been kicked")
//	TokenDifferentPlatformIDMsg = errors.New("different platformID")
//	TokenDifferentUserIDMsg     = errors.New("different userID")
//	StatusMsg                   = errors.New("status is abnormal")
//	ArgsMsg                     = errors.New("args failed")
//	CallBackMsg                 = errors.New("callback failed")
//	InvitationMsg               = errors.New("invitationCode error")
//
//	ThirdPartyMsg = errors.New("third party error")
//)
//
//// 通用错误码
//const (
//	NoError         = 0     //无错误
//	FormattingError = 10001 //输入参数错误
//	DatabaseError   = 10002 //redis/mysql等db错误
//
//	ServerError     = 10003 //
//	HttpError       = 10004 //
//	GetIMTokenError = 10005 //获取OpenIM token失败
//	RecordNotFound  = 10006 //记录不存在
//
//	NoPermission = 10007
//)
//
//// 注册相关错误
//const (
//	HasRegistered     = 20001 //账号已经注册
//	RepeatSendCode    = 20002 //重复发送验证码
//	InvitationInvalid = 20003 //邀请码错误
//	RegisterLimit     = 20004 //注册受ip限制
//)
//
//// 验证码 邀请码相关
//const (
//	CodeInvalid = 30001 //验证码错误
//	CodeExpired = 30002 //验证码已过期
//	//InvitationUsed = 30003 //邀请码已被使用
//
//	InvitationCodeUsed        = 30003 // 邀请码被使用
//	InvitationCodeNonExistent = 30004 // 邀请码不存在
//)
//
//// 登录相关
//const (
//	HasNotRegistered = 40001 //账号未注册
//	PasswordError    = 40002 //密码错误
//	LoginLimit       = 40003 //登录受ip限制
//	IPForbidden      = 40004 //ip禁止 登录 注册
//	DisableLogin     = 40005 // 账号封禁
//	NotAdmin         = 41000 //非管理员
//)
//
//const (
//	TokenExpired     = 50001
//	TokenMalformed   = 50002
//	TokenNotValidYet = 50003
//	TokenUnknown     = 50004
//	CreateToken      = 50005
//)

var (
	OK        = ErrInfo{0, ""}
	ErrServer = ErrInfo{500, "server error"}

	ErrParseToken = ErrInfo{700, ParseTokenMsg.Error()}

	ErrTencentCredential = ErrInfo{400, ThirdPartyMsg.Error()}

	ErrTokenKicked              = ErrInfo{706, TokenUserKickedMsg.Error()}
	ErrTokenDifferentPlatformID = ErrInfo{707, TokenDifferentPlatformIDMsg.Error()}
	ErrTokenDifferentUserID     = ErrInfo{708, TokenDifferentUserIDMsg.Error()}

	ErrStatus                = ErrInfo{ErrCode: 804, ErrMsg: StatusMsg.Error()}
	ErrCallback              = ErrInfo{ErrCode: 809, ErrMsg: CallBackMsg.Error()}
	ErrSendLimit             = ErrInfo{ErrCode: 810, ErrMsg: "send msg limit, to many request, try again later"}
	ErrMessageHasReadDisable = ErrInfo{ErrCode: 811, ErrMsg: "message has read disable"}

	//通用错误
	ErrArgs           = ErrInfo{ErrCode: FormattingError, ErrMsg: "Parameter failed"}
	ErrDB             = ErrInfo{ErrCode: DatabaseError, ErrMsg: " operation database  failed "}
	ErrRpcProcess     = ErrInfo{ErrCode: ServerError, ErrMsg: " rpc process failed "}
	ErrGetRpcConn     = ErrInfo{ErrCode: ServerError, ErrMsg: " get grpc node failed "}
	ErrRecordNotFound = ErrInfo{ErrCode: RecordNotFound, ErrMsg: " record not found "}
	//登录相关
	ErrNotAdmin         = ErrInfo{ErrCode: NotAdmin, ErrMsg: "not admin"}
	ErrHasNotRegistered = ErrInfo{ErrCode: HasNotRegistered, ErrMsg: "account is not registered"}
	ErrPasswordError    = ErrInfo{ErrCode: PasswordError, ErrMsg: "password error"}

	//token相关
	ErrTokenExpired     = ErrInfo{TokenExpired, "token expired"}
	ErrTokenMalformed   = ErrInfo{TokenMalformed, "token format error"}
	ErrTokenNotValidYet = ErrInfo{TokenNotValidYet, "token not valid yet"}
	ErrTokenUnknown     = ErrInfo{TokenUnknown, "token unknown error"}
	ErrCreateToken      = ErrInfo{CreateToken, "create token error"}

	//没有权限
	ErrNoPermission = ErrInfo{ErrCode: NoPermission, ErrMsg: "no permission"}
)

var (
	ParseTokenMsg               = errors.New("parse token failed")
	TokenUserKickedMsg          = errors.New("user has been kicked")
	TokenDifferentPlatformIDMsg = errors.New("different platformID")
	TokenDifferentUserIDMsg     = errors.New("different userID")
	StatusMsg                   = errors.New("status is abnormal")
	ArgsMsg                     = errors.New("args failed")
	CallBackMsg                 = errors.New("callback failed")
	InvitationMsg               = errors.New("invitationCode error")
	ThirdPartyMsg               = errors.New("third party error")
)

// 通用错误码
const (
	NoError         = 0     //无错误
	FormattingError = 10001 //输入参数错误
	DatabaseError   = 10002 //redis/mysql等db错误

	ServerError     = 10003 //
	HttpError       = 10004 //
	GetIMTokenError = 10005 //获取OpenIM token失败
	RecordNotFound  = 10006 //记录不存在

	NoPermission = 10007
)

// 注册相关错误
const (
	HasRegistered     = 20001 //账号已经注册
	RepeatSendCode    = 20002 //重复发送验证码
	InvitationInvalid = 20003 //邀请码错误
	RegisterLimit     = 20004 //注册受ip限制
)

// 验证码 邀请码相关
const (
	CodeInvalid = 30001 //验证码错误
	CodeExpired = 30002 //验证码已过期
	//InvitationUsed = 30003 //邀请码已被使用
	InvitationCodeNonExistent = 30004 // 邀请码不存在
)

// 登录相关
const (
	HasNotRegistered      = 40001 //账号未注册
	PasswordError         = 40002 //密码错误
	LoginLimit            = 40003 //登录受ip限制
	IPForbidden           = 40004 //ip禁止 登录 注册
	DisableLogin          = 40005 // 账号封禁
	VerificationCodeError = 40006
	NotAdmin              = 41000 //非管理员
	NotChat               = 41001 //非用户
)

// token相关错误
const (
	TokenExpired     = 50001 //过期
	TokenMalformed   = 50002 //格式错误
	TokenNotValidYet = 50003 //不存在
	TokenUnknown     = 50004 //未知错误
	CreateToken      = 50005 //创建错误
)

const (
	GenderFemale  = 0 // 女
	GenderMale    = 1 // 男
	GenderUnknown = 2 // 未知
)

const NilTimestamp = 0 // *time.Time == nil 对应的时间戳
