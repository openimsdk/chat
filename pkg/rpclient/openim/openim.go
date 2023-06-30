package openim

import (
	"context"
	"fmt"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/auth"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/friend"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/group"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/user"
)

func NewOpenIMClient(zk discoveryregistry.SvcDiscoveryRegistry) *OpenIMClient {
	return &OpenIMClient{
		zk: zk,
	}
}

type OpenIMClient struct {
	zk discoveryregistry.SvcDiscoveryRegistry
}

func (o *OpenIMClient) getUserClient(ctx context.Context) (user.UserClient, error) {
	conn, err := o.zk.GetConn(ctx, config.Config.RpcRegisterName.OpenImUserName)
	if err != nil {
		return nil, err
	}
	return user.NewUserClient(conn), nil
}

func (o *OpenIMClient) getFriendClient(ctx context.Context) (friend.FriendClient, error) {
	conn, err := o.zk.GetConn(ctx, config.Config.RpcRegisterName.OpenImFriendName)
	if err != nil {
		return nil, err
	}
	return friend.NewFriendClient(conn), nil
}

func (o *OpenIMClient) getGroupClient(ctx context.Context) (group.GroupClient, error) {
	conn, err := o.zk.GetConn(ctx, config.Config.RpcRegisterName.OpenImGroupName)
	if err != nil {
		return nil, err
	}
	return group.NewGroupClient(conn), nil
}

func (o *OpenIMClient) getAuthClient(ctx context.Context) (auth.AuthClient, error) {
	name := config.Config.RpcRegisterName.OpenImAuthName
	conn, err := o.zk.GetConn(ctx, name)
	if err != nil {
		return nil, errs.ErrInternalServer.Wrap(fmt.Sprintf("get auth <%s> client failed: %s", name, err))
	}
	return auth.NewAuthClient(conn), nil
}

func (o *OpenIMClient) UpdateUser(ctx context.Context, req *user.UpdateUserInfoReq) error {
	client, err := o.getUserClient(ctx)
	if err != nil {
		return err
	}
	_, err = client.UpdateUserInfo(ctx, req)
	return err
}

func (o *OpenIMClient) UserRegister(ctx context.Context, req *sdkws.UserInfo) error {
	client, err := o.getUserClient(ctx)
	if err != nil {
		return err
	}
	_, err = client.UserRegister(ctx, &user.UserRegisterReq{Secret: config.Config.Secret, Users: []*sdkws.UserInfo{req}})
	return err
}

func (o *OpenIMClient) AddDefaultFriend(ctx context.Context, userID string, friendUserIDs []string) error {
	client, err := o.getFriendClient(ctx)
	if err != nil {
		return err
	}
	_, err = client.ImportFriends(ctx, &friend.ImportFriendReq{
		OwnerUserID:   userID,
		FriendUserIDs: friendUserIDs,
	})
	return err
}

func (o *OpenIMClient) AddDefaultGroup(ctx context.Context, userID string, groupID string) error {
	client, err := o.getGroupClient(ctx)
	if err != nil {
		return err
	}
	_, err = client.InviteUserToGroup(ctx, &group.InviteUserToGroupReq{
		GroupID:        groupID,
		Reason:         "",
		InvitedUserIDs: []string{userID},
	})
	return err
}

func (o *OpenIMClient) UserToken(ctx context.Context, userID string, platformID int32) (*auth.UserTokenResp, error) {
	client, err := o.getAuthClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.UserToken(ctx, &auth.UserTokenReq{Secret: config.Config.Secret, PlatformID: platformID, UserID: userID})
}

func (o *OpenIMClient) FindGroup(ctx context.Context, groupIDs []string) ([]*sdkws.GroupInfo, error) {
	client, err := o.getGroupClient(ctx)
	if err != nil {
		return nil, err
	}
	resp, err := client.GetGroupsInfo(ctx, &group.GetGroupsInfoReq{GroupIDs: groupIDs})
	if err != nil {
		return nil, err
	}
	return resp.GroupInfos, nil
}

func (o *OpenIMClient) MapGroup(ctx context.Context, groupIDs []string) (map[string]*sdkws.GroupInfo, error) {
	groups, err := o.FindGroup(ctx, groupIDs)
	if err != nil {
		return nil, err
	}
	groupMap := make(map[string]*sdkws.GroupInfo)
	for i, info := range groups {
		groupMap[info.GroupID] = groups[i]
	}
	return groupMap, nil
}

func (o *OpenIMClient) ForceOffline(ctx context.Context, userID string) error {
	client, err := o.getAuthClient(ctx)
	if err != nil {
		return err
	}
	for id := range constant.PlatformID2Name {
		_, err := client.ForceLogout(ctx, &auth.ForceLogoutReq{
			PlatformID: int32(id),
			UserID:     userID,
		})
		if err != nil {
			log.ZError(ctx, "ForceOffline", err, "userID", userID, "platformID", id)
		}
	}
	return nil
}

func (o *OpenIMClient) GetGroupMemberID(ctx context.Context, groupID string) ([]string, error) {
	client, err := o.getGroupClient(ctx)
	if err != nil {
		return nil, err
	}
	resp, err := client.GetGroupMemberUserIDs(ctx, &group.GetGroupMemberUserIDsReq{GroupID: groupID})
	if err != nil {
		return nil, err
	}
	return resp.UserIDs, nil
}
