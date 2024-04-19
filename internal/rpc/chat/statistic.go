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
	"time"

	"github.com/openimsdk/chat/pkg/protocol/chat"
	"github.com/openimsdk/tools/errs"
)

func (o *chatSvr) UserLoginCount(ctx context.Context, req *chat.UserLoginCountReq) (*chat.UserLoginCountResp, error) {
	resp := &chat.UserLoginCountResp{}
	if req.Start > req.End {
		return nil, errs.ErrArgs.WrapMsg("start > end")
	}
	total, err := o.Database.NewUserCountTotal(ctx, nil)
	if err != nil {
		return nil, err
	}
	start := time.UnixMilli(req.Start)
	end := time.UnixMilli(req.End)
	count, loginCount, err := o.Database.UserLoginCountRangeEverydayTotal(ctx, &start, &end)
	if err != nil {
		return nil, err
	}
	resp.LoginCount = loginCount
	resp.UnloginCount = total - loginCount
	resp.Count = count
	return resp, nil
}
