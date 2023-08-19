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
	"context"

	"github.com/OpenIMSDK/tools/log"

	"github.com/OpenIMSDK/tools/errs"

	"github.com/OpenIMSDK/chat/pkg/common/mctx"
	"github.com/OpenIMSDK/chat/pkg/proto/admin"
)

func (o *adminServer) GetClientConfig(ctx context.Context, req *admin.GetClientConfigReq) (*admin.GetClientConfigResp, error) {
	defer log.ZDebug(ctx, "return")
	conf, err := o.Database.GetConfig(ctx)
	if err != nil {
		return nil, err
	}
	return &admin.GetClientConfigResp{Config: conf}, nil
}

func (o *adminServer) SetClientConfig(ctx context.Context, req *admin.SetClientConfigReq) (*admin.SetClientConfigResp, error) {
	defer log.ZDebug(ctx, "return")
	if _, err := mctx.CheckAdmin(ctx); err != nil {
		return nil, err
	}
	if len(req.Config) == 0 {
		return nil, errs.ErrArgs.Wrap("update config empty")
	}
	if err := o.Database.SetConfig(ctx, req.Config); err != nil {
		return nil, err
	}
	return &admin.SetClientConfigResp{}, nil
}

func (o *adminServer) DelClientConfig(ctx context.Context, req *admin.DelClientConfigReq) (*admin.DelClientConfigResp, error) {
	defer log.ZDebug(ctx, "return")
	if _, err := mctx.CheckAdmin(ctx); err != nil {
		return nil, err
	}
	if err := o.Database.DelConfig(ctx, req.Keys); err != nil {
		return nil, err
	}
	return &admin.DelClientConfigResp{}, nil
}
