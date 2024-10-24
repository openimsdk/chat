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

package chat

import (
	"github.com/openimsdk/chat/pkg/common/constant"
	chatdb "github.com/openimsdk/chat/pkg/common/db/table/chat"
	"time"

	"github.com/openimsdk/tools/errs"

	"github.com/openimsdk/chat/pkg/protocol/chat"
)

func ToDBAttributeUpdate(req *chat.UpdateUserInfoReq) (map[string]any, error) {
	update := make(map[string]any)
	if req.Account != nil {
		update["account"] = req.Account.Value
	}
	if req.AreaCode != nil {
		update["area_code"] = req.AreaCode.Value
	}
	if req.Email != nil {
		update["email"] = req.Email.Value
	}
	if req.Nickname != nil {
		if req.Nickname.Value == "" {
			return nil, errs.ErrArgs.WrapMsg("nickname can not be empty")
		}
		update["nickname"] = req.Nickname.Value
	}
	if req.FaceURL != nil {
		update["face_url"] = req.FaceURL.Value
	}
	if req.Gender != nil {
		update["gender"] = req.Gender.Value
	}
	if req.Level != nil {
		update["level"] = req.Level.Value
	}
	if req.Birth != nil {
		update["birth_time"] = time.UnixMilli(req.Birth.Value)
	}
	if req.AllowAddFriend != nil {
		update["allow_add_friend"] = req.AllowAddFriend.Value
	}
	if req.AllowBeep != nil {
		update["allow_beep"] = req.AllowBeep.Value
	}
	if req.AllowVibration != nil {
		update["allow_vibration"] = req.AllowVibration.Value
	}
	if req.GlobalRecvMsgOpt != nil {
		update["global_recv_msg_opt"] = req.GlobalRecvMsgOpt.Value
	}
	//if len(update) == 0 {
	//	return nil, errs.ErrArgs.WrapMsg("no update info")
	//}
	return update, nil
}

func ToDBCredentialUpdate(req *chat.UpdateUserInfoReq, allowChange bool) ([]*chatdb.Credential, []*chatdb.Credential, error) {
	update := make([]*chatdb.Credential, 0)
	del := make([]*chatdb.Credential, 0)
	if req.Account != nil {
		if req.Account.GetValue() == "" {
			del = append(del, &chatdb.Credential{
				UserID: req.UserID,
				Type:   constant.CredentialAccount,
			})
		} else {
			update = append(update, &chatdb.Credential{
				UserID:      req.UserID,
				Account:     req.Account.GetValue(),
				Type:        constant.CredentialAccount,
				AllowChange: allowChange,
			})
		}
	}

	if req.Email != nil {
		if req.Email.GetValue() == "" {
			del = append(del, &chatdb.Credential{
				UserID: req.UserID,
				Type:   constant.CredentialEmail,
			})
		} else {
			update = append(update, &chatdb.Credential{
				UserID:      req.UserID,
				Account:     req.Email.GetValue(),
				Type:        constant.CredentialEmail,
				AllowChange: allowChange,
			})
		}
	}
	if req.PhoneNumber != nil {
		if req.PhoneNumber.GetValue() == "" {
			del = append(del, &chatdb.Credential{
				UserID: req.UserID,
				Type:   constant.CredentialPhone,
			})
		} else {
			update = append(update, &chatdb.Credential{
				UserID:      req.UserID,
				Account:     BuildCredentialPhone(req.AreaCode.GetValue(), req.PhoneNumber.GetValue()),
				Type:        constant.CredentialPhone,
				AllowChange: allowChange,
			})
		}
	}

	return update, del, nil
}
