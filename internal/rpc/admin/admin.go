package admin

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/mcontext"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/chat/pkg/common/constant"
	"github.com/OpenIMSDK/chat/pkg/common/db/database"
	"github.com/OpenIMSDK/chat/pkg/common/db/dbutil"
	admin2 "github.com/OpenIMSDK/chat/pkg/common/db/table/admin"
	"github.com/OpenIMSDK/chat/pkg/common/dbconn"
	"github.com/OpenIMSDK/chat/pkg/common/mctx"
	"github.com/OpenIMSDK/chat/pkg/eerrs"
	"github.com/OpenIMSDK/chat/pkg/proto/admin"
	"github.com/OpenIMSDK/chat/pkg/rpclient/chat"
	"github.com/OpenIMSDK/chat/pkg/rpclient/openim"
	"google.golang.org/grpc"
)

func Start(zk discoveryregistry.SvcDiscoveryRegistry, server *grpc.Server) error {
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
	admin.RegisterAdminServer(server, &adminServer{
		Database: database.NewAdminDatabase(db),
		Chat:     chat.NewChatClient(zk),
		OpenIM:   openim.NewOpenIM(zk),
	})
	return nil
}

type adminServer struct {
	Database database.AdminDatabaseInterface
	Chat     *chat.ChatClient
	OpenIM   *openim.OpenIM
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

func (o *adminServer) AdminUpdateInfo(ctx context.Context, req *admin.AdminUpdateInfoReq) (*admin.AdminUpdateInfoResp, error) {
	userID, err := mctx.CheckAdmin(ctx)
	if err != nil {
		return nil, err
	}
	update, err := ToDBAdminUpdate(req)
	if err != nil {
		return nil, err
	}
	_, err = o.Database.GetAdminUserID(ctx, mcontext.GetOpUserID(ctx))
	if err != nil {
		return nil, err
	}
	if err := o.Database.UpdateAdmin(ctx, userID, update); err != nil {
		return nil, err
	}
	return &admin.AdminUpdateInfoResp{}, nil
}

func (o *adminServer) Login(ctx context.Context, req *admin.LoginReq) (*admin.LoginResp, error) {
	a, err := o.Database.GetAdmin(ctx, req.Account)
	if err != nil {
		if dbutil.IsNotFound(err) {
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
	imToken, err := o.OpenIM.UserToken(ctx, a.UserID, 1)
	if err != nil {
		return nil, err
	}
	return &admin.LoginResp{
		AdminAccount: a.Account,
		AdminToken:   adminToken.Token,
		Nickname:     a.Nickname,
		FaceURL:      a.FaceURL,
		Level:        a.Level,
		IMUserID:     a.UserID,
		IMToken:      imToken.Token,
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
