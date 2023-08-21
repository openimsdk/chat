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

package constant

import "github.com/OpenIMSDK/protocol/constant"

const (
	// verificationCode used for.
	VerificationCodeForRegister      = 1 // 注册
	VerificationCodeForResetPassword = 2 // 重置密码
	VerificationCodeForLogin         = 3 // 登录

	VerificationCodeForRegisterSuffix = "_forRegister"
	VerificationCodeForResetSuffix    = "_forReset"
	VerificationCodeForLoginSuffix    = "_forLogin"
)

const LogFileName = "chat.log"

// block unblock.
const (
	BlockUser   = 1
	UnblockUser = 2
)

// AccountType.
const (
	Email   = "email"
	Phone   = "phone"
	Account = "account"
)

// Mode.
const (
	UserMode  = "user"
	AdminMode = "admin"
)

// user level.
const (
	OrdinaryUserLevel = 1
	AdvancedUserLevel = 100
)

// AddFriendCtrl.
const (
	OrdinaryUserAddFriendEnable  = 1  // 允许普通用户添加好友
	OrdinaryUserAddFriendDisable = -1 // 不允许普通用户添加好友
)

// minioUpload.
const (
	OtherType = 1
	VideoType = 2
	ImageType = 3
)

// callback Action.
const (
	ActionAllow     = 0
	ActionForbidden = 1
)

const (
	ScreenInvitationRegisterAll     = 0 // 全部
	ScreenInvitationRegisterUsed    = 1 // 已使用
	ScreenInvitationRegisterNotUsed = 2 // 未使用
)

// 1 block; 2 unblock.
const (
	UserBlock   = 1 // 封号
	UserUnblock = 2 // 解封
)

const (
	NormalUser = 1
	AdminUser  = 2
)

const (
	DoNotDisturbModeDisable = 1
	DoNotDisturbModeEnable  = 2
)

const (
	AllowAddFriend    = 1
	NotAllowAddFriend = 2
)

const (
	AllowBeep    = 1
	NotAllowBeep = 2
)

const (
	AllowVibration    = 1
	NotAllowVibration = 2
)

const (
	AllowSendMsgNotFriend    = 1
	NotAllowSendMsgNotFriend = 2
)

const (
	NotNeedInvitationCodeRegister = 0 // 不需要邀请码
	NeedInvitationCodeRegister    = 1 // 需要邀请码
)

// 小程序.
const (
	StatusOnShelf = 1 // 上架
	StatusUnShelf = 2 // 下架
)

const (
	LimitIP    = 1
	NotLimitIP = 0
)

const (
	LimitNil             = 0 // 无
	LimitEmpty           = 1 // 都不限制
	LimitOnlyLoginIP     = 2 // 仅限制登录
	LimitOnlyRegisterIP  = 3 // 仅限制注册
	LimitLoginIP         = 4 // 限制登录
	LimitRegisterIP      = 5 // 限制注册
	LimitLoginRegisterIP = 6 // 限制登录注册
)

const (
	InvitationCodeAll    = 0 // 全部
	InvitationCodeUsed   = 1 // 已使用
	InvitationCodeUnused = 2 // 未使用
)

// 默认发现页面.
const DefaultDiscoverPageURL = "https://doc.rentsoft.cn/#/"

// const OperationID = "operationID"
// const OpUserID = "opUserID".
const (
	RpcOperationID = constant.OperationID
	RpcOpUserID    = constant.OpUserID
	RpcOpUserType  = "opUserType"
)

const RpcCustomHeader = constant.RpcCustomHeader

const NeedInvitationCodeRegisterConfigKey = "needInvitationCodeRegister"

const (
	DefaultAllowVibration = 1
	DefaultAllowBeep      = 1
	DefaultAllowAddFriend = 1
)

const (
	FinDAllUser    = 0
	FindNormalUser = 1
)

const DefaultPlatform = 1

const CtxApiToken = "api-token"
