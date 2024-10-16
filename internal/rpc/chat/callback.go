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
	"context"
	"encoding/json"
	"fmt"

	"github.com/openimsdk/chat/pkg/common/constant"
	"github.com/openimsdk/chat/pkg/eerrs"
	"github.com/openimsdk/chat/pkg/protocol/chat"
	constantpb "github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/tools/errs"
)

type CallbackBeforeAddFriendReq struct {
	CallbackCommand `json:"callbackCommand"`
	FromUserID      string `json:"fromUserID" `
	ToUserID        string `json:"toUserID"`
	ReqMsg          string `json:"reqMsg"`
	OperationID     string `json:"operationID"`
}

type CallbackCommand string

func (c CallbackCommand) GetCallbackCommand() string {
	return string(c)
}

func (o *chatSvr) OpenIMCallback(ctx context.Context, req *chat.OpenIMCallbackReq) (*chat.OpenIMCallbackResp, error) {
	switch req.Command {
	case constantpb.CallbackBeforeAddFriendCommand:
		var data CallbackBeforeAddFriendReq
		if err := json.Unmarshal([]byte(req.Body), &data); err != nil {
			return nil, errs.Wrap(err)
		}
		user, err := o.Database.TakeAttributeByUserID(ctx, data.ToUserID)
		if err != nil {
			return nil, err
		}
		if user.AllowAddFriend != constant.OrdinaryUserAddFriendEnable {
			return nil, eerrs.ErrRefuseFriend.WrapMsg(fmt.Sprintf("state %d", user.AllowAddFriend))
		}
		return &chat.OpenIMCallbackResp{}, nil
	default:
		return nil, errs.ErrArgs.WrapMsg(fmt.Sprintf("invalid command %s", req.Command))
	}
}
