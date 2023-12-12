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
	"crypto/md5"
	"encoding/hex"
	"github.com/OpenIMSDK/chat/pkg/common/db/cache"
	"github.com/OpenIMSDK/tools/discoveryregistry"
	"github.com/OpenIMSDK/tools/errs"
	"github.com/OpenIMSDK/tools/log"
	"github.com/OpenIMSDK/tools/mcontext"
	"google.golang.org/grpc"
	"math/rand"
	"time"

	"github.com/OpenIMSDK/chat/pkg/common/config"
	"github.com/OpenIMSDK/chat/pkg/common/constant"
	"github.com/OpenIMSDK/chat/pkg/common/db/database"
	"github.com/OpenIMSDK/chat/pkg/common/db/dbutil"
	admin2 "github.com/OpenIMSDK/chat/pkg/common/db/table/admin"
	"github.com/OpenIMSDK/chat/pkg/common/dbconn"
	"github.com/OpenIMSDK/chat/pkg/common/mctx"
	"github.com/OpenIMSDK/chat/pkg/eerrs"
	"github.com/OpenIMSDK/chat/pkg/proto/admin"
	"github.com/OpenIMSDK/chat/pkg/rpclient/chat"
)

func Start(discov discoveryregistry.SvcDiscoveryRegistry, server *grpc.Server) error {
	db, err := dbconn.NewGormDB()
	if err != nil {
		return err
	}
	tables := []any{
		admin2.Admin{},
		admin2.Applet{},
		admin2.ForbiddenAccount{},
		admin2.InvitationRegister{},
		admin2.IPForbidden{},
		admin2.LimitUserLoginIP{},
		admin2.RegisterAddFriend{},
		admin2.RegisterAddGroup{},
		admin2.ClientConfig{},
	}
	if err := db.AutoMigrate(tables...); err != nil {
		return err
	}
	rdb, err := cache.NewRedis()
	if err != nil {
		return err
	}
	if err := database.NewAdminDatabase(db, rdb).InitAdmin(context.Background()); err != nil {
		return err
	}
	if err := discov.CreateRpcRootNodes([]string{config.Config.RpcRegisterName.OpenImAdminName, config.Config.RpcRegisterName.OpenImChatName}); err != nil {
		panic(err)
	}

	admin.RegisterAdminServer(server, &adminServer{Database: database.NewAdminDatabase(db, rdb),
		Chat: chat.NewChatClient(discov),
	})
	return nil
}

type adminServer struct {
	Database database.AdminDatabaseInterface
	Chat     *chat.ChatClient
}

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
	_, err := o.Database.GetAdmin(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	if err := o.Database.ChangePassword(ctx, req.UserID, o.passwordEncryption(req.NewPassword)); err != nil {
		return nil, err
	}
	return &admin.ChangeAdminPasswordResp{}, nil
}

func (o *adminServer) AddAdminAccount(ctx context.Context, req *admin.AddAdminAccountReq) (*admin.AddAdminAccountResp, error) {
	userID, err := mctx.CheckAdmin(ctx)
	if err != nil {
		return nil, err
	}

	adminUser, err := o.Database.GetAdminUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if adminUser.Level != constant.AdvancedUserLevel {
		return nil, errs.ErrNoPermission.Wrap()
	}

	_, err = o.Database.GetAdmin(ctx, req.Account)
	if err == nil {
		return nil, errs.ErrRegisteredAlready.Wrap("the account is registered")
	}

	adm := &admin2.Admin{
		Account:    req.Account,
		Password:   o.passwordEncryption(req.Password),
		FaceURL:    req.FaceURL,
		Nickname:   req.Nickname,
		UserID:     o.genUserID(),
		Level:      80,
		CreateTime: time.Now(),
	}
	if err = o.Database.AddAdminAccount(ctx, adm); err != nil {
		return nil, err
	}
	return &admin.AddAdminAccountResp{}, nil
}

func (o *adminServer) DelAdminAccount(ctx context.Context, req *admin.DelAdminAccountReq) (*admin.DelAdminAccountResp, error) {
	if err := o.Database.DelAdminAccount(ctx, req.UserIDs); err != nil {
		return nil, err
	}
	return &admin.DelAdminAccountResp{}, nil
}

func (o *adminServer) SearchAdminAccount(ctx context.Context, req *admin.SearchAdminAccountReq) (*admin.SearchAdminAccountResp, error) {
	defer log.ZDebug(ctx, "return")
	if _, err := mctx.CheckAdmin(ctx); err != nil {
		return nil, err
	}
	var resp *admin.SearchAdminAccountResp
	total, adminAccounts, err := o.Database.SearchAdminAccount(ctx, req.Keyword, req.Pagination.ShowNumber, req.Pagination.PageNumber)
	if err != nil {
		return nil, err
	}
	resp.Total = total
	accounts := make([]*admin.GetAdminInfoResp, len(adminAccounts))
	for _, v := range adminAccounts {
		temp := &admin.GetAdminInfoResp{
			Account:    v.Account,
			Password:   v.Password,
			FaceURL:    v.FaceURL,
			Nickname:   v.Nickname,
			UserID:     v.UserID,
			Level:      v.Level,
			CreateTime: v.CreateTime.Unix(),
		}
		accounts = append(accounts, temp)
	}
	resp.AdminAccounts = accounts
	return resp, nil
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
		resp.Nickname = info.FaceURL
	} else {
		resp.FaceURL = req.FaceURL.Value
	}
	return resp, nil
}

func (o *adminServer) Login(ctx context.Context, req *admin.LoginReq) (*admin.LoginResp, error) {
	a, err := o.Database.GetAdmin(ctx, req.Account)
	if err != nil {
		if dbutil.IsGormNotFound(err) {
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

func (o *adminServer) passwordEncryption(password string) string {
	paswd := md5.Sum([]byte(password))
	return hex.EncodeToString(paswd[:])
}
