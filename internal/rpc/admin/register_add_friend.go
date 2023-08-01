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

	"github.com/OpenIMSDK/tools/log"

	"github.com/OpenIMSDK/tools/errs"
	"github.com/OpenIMSDK/tools/utils"

	admin2 "github.com/OpenIMSDK/chat/pkg/common/db/table/admin"
	"github.com/OpenIMSDK/chat/pkg/common/mctx"
	"github.com/OpenIMSDK/chat/pkg/proto/admin"
	"github.com/OpenIMSDK/chat/pkg/proto/common"
)

func (o *adminServer) AddDefaultFriend(ctx context.Context, req *admin.AddDefaultFriendReq) (*admin.AddDefaultFriendResp, error) {
	defer log.ZDebug(ctx, "return")
	if _, err := mctx.CheckAdmin(ctx); err != nil {
		return nil, err
	}
	if len(req.UserIDs) == 0 {
		return nil, errs.ErrArgs.Wrap("user ids is empty")
	}
	if utils.Duplicate(req.UserIDs) {
		return nil, errs.ErrArgs.Wrap("user ids is duplicate")
	}
	users, err := o.Chat.FindUserPublicInfo(ctx, req.UserIDs)
	if err != nil {
		return nil, err
	}
	if ids := utils.Single(req.UserIDs, utils.Slice(users, func(user *common.UserPublicInfo) string { return user.UserID })); len(ids) > 0 {
		return nil, errs.ErrUserIDNotFound.Wrap(strings.Join(ids, ", "))
	}
	exists, err := o.Database.FindDefaultFriend(ctx, req.UserIDs)
	if err != nil {
		return nil, err
	}
	if len(exists) > 0 {
		return nil, errs.ErrDuplicateKey.Wrap(strings.Join(exists, ", "))
	}
	now := time.Now()
	ms := make([]*admin2.RegisterAddFriend, 0, len(req.UserIDs))
	for _, userID := range req.UserIDs {
		ms = append(ms, &admin2.RegisterAddFriend{
			UserID:     userID,
			CreateTime: now,
		})
	}
	if err := o.Database.AddDefaultFriend(ctx, ms); err != nil {
		return nil, err
	}
	return &admin.AddDefaultFriendResp{}, nil
}

func (o *adminServer) DelDefaultFriend(ctx context.Context, req *admin.DelDefaultFriendReq) (*admin.DelDefaultFriendResp, error) {
	defer log.ZDebug(ctx, "return")
	if _, err := mctx.CheckAdmin(ctx); err != nil {
		return nil, err
	}
	if len(req.UserIDs) == 0 {
		return nil, errs.ErrArgs.Wrap("user ids is empty")
	}
	if utils.Duplicate(req.UserIDs) {
		return nil, errs.ErrArgs.Wrap("user ids is duplicate")
	}
	exists, err := o.Database.FindDefaultFriend(ctx, req.UserIDs)
	if err != nil {
		return nil, err
	}
	if ids := utils.Single(req.UserIDs, exists); len(ids) > 0 {
		return nil, errs.ErrUserIDNotFound.Wrap(strings.Join(ids, ", "))
	}
	now := time.Now()
	ms := make([]*admin2.RegisterAddFriend, 0, len(req.UserIDs))
	for _, userID := range req.UserIDs {
		ms = append(ms, &admin2.RegisterAddFriend{
			UserID:     userID,
			CreateTime: now,
		})
	}
	if err := o.Database.DelDefaultFriend(ctx, req.UserIDs); err != nil {
		return nil, err
	}
	return &admin.DelDefaultFriendResp{}, nil
}

func (o *adminServer) FindDefaultFriend(ctx context.Context, req *admin.FindDefaultFriendReq) (*admin.FindDefaultFriendResp, error) {
	defer log.ZDebug(ctx, "return")
	if _, _, err := mctx.Check(ctx); err != nil {
		return nil, err
	}
	userIDs, err := o.Database.FindDefaultFriend(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &admin.FindDefaultFriendResp{UserIDs: userIDs}, nil
}

func (o *adminServer) SearchDefaultFriend(ctx context.Context, req *admin.SearchDefaultFriendReq) (*admin.SearchDefaultFriendResp, error) {
	defer log.ZDebug(ctx, "return")
	if _, err := mctx.CheckAdmin(ctx); err != nil {
		return nil, err
	}
	total, infos, err := o.Database.SearchDefaultFriend(ctx, req.Keyword, req.Pagination.PageNumber, req.Pagination.ShowNumber)
	if err != nil {
		return nil, err
	}
	userIDs := utils.Slice(infos, func(info *admin2.RegisterAddFriend) string { return info.UserID })
	userMap, err := o.Chat.MapUserPublicInfo(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	attributes := make([]*admin.DefaultFriendAttribute, 0, len(infos))
	for _, info := range infos {
		attribute := &admin.DefaultFriendAttribute{
			UserID:     info.UserID,
			CreateTime: info.CreateTime.UnixMilli(),
			User:       userMap[info.UserID],
		}
		attributes = append(attributes, attribute)
	}
	return &admin.SearchDefaultFriendResp{Total: total, Users: attributes}, nil
}
