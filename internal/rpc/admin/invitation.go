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
	"math/rand"
	"strings"
	"time"

	"github.com/openimsdk/tools/utils/datautil"

	"github.com/openimsdk/tools/errs"

	admindb "github.com/openimsdk/chat/pkg/common/db/table/admin"
	"github.com/openimsdk/chat/pkg/common/mctx"
	"github.com/openimsdk/chat/pkg/eerrs"
	"github.com/openimsdk/chat/pkg/protocol/admin"
)

func (o *adminServer) AddInvitationCode(ctx context.Context, req *admin.AddInvitationCodeReq) (*admin.AddInvitationCodeResp, error) {
	if _, err := mctx.CheckAdmin(ctx); err != nil {
		return nil, err
	}
	if len(req.Codes) == 0 {
		return nil, errs.ErrArgs.WrapMsg("codes is empty")
	}
	if datautil.Duplicate(req.Codes) {
		return nil, errs.ErrArgs.WrapMsg("codes is duplicate")
	}
	irs, err := o.Database.FindInvitationRegister(ctx, req.Codes)
	if err != nil {
		return nil, err
	}
	if len(irs) > 0 {
		ids := datautil.Slice(irs, func(info *admindb.InvitationRegister) string { return info.InvitationCode })
		return nil, errs.ErrArgs.WrapMsg("code existed", "ids", ids)
	}
	now := time.Now()
	codes := make([]*admindb.InvitationRegister, 0, len(req.Codes))
	for _, code := range req.Codes {
		codes = append(codes, &admindb.InvitationRegister{
			InvitationCode: code,
			UsedByUserID:   "",
			CreateTime:     now,
		})
	}
	if err := o.Database.CreatInvitationRegister(ctx, codes); err != nil {
		return nil, err
	}
	return &admin.AddInvitationCodeResp{}, nil
}

func (o *adminServer) GenInvitationCode(ctx context.Context, req *admin.GenInvitationCodeReq) (*admin.GenInvitationCodeResp, error) {
	if _, err := mctx.CheckAdmin(ctx); err != nil {
		return nil, err
	}
	if req.Num <= 0 || req.Len <= 0 {
		return nil, errs.ErrArgs.WrapMsg("num or len <= 0")
	}
	if len(req.Chars) == 0 {
		req.Chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	}
	now := time.Now()
	invitationRegisters := make([]*admindb.InvitationRegister, 0, req.Num)
	codes := make([]string, 0, req.Num)
	for i := int32(0); i < req.Num; i++ {
		buf := make([]byte, req.Len)
		rand.Read(buf)
		for i, b := range buf {
			buf[i] = req.Chars[b%byte(len(req.Chars))]
		}
		codes = append(codes, string(buf))
		invitationRegisters = append(invitationRegisters, &admindb.InvitationRegister{
			InvitationCode: string(buf),
			UsedByUserID:   "",
			CreateTime:     now,
		})
	}
	if datautil.Duplicate(codes) {
		return nil, errs.ErrArgs.WrapMsg("gen duplicate codes")
	}
	irs, err := o.Database.FindInvitationRegister(ctx, codes)
	if err != nil {
		return nil, err
	}
	if len(irs) > 0 {
		ids := datautil.Single(codes, datautil.Slice(irs, func(ir *admindb.InvitationRegister) string { return ir.InvitationCode }))
		return nil, errs.ErrArgs.WrapMsg(strings.Join(ids, ", "))
	}
	if err := o.Database.CreatInvitationRegister(ctx, invitationRegisters); err != nil {
		return nil, err
	}
	return &admin.GenInvitationCodeResp{}, nil
}

