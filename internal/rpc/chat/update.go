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
	"time"

	"github.com/OpenIMSDK/tools/errs"

	"github.com/OpenIMSDK/chat/pkg/proto/chat"
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
			return nil, errs.ErrArgs.Wrap("nickname can not be empty")
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
	if len(update) == 0 {
		return nil, errs.ErrArgs.Wrap("no update info")
	}
	return update, nil
}
