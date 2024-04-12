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

package admin

import (
	"crypto/md5"
	"encoding/hex"
	"time"

	"github.com/openimsdk/tools/errs"

	"github.com/openimsdk/chat/pkg/protocol/admin"
)

type Admin struct {
	Account    string    `gorm:"column:account;primary_key;type:char(64)"`
	Password   string    `gorm:"column:password;type:char(64)"`
	FaceURL    string    `gorm:"column:face_url;type:char(64)"`
	Nickname   string    `gorm:"column:nickname;type:char(64)"`
	UserID     string    `gorm:"column:user_id;type:char(64)"` // openIM userID
	Level      int32     `gorm:"column:level;default:1"  `
	CreateTime time.Time `gorm:"column:create_time"`
}

func ToDBAdminUpdate(req *admin.AdminUpdateInfoReq) (map[string]any, error) {
	update := make(map[string]any)
	if req.Account != nil {
		if req.Account.Value == "" {
			return nil, errs.ErrArgs.WrapMsg("account is empty")
		}
		update["account"] = req.Account.Value
	}
	if req.Password != nil {
		if req.Password.Value == "" {
			return nil, errs.ErrArgs.WrapMsg("password is empty")
		}
		update["password"] = req.Password.Value
	}
	if req.FaceURL != nil {
		update["face_url"] = req.FaceURL.Value
	}
	if req.Nickname != nil {
		if req.Nickname.Value == "" {
			return nil, errs.ErrArgs.WrapMsg("nickname is empty")
		}
		update["nickname"] = req.Nickname.Value
	}
	//if req.UserID != nil {
	//	update["user_id"] = req.UserID.Value
	//}
	if req.Level != nil {
		update["level"] = req.Level.Value
	}
	if len(update) == 0 {
		return nil, errs.ErrArgs.WrapMsg("no update info")
	}
	return update, nil
}

func ToDBAdminUpdatePassword(password string) (map[string]any, error) {
	if password == "" {
		return nil, errs.ErrArgs.WrapMsg("password is empty")
	}
	return map[string]any{"password": password}, nil
}

func ToDBAppletUpdate(req *admin.UpdateAppletReq) (map[string]any, error) {
	update := make(map[string]any)
	if req.Name != nil {
		if req.Name.Value == "" {
			return nil, errs.ErrArgs.WrapMsg("name is empty")
		}
		update["name"] = req.Name.Value
	}
	if req.AppID != nil {
		if req.AppID.Value == "" {
			return nil, errs.ErrArgs.WrapMsg("appID is empty")
		}
		update["app_id"] = req.AppID.Value
	}
	if req.Icon != nil {
		update["icon"] = req.Icon.Value
	}
	if req.Url != nil {
		if req.Url.Value == "" {
			return nil, errs.ErrArgs.WrapMsg("url is empty")
		}
		update["url"] = req.Url.Value
	}
	if req.Md5 != nil {
		if hash, _ := hex.DecodeString(req.Md5.Value); len(hash) != md5.Size {
			return nil, errs.ErrArgs.WrapMsg("md5 is invalid")
		}
		update["md5"] = req.Md5.Value
	}
	if req.Size != nil {
		if req.Size.Value <= 0 {
			return nil, errs.ErrArgs.WrapMsg("size is invalid")
		}
		update["size"] = req.Size.Value
	}
	if req.Version != nil {
		if req.Version.Value == "" {
			return nil, errs.ErrArgs.WrapMsg("version is empty")
		}
		update["version"] = req.Version.Value
	}
	if req.Priority != nil {
		update["priority"] = req.Priority.Value
	}
	if req.Status != nil {
		update["status"] = req.Status.Value
	}
	if len(update) == 0 {
		return nil, errs.ErrArgs.WrapMsg("no update info")
	}
	return update, nil
}

func ToDBInvitationRegisterUpdate(userID string) map[string]any {
	return map[string]any{"user_id": userID}
}
