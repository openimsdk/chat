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
	"github.com/openimsdk/protocol/wrapperspb"
	"github.com/openimsdk/tools/utils/datautil"
	"strings"
	"time"

	"github.com/openimsdk/chat/pkg/common/db/dbutil"
	admindb "github.com/openimsdk/chat/pkg/common/db/table/admin"
	"github.com/openimsdk/chat/pkg/common/mctx"
	"github.com/openimsdk/chat/pkg/protocol/admin"
	"github.com/openimsdk/chat/pkg/protocol/chat"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/mcontext"
)

func (o *adminServer) CancellationUser(ctx context.Context, req *admin.CancellationUserReq) (*admin.CancellationUserResp, error) {
	if _, err := mctx.CheckAdmin(ctx); err != nil {
		return nil, err
	}
	empty := wrapperspb.String("")
	update := &chat.UpdateUserInfoReq{UserID: req.UserID, Account: empty, AreaCode: empty, PhoneNumber: empty, Email: empty}
	if err := o.Chat.UpdateUser(ctx, update); err != nil {
		return nil, err
	}
	return &admin.CancellationUserResp{}, nil
}

func (o *adminServer) BlockUser(ctx context.Context, req *admin.BlockUserReq) (*admin.BlockUserResp, error) {
	_, err := mctx.CheckAdmin(ctx)
	if err != nil {
		return nil, err
	}
	_, err = o.Database.GetBlockInfo(ctx, req.UserID)
	if err == nil {
		return nil, errs.ErrArgs.WrapMsg("user already blocked")
	} else if !dbutil.IsDBNotFound(err) {
		return nil, err
	}

	t := &admindb.ForbiddenAccount{
		UserID:         req.UserID,
		Reason:         req.Reason,
		OperatorUserID: mcontext.GetOpUserID(ctx),
		CreateTime:     time.Now(),
	}
	if err := o.Database.BlockUser(ctx, []*admindb.ForbiddenAccount{t}); err != nil {
		return nil, err
	}
	return &admin.BlockUserResp{}, nil
}

func (o *adminServer) UnblockUser(ctx context.Context, req *admin.UnblockUserReq) (*admin.UnblockUserResp, error) {
	if _, err := mctx.CheckAdmin(ctx); err != nil {
		return nil, err
	}
	if len(req.UserIDs) == 0 {
		return nil, errs.ErrArgs.WrapMsg("empty user id")
	}
	if datautil.Duplicate(req.UserIDs) {
		return nil, errs.ErrArgs.WrapMsg("duplicate user id")
	}
	bs, err := o.Database.FindBlockInfo(ctx, req.UserIDs)
	if err != nil {
		return nil, err
	}
	if len(req.UserIDs) != len(bs) {
		ids := datautil.Single(req.UserIDs, datautil.Slice(bs, func(info *admindb.ForbiddenAccount) string { return info.UserID }))
		return nil, errs.ErrArgs.WrapMsg("user not blocked " + strings.Join(ids, ", "))
	}
	if err := o.Database.DelBlockUser(ctx, req.UserIDs); err != nil {
		return nil, err
	}
	return &admin.UnblockUserResp{}, nil
}

func (o *adminServer) SearchBlockUser(ctx context.Context, req *admin.SearchBlockUserReq) (*admin.SearchBlockUserResp, error) {
	if _, err := mctx.CheckAdmin(ctx); err != nil {
		return nil, err
	}
	total, infos, err := o.Database.SearchBlockUser(ctx, req.Keyword, req.Pagination)
	if err != nil {
		return nil, err
	}
	userIDs := datautil.Slice(infos, func(info *admindb.ForbiddenAccount) string { return info.UserID })
	userMap, err := o.Chat.MapUserFullInfo(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	users := make([]*admin.BlockUserInfo, 0, len(infos))
	for _, info := range infos {
		user := &admin.BlockUserInfo{
			UserID:     info.UserID,
			Reason:     info.Reason,
			OpUserID:   info.OperatorUserID,
			CreateTime: info.CreateTime.UnixMilli(),
		}
		if userFull := userMap[info.UserID]; userFull != nil {
			user.Account = userFull.Account
			user.PhoneNumber = userFull.PhoneNumber
			user.AreaCode = userFull.AreaCode
			user.Email = userFull.Email
			user.Nickname = userFull.Nickname
			user.FaceURL = userFull.FaceURL
			user.Gender = userFull.Gender
		}
		users = append(users, user)
	}
	return &admin.SearchBlockUserResp{Total: uint32(total), Users: users}, nil
}

func (o *adminServer) FindUserBlockInfo(ctx context.Context, req *admin.FindUserBlockInfoReq) (*admin.FindUserBlockInfoResp, error) {
	if _, err := mctx.CheckAdmin(ctx); err != nil {
		return nil, err
	}
	list, err := o.Database.FindBlockUser(ctx, req.UserIDs)
	if err != nil {
		return nil, err
	}
	blocks := make([]*admin.BlockInfo, 0, len(list))
	for _, info := range list {
		blocks = append(blocks, &admin.BlockInfo{
			UserID:     info.UserID,
			Reason:     info.Reason,
			OpUserID:   info.OperatorUserID,
			CreateTime: info.CreateTime.UnixMilli(),
		})
	}
	return &admin.FindUserBlockInfoResp{Blocks: blocks}, nil
}
