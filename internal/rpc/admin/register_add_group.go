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
	"github.com/openimsdk/tools/utils/datautil"
	"time"

	"github.com/openimsdk/tools/errs"

	admindb "github.com/openimsdk/chat/pkg/common/db/table/admin"
	"github.com/openimsdk/chat/pkg/common/mctx"
	"github.com/openimsdk/chat/pkg/protocol/admin"
)

func (o *adminServer) AddDefaultGroup(ctx context.Context, req *admin.AddDefaultGroupReq) (*admin.AddDefaultGroupResp, error) {
	if _, err := mctx.CheckAdmin(ctx); err != nil {
		return nil, err
	}
	if len(req.GroupIDs) == 0 {
		return nil, errs.ErrArgs.WrapMsg("group ids is empty")
	}
	if datautil.Duplicate(req.GroupIDs) {
		return nil, errs.ErrArgs.WrapMsg("group ids is duplicate")
	}
	exists, err := o.Database.FindDefaultGroup(ctx, req.GroupIDs)
	if err != nil {
		return nil, err
	}
	if len(exists) > 0 {
		return nil, errs.ErrDuplicateKey.WrapMsg("group id existed", "groupID", exists)
	}
	now := time.Now()
	ms := make([]*admindb.RegisterAddGroup, 0, len(req.GroupIDs))
	for _, groupID := range req.GroupIDs {
		ms = append(ms, &admindb.RegisterAddGroup{
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
	if _, err := mctx.CheckAdmin(ctx); err != nil {
		return nil, err
	}
	if len(req.GroupIDs) == 0 {
		return nil, errs.ErrArgs.WrapMsg("group ids is empty")
	}
	if datautil.Duplicate(req.GroupIDs) {
		return nil, errs.ErrArgs.WrapMsg("group ids is duplicate")
	}
	exists, err := o.Database.FindDefaultGroup(ctx, req.GroupIDs)
	if err != nil {
		return nil, err
	}
	if ids := datautil.Single(req.GroupIDs, exists); len(ids) > 0 {
		return nil, errs.ErrRecordNotFound.WrapMsg("group id not found", "groupID", ids)
	}
	now := time.Now()
	ms := make([]*admindb.RegisterAddGroup, 0, len(req.GroupIDs))
	for _, groupID := range req.GroupIDs {
		ms = append(ms, &admindb.RegisterAddGroup{
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
	if _, err := mctx.CheckAdmin(ctx); err != nil {
		return nil, err
	}
	total, infos, err := o.Database.SearchDefaultGroup(ctx, req.Keyword, req.Pagination)
	if err != nil {
		return nil, err
	}
	return &admin.SearchDefaultGroupResp{Total: uint32(total), GroupIDs: datautil.Slice(infos, func(info *admindb.RegisterAddGroup) string { return info.GroupID })}, nil
}
