// Copyright © 2023 OpenIM. All rights reserved.
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

package internal

import (
	"github.com/openimsdk/chat/pkg/common/db/table/chat"
	oldchat "github.com/openimsdk/chat/tools/mysql2mongo/internal/mysql/table/chat"
)

type convertChat struct{}

func (convertChat) Account(v oldchat.Account) chat.Account {
	return chat.Account{
		UserID:         v.UserID,
		Password:       v.Password,
		CreateTime:     v.CreateTime,
		ChangeTime:     v.ChangeTime,
		OperatorUserID: v.OperatorUserID,
	}
}

func (convertChat) Attribute(v oldchat.Attribute) chat.Attribute {
	return chat.Attribute{
		UserID:           v.UserID,
		Account:          v.Account,
		PhoneNumber:      v.PhoneNumber,
		AreaCode:         v.AreaCode,
		Email:            v.Email,
		Nickname:         v.Nickname,
		FaceURL:          v.FaceURL,
		Gender:           v.Gender,
		CreateTime:       v.CreateTime,
		ChangeTime:       v.ChangeTime,
		BirthTime:        v.BirthTime,
		Level:            v.Level,
		AllowVibration:   v.AllowVibration,
		AllowBeep:        v.AllowBeep,
		AllowAddFriend:   v.AllowAddFriend,
		GlobalRecvMsgOpt: v.GlobalRecvMsgOpt,
		RegisterType:     v.RegisterType,
	}
}

func (convertChat) Log(v oldchat.Log) chat.Log {
	return chat.Log{
		LogID:      v.LogID,
		Platform:   v.Platform,
		UserID:     v.UserID,
		CreateTime: v.CreateTime,
		Url:        v.Url,
		FileName:   v.FileName,
		SystemType: v.SystemType,
		Version:    v.Version,
		Ex:         v.Ex,
	}
}

func (convertChat) Register(v oldchat.Register) chat.Register {
	return chat.Register{
		UserID:      v.UserID,
		DeviceID:    v.DeviceID,
		IP:          v.IP,
		Platform:    v.Platform,
		AccountType: v.AccountType,
		Mode:        v.Mode,
		CreateTime:  v.CreateTime,
	}
}

func (convertChat) UserLoginRecord(v oldchat.UserLoginRecord) chat.UserLoginRecord {
	return chat.UserLoginRecord{
		UserID:    v.UserID,
		LoginTime: v.LoginTime,
		IP:        v.IP,
		DeviceID:  v.DeviceID,
		Platform:  v.Platform,
	}
}
