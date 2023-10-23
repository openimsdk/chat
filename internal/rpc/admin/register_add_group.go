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
	"strings"
	"time"

	"github.com/OpenIMSDK/tools/errs"
	"github.com/OpenIMSDK/tools/log"
	"github.com/OpenIMSDK/tools/utils"

	admin2 "github.com/OpenIMSDK/chat/pkg/common/db/table/admin"
	"github.com/OpenIMSDK/chat/pkg/common/mctx"
	"github.com/OpenIMSDK/chat/pkg/proto/admin"
)

func (o *adminServer) AddDefaultGroup(ctx context.Context, req *admin.AddDefaultGroupReq) (*admin.AddDefaultGroupResp, error) {
	defer log.ZDebug(ctx, "return")
	if _, err := mctx.CheckAdmin(ctx); err != nil {
		return nil, err
	}
	if len(req.GroupIDs) == 0 {
		return nil, errs.ErrArgs.Wrap("group ids is empty")
	}
	if utils.Duplicate(req.GroupIDs) {
		return nil, errs.ErrArgs.Wrap("group ids is duplicate")
	}
	exists, err := o.Database.FindDefaultGroup(ctx, req.GroupIDs)
	if err != nil {
		return nil, err
	}
	if len(exists) > 0 {
		return nil, errs.ErrGroupIDExisted.Wrap(strings.Join(exists, ", "))
	}
	now := time.Now()
	ms := make([]*admin2.RegisterAddGroup, 0, len(req.GroupIDs))
	for _, groupID := range req.GroupIDs {
		ms = append(ms, &admin2.RegisterAddGroup{
			GroupID:    groupID,
			CreateTime: now,
		})
	}
	if err := o.Database.AddDefaultGroup(ctx, ms); err != nil {
		return nil, err
	}
	return &admin.AddDefaultGroupResp{}, nil
}

func (o *adminServer) DelDefaultGroup(ctx context.Context, req *admin.DelDefaultGroupReq) (*admin.DelDefaultGroupResp, error) {
	defer log.ZDebug(ctx, "return")
	if _, err := mctx.CheckAdmin(ctx); err != nil {
		return nil, err
	}
	if len(req.GroupIDs) == 0 {
		return nil, errs.ErrArgs.Wrap("group ids is empty")
	}
	if utils.Duplicate(req.GroupIDs) {
		return nil, errs.ErrArgs.Wrap("group ids is duplicate")
	}
	exists, err := o.Database.FindDefaultGroup(ctx, req.GroupIDs)
	if err != nil {
		return nil, err
	}
	if ids := utils.Single(req.GroupIDs, exists); len(ids) > 0 {
		return nil, errs.ErrGroupIDNotFound.Wrap(strings.Join(ids, ", "))
	}
	now := time.Now()
	ms := make([]*admin2.RegisterAddGroup, 0, len(req.GroupIDs))
	for _, groupID := range req.GroupIDs {
		ms = append(ms, &admin2.RegisterAddGroup{
			GroupID:    groupID,
			CreateTime: now,
		})
	}
	if err := o.Database.DelDefaultGroup(ctx, req.GroupIDs); err != nil {
		return nil, err
	}
	return &admin.DelDefaultGroupResp{}, nil
}

func (o *adminServer) FindDefaultGroup(ctx context.Context, req *admin.FindDefaultGroupReq) (*admin.FindDefaultGroupResp, error) {
	defer log.ZDebug(ctx, "return")
	if _, _, err := mctx.Check(ctx); err != nil {
		return nil, err
	}
	groupIDs, err := o.Database.FindDefaultGroup(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &admin.FindDefaultGroupResp{GroupIDs: groupIDs}, nil
}

func (o *adminServer) SearchDefaultGroup(ctx context.Context, req *admin.SearchDefaultGroupReq) (*admin.SearchDefaultGroupResp, error) {
	defer log.ZDebug(ctx, "return")
	if _, err := mctx.CheckAdmin(ctx); err != nil {
		return nil, err
	}
	total, infos, err := o.Database.SearchDefaultGroup(ctx, req.Keyword, req.Pagination.PageNumber, req.Pagination.ShowNumber)
	if err != nil {
		return nil, err
	}
	return &admin.SearchDefaultGroupResp{Total: total, GroupIDs: utils.Slice(infos, func(info *admin2.RegisterAddGroup) string { return info.GroupID })}, nil
}
