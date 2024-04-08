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
	"github.com/openimsdk/chat/pkg/common/db/table/admin"
	oldadmin "github.com/openimsdk/chat/tools/mysql2mongo/internal/mysql/table/admin"
)

type convertAdmin struct{}

func (convertAdmin) Admin(v oldadmin.Admin) admin.Admin {
	return admin.Admin{
		Account:    v.Account,
		Password:   v.Password,
		FaceURL:    v.FaceURL,
		Nickname:   v.Nickname,
		UserID:     v.UserID,
		Level:      v.Level,
		CreateTime: v.CreateTime,
	}
}

func (convertAdmin) Applet(v oldadmin.Applet) admin.Applet {
	return admin.Applet{
		ID:         v.ID,
		Name:       v.Name,
		AppID:      v.AppID,
		Icon:       v.Icon,
		URL:        v.URL,
		MD5:        v.MD5,
		Size:       v.Size,
		Version:    v.Version,
		Priority:   v.Priority,
		Status:     v.Status,
		CreateTime: v.CreateTime,
	}
}

func (convertAdmin) ClientConfig(v oldadmin.ClientConfig) admin.ClientConfig {
	return admin.ClientConfig{
		Key:   v.Key,
		Value: v.Value,
	}
}

func (convertAdmin) ForbiddenAccount(v oldadmin.ForbiddenAccount) admin.ForbiddenAccount {
	return admin.ForbiddenAccount{
		UserID:         v.UserID,
		Reason:         v.Reason,
		OperatorUserID: v.OperatorUserID,
		CreateTime:     v.CreateTime,
	}
}

func (convertAdmin) InvitationRegister(v oldadmin.InvitationRegister) admin.InvitationRegister {
	return admin.InvitationRegister{
		InvitationCode: v.InvitationCode,
		UsedByUserID:   v.UsedByUserID,
		CreateTime:     v.CreateTime,
	}
}

func (convertAdmin) IPForbidden(v oldadmin.IPForbidden) admin.IPForbidden {
	return admin.IPForbidden{
		IP:            v.IP,
		LimitRegister: v.LimitRegister,
		LimitLogin:    v.LimitLogin,
		CreateTime:    v.CreateTime,
	}
}

func (convertAdmin) LimitUserLoginIP(v oldadmin.LimitUserLoginIP) admin.LimitUserLoginIP {
	return admin.LimitUserLoginIP{
		UserID:     v.UserID,
		IP:         v.IP,
		CreateTime: v.CreateTime,
	}
}

func (convertAdmin) RegisterAddFriend(v oldadmin.RegisterAddFriend) admin.RegisterAddFriend {
	return admin.RegisterAddFriend{
		UserID:     v.UserID,
		CreateTime: v.CreateTime,
	}
}

func (convertAdmin) RegisterAddGroup(v oldadmin.RegisterAddGroup) admin.RegisterAddGroup {
	return admin.RegisterAddGroup{
		GroupID:    v.GroupID,
		CreateTime: v.CreateTime,
	}
}
