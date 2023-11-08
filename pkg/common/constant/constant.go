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

// config path
const (
	ConfigPath = "/config/config.yaml"

	OpenIMConfig = "OpenIMConfig" // environment variables
)

const (
	// verificationCode used for.
	VerificationCodeForRegister      = 1 // Register
	VerificationCodeForResetPassword = 2 // Reset password
	VerificationCodeForLogin         = 3 // Login

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
	OrdinaryUserAddFriendEnable  = 1  // Allow ordinary users to add friends
	OrdinaryUserAddFriendDisable = -1 // Do not allow ordinary users to add friends
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
	ScreenInvitationRegisterAll     = 0 // All
	ScreenInvitationRegisterUsed    = 1 // Used
	ScreenInvitationRegisterNotUsed = 2 // Unused
)

// 1 block; 2 unblock.
const (
	UserBlock   = 1 // Account ban
	UserUnblock = 2 // Unban
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
	NotNeedInvitationCodeRegister = 0 // No invitation code required
	NeedInvitationCodeRegister    = 1 // Invitation code required
)

// mini-app
const (
	StatusOnShelf = 1 // OnShelf
	StatusUnShelf = 2 // UnShelf
)

const (
	LimitIP    = 1
	NotLimitIP = 0
)

const (
	LimitNil             = 0 // None
	LimitEmpty           = 1 // Neither are restricted
	LimitOnlyLoginIP     = 2 // Only login is restricted
	LimitOnlyRegisterIP  = 3 // Only registration is restricted
	LimitLoginIP         = 4 // Restrict login
	LimitRegisterIP      = 5 // Restrict registration
	LimitLoginRegisterIP = 6 // Restrict both login and registration
)

const (
	InvitationCodeAll    = 0 // All
	InvitationCodeUsed   = 1 // Used
	InvitationCodeUnused = 2 // Unused
)

// Default discovery page
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
