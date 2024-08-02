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
	"fmt"
	"math/rand"
	"time"

	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/mcontext"
	"github.com/openimsdk/tools/utils/datautil"

	"github.com/openimsdk/chat/pkg/common/constant"
	"github.com/openimsdk/chat/pkg/common/db/dbutil"
	admindb "github.com/openimsdk/chat/pkg/common/db/table/admin"
	"github.com/openimsdk/chat/pkg/common/mctx"
	"github.com/openimsdk/chat/pkg/eerrs"
	"github.com/openimsdk/chat/pkg/protocol/admin"
)

func (o *adminServer) GetAdminInfo(ctx context.Context, req *admin.GetAdminInfoReq) (*admin.GetAdminInfoResp, error) {
	userID, err := mctx.CheckAdmin(ctx)
	if err != nil {
		return nil, err
	}
	a, err := o.Database.GetAdminUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &admin.GetAdminInfoResp{
		Account:    a.Account,
		Password:   a.Password,
		FaceURL:    a.FaceURL,
		Nickname:   a.Nickname,
		UserID:     a.UserID,
		Level:      a.Level,
		CreateTime: a.CreateTime.UnixMilli(),
	}, nil
}

func (o *adminServer) ChangeAdminPassword(ctx context.Context, req *admin.ChangeAdminPasswordReq) (*admin.ChangeAdminPasswordResp, error) {
	user, err := o.Database.GetAdminUserID(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	if user.Password != req.CurrentPassword {
		return nil, errs.ErrInternalServer.WrapMsg("password error")
	}

	if err := o.Database.ChangePassword(ctx, req.UserID, req.NewPassword); err != nil {
		return nil, err
	}
	return &admin.ChangeAdminPasswordResp{}, nil
}

func (o *adminServer) AddAdminAccount(ctx context.Context, req *admin.AddAdminAccountReq) (*admin.AddAdminAccountResp, error) {
	if err := o.CheckSuperAdmin(ctx); err != nil {
		return nil, err
	}

	_, err := o.Database.GetAdmin(ctx, req.Account)
	if err == nil {
		return nil, errs.ErrDuplicateKey.WrapMsg("the account is registered")
	}

	adm := &admindb.Admin{
		Account:    req.Account,
		Password:   req.Password,
		FaceURL:    req.FaceURL,
		Nickname:   req.Nickname,
		UserID:     o.genUserID(),
		Level:      80,
		CreateTime: time.Now(),
	}
	if err = o.Database.AddAdminAccount(ctx, []*admindb.Admin{adm}); err != nil {
		return nil, err
	}
	return &admin.AddAdminAccountResp{}, nil
}

func (o *adminServer) DelAdminAccount(ctx context.Context, req *admin.DelAdminAccountReq) (*admin.DelAdminAccountResp, error) {
	if err := o.CheckSuperAdmin(ctx); err != nil {
		return nil, err
	}

	if datautil.Duplicate(req.UserIDs) {
		return nil, errs.ErrArgs.WrapMsg("user ids is duplicate")
	}

	for _, userID := range req.UserIDs {
		superAdmin, err := o.Database.GetAdminUserID(ctx, userID)
		if err != nil {
			return nil, err
		}
		if superAdmin.Level == constant.AdvancedUserLevel {
			return nil, errs.ErrNoPermission.WrapMsg(fmt.Sprintf("%s is superAdminID", userID))
		}
	}

	if err := o.Database.DelAdminAccount(ctx, req.UserIDs); err != nil {
		return nil, err
	}
	return &admin.DelAdminAccountResp{}, nil
}

func (o *adminServer) SearchAdminAccount(ctx context.Context, req *admin.SearchAdminAccountReq) (*admin.SearchAdminAccountResp, error) {
	if err := o.CheckSuperAdmin(ctx); err != nil {
		return nil, err
	}

	total, adminAccounts, err := o.Database.SearchAdminAccount(ctx, req.Pagination)
	if err != nil {
		return nil, err
	}
	accounts := make([]*admin.GetAdminInfoResp, 0, len(adminAccounts))
	for _, v := range adminAccounts {
		temp := &admin.GetAdminInfoResp{
			Account:    v.Account,
			FaceURL:    v.FaceURL,
			Nickname:   v.Nickname,
			UserID:     v.UserID,
			Level:      v.Level,
			CreateTime: v.CreateTime.Unix(),
		}
		accounts = append(accounts, temp)
	}
	return &admin.SearchAdminAccountResp{Total: uint32(total), AdminAccounts: accounts}, nil
}

func (o *adminServer) AdminUpdateInfo(ctx context.Context, req *admin.AdminUpdateInfoReq) (*admin.AdminUpdateInfoResp, error) {
	userID, err := mctx.CheckAdmin(ctx)
	if err != nil {
		return nil, err
	}
	update, err := ToDBAdminUpdate(req)
	if err != nil {
		return nil, err
	}
	info, err := o.Database.GetAdminUserID(ctx, mcontext.GetOpUserID(ctx))
	if err != nil {
		return nil, err
	}
	if err := o.Database.UpdateAdmin(ctx, userID, update); err != nil {
		return nil, err
	}
	resp := &admin.AdminUpdateInfoResp{UserID: info.UserID}
	if req.Nickname == nil {
		resp.Nickname = info.Nickname
	} else {
		resp.Nickname = req.Nickname.Value
	}
	if req.FaceURL == nil {
		resp.FaceURL = info.FaceURL
	} else {
		resp.FaceURL = req.FaceURL.Value
	}
	return resp, nil
}

func (o *adminServer) Login(ctx context.Context, req *admin.LoginReq) (*admin.LoginResp, error) {
	a, err := o.Database.GetAdmin(ctx, req.Account)
	if err != nil {
		if dbutil.IsDBNotFound(err) {
			return nil, eerrs.ErrAccountNotFound.Wrap()
		}
		return nil, err
	}
	if a.Password != req.Password {
		return nil, eerrs.ErrPassword.Wrap()
	}
	adminToken, err := o.CreateToken(ctx, &admin.CreateTokenReq{UserID: a.UserID, UserType: constant.AdminUser})
	if err != nil {
		return nil, err
	}
	return &admin.LoginResp{
		AdminUserID:  a.UserID,
		AdminAccount: a.Account,
		AdminToken:   adminToken.Token,
		Nickname:     a.Nickname,
		FaceURL:      a.FaceURL,
		Level:        a.Level,
	}, nil
}

func (o *adminServer) ChangePassword(ctx context.Context, req *admin.ChangePasswordReq) (*admin.ChangePasswordResp, error) {
	userID, err := mctx.CheckAdmin(ctx)
	if err != nil {
		return nil, err
	}
	a, err := o.Database.GetAdmin(ctx, userID)
	if err != nil {
		return nil, err
	}
	update, err := ToDBAdminUpdatePassword(req.Password)
	if err != nil {
		return nil, err
	}
	if err := o.Database.UpdateAdmin(ctx, a.UserID, update); err != nil {
		return nil, err
	}
	return &admin.ChangePasswordResp{}, nil
}

func (o *adminServer) genUserID() string {
	const l = 10
	data := make([]byte, l)
	rand.Read(data)
	chars := []byte("0123456789")
	for i := 0; i < len(data); i++ {
		if i == 0 {
			data[i] = chars[1:][data[i]%9]
		} else {
			data[i] = chars[data[i]%10]
		}
	}
	return string(data)
}

func (o *adminServer) CheckSuperAdmin(ctx context.Context) error {
	userID, err := mctx.CheckAdmin(ctx)
	if err != nil {
		return err
	}

	adminUser, err := o.Database.GetAdminUserID(ctx, userID)
	if err != nil {
		return err
	}

	if adminUser.Level != constant.AdvancedUserLevel {
		return errs.ErrNoPermission.Wrap()
	}
	return nil
}