func (o *adminServer) FindInvitationCode(ctx context.Context, req *admin.FindInvitationCodeReq) (*admin.FindInvitationCodeResp, error) {
	if _, _, err := mctx.Check(ctx); err != nil {
		return nil, err
	}
	if len(req.Codes) == 0 {
		return nil, errs.ErrArgs.WrapMsg("codes is empty")
	}
	invitationRegisters, err := o.Database.FindInvitationRegister(ctx, req.Codes)
	if err != nil {
		return nil, err
	}
	userIDs := make([]string, 0, len(invitationRegisters))
	for _, register := range invitationRegisters {
		if register.UsedByUserID != "" {
			userIDs = append(userIDs, register.UsedByUserID)
		}
	}
	userMap, err := o.Chat.MapUserPublicInfo(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	resp := &admin.FindInvitationCodeResp{Codes: make([]*admin.InvitationRegister, 0, len(invitationRegisters))}
	for _, register := range invitationRegisters {
		resp.Codes = append(resp.Codes, &admin.InvitationRegister{
			InvitationCode: register.InvitationCode,
			CreateTime:     register.CreateTime.UnixMilli(),
			UsedUserID:     register.UsedByUserID,
			UsedUser:       userMap[register.UsedByUserID],
		})
	}
	return resp, nil
}

func (o *adminServer) UseInvitationCode(ctx context.Context, req *admin.UseInvitationCodeReq) (*admin.UseInvitationCodeResp, error) {
	if _, _, err := mctx.Check(ctx); err != nil {
		return nil, err
	}
	codes, err := o.Database.FindInvitationRegister(ctx, []string{req.Code})
	if err != nil {
		return nil, err
	}
	if len(codes) == 0 {
		return nil, eerrs.ErrInvitationNotFound.Wrap()
	}
	if codes[0].UsedByUserID != "" {
		return nil, eerrs.ErrInvitationCodeUsed.Wrap()
	}
	if err := o.Database.UpdateInvitationRegister(ctx, req.Code, ToDBInvitationRegisterUpdate(req.UserID)); err != nil {
		return nil, err
	}
	return &admin.UseInvitationCodeResp{}, nil
}

func (o *adminServer) DelInvitationCode(ctx context.Context, req *admin.DelInvitationCodeReq) (*admin.DelInvitationCodeResp, error) {
	if _, err := mctx.CheckAdmin(ctx); err != nil {
		return nil, err
	}
	if len(req.Codes) == 0 {
		return nil, errs.ErrArgs.WrapMsg("codes is empty")
	}
	if datautil.Duplicate(req.Codes) {
		return nil, errs.ErrArgs.WrapMsg("codes is duplicate")
	}
	irs, err := o.Database.FindInvitationRegister(ctx, req.Codes)
	if err != nil {
		return nil, err
	}
	if len(irs) != len(req.Codes) {
		ids := datautil.Single(req.Codes, datautil.Slice(irs, func(ir *admindb.InvitationRegister) string { return ir.InvitationCode }))
		return nil, errs.ErrArgs.WrapMsg("code not found " + strings.Join(ids, ", "))
	}
	if err := o.Database.DelInvitationRegister(ctx, req.Codes); err != nil {
		return nil, err
	}
	return &admin.DelInvitationCodeResp{}, nil
}

func (o *adminServer) SearchInvitationCode(ctx context.Context, req *admin.SearchInvitationCodeReq) (*admin.SearchInvitationCodeResp, error) {
	if _, err := mctx.CheckAdmin(ctx); err != nil {
		return nil, err
	}
	total, list, err := o.Database.SearchInvitationRegister(ctx, req.Keyword, req.Status, req.UserIDs, req.Codes, req.Pagination)
	if err != nil {
		return nil, err
	}
	userIDs := make([]string, 0, len(list))
	for _, register := range list {
		if register.UsedByUserID != "" {
			userIDs = append(userIDs, register.UsedByUserID)
		}
	}
	userMap, err := o.Chat.MapUserPublicInfo(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	invitationRegisters := make([]*admin.InvitationRegister, 0, len(list))
	for _, register := range list {
		invitationRegisters = append(invitationRegisters, &admin.InvitationRegister{
			InvitationCode: register.InvitationCode,
			CreateTime:     register.CreateTime.UnixMilli(),
			UsedUserID:     register.UsedByUserID,
			UsedUser:       userMap[register.UsedByUserID],
		})
	}
	return &admin.SearchInvitationCodeResp{
		Total: uint32(total),
		List:  invitationRegisters,
	}, nil
}
