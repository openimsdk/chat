package chat

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/chat/pkg/common/config"
	"github.com/OpenIMSDK/chat/pkg/common/mctx"
	"github.com/OpenIMSDK/chat/pkg/eerrs"
	"github.com/OpenIMSDK/chat/pkg/proto/admin"
)

func NewAdmin(zk discoveryregistry.SvcDiscoveryRegistry) *Admin {
	return &Admin{
		zk: zk,
	}
}

type Admin struct {
	zk discoveryregistry.SvcDiscoveryRegistry
}

func (o *Admin) client(ctx context.Context) (admin.AdminClient, error) {
	conn, err := o.zk.GetConn(ctx, config.Config.RpcRegisterName.OpenImAdminName)
	if err != nil {
		return nil, err
	}
	return admin.NewAdminClient(conn), nil
}

func (o *Admin) GetConfig(ctx context.Context) (map[string]string, error) {
	client, err := o.client(ctx)
	if err != nil {
		return nil, err
	}
	conf, err := client.GetClientConfig(ctx, &admin.GetClientConfigReq{})
	if err != nil {
		return nil, err
	}
	if conf.Config == nil {
		return map[string]string{}, nil
	}
	return conf.Config, nil
}

func (o *Admin) CheckInvitationCode(ctx context.Context, invitationCode string) error {
	client, err := o.client(ctx)
	if err != nil {
		return err
	}
	resp, err := client.FindInvitationCode(ctx, &admin.FindInvitationCodeReq{Codes: []string{invitationCode}})
	if err != nil {
		return err
	}
	if len(resp.Codes) == 0 {
		return eerrs.ErrInvitationNotFound.Wrap()
	}
	if resp.Codes[0].UsedUserID != "" {
		return eerrs.ErrInvitationCodeUsed.Wrap()
	}
	return nil
}

func (o *Admin) CheckRegister(ctx context.Context, ip string) error {
	client, err := o.client(ctx)
	if err != nil {
		return err
	}
	_, err = client.CheckRegisterForbidden(ctx, &admin.CheckRegisterForbiddenReq{Ip: ip})
	return err
}

func (o *Admin) CheckLogin(ctx context.Context, userID string, ip string) error {
	client, err := o.client(ctx)
	if err != nil {
		return err
	}
	_, err = client.CheckLoginForbidden(ctx, &admin.CheckLoginForbiddenReq{Ip: ip, UserID: userID})
	return err
}

func (o *Admin) UseInvitationCode(ctx context.Context, userID string, invitationCode string) error {
	client, err := o.client(ctx)
	if err != nil {
		return err
	}
	_, err = client.UseInvitationCode(ctx, &admin.UseInvitationCodeReq{UserID: userID, Code: invitationCode})
	return err
}

func (o *Admin) CheckNilOrAdmin(ctx context.Context) (bool, error) {
	if !mctx.HaveOpUser(ctx) {
		return false, nil
	}
	_, err := mctx.CheckAdmin(ctx)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (o *Admin) CreateToken(ctx context.Context, userID string, userType int32) (*admin.CreateTokenResp, error) {
	client, err := o.client(ctx)
	if err != nil {
		return nil, err
	}
	return client.CreateToken(ctx, &admin.CreateTokenReq{UserID: userID, UserType: userType})
}

func (o *Admin) GetDefaultFriendUserID(ctx context.Context) ([]string, error) {
	client, err := o.client(ctx)
	if err != nil {
		return nil, err
	}
	resp, err := client.FindDefaultFriend(ctx, &admin.FindDefaultFriendReq{})
	if err != nil {
		return nil, err
	}
	return resp.UserIDs, nil
}

func (o *Admin) GetDefaultGroupID(ctx context.Context) ([]string, error) {
	client, err := o.client(ctx)
	if err != nil {
		return nil, err
	}
	resp, err := client.FindDefaultGroup(ctx, &admin.FindDefaultGroupReq{})
	if err != nil {
		return nil, err
	}
	return resp.GroupIDs, nil
}
