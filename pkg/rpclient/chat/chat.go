package chat

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
	"github.com/OpenIMSDK/chat/pkg/common/config"
	"github.com/OpenIMSDK/chat/pkg/proto/chat"
	"github.com/OpenIMSDK/chat/pkg/proto/common"
)

func NewChat(zk discoveryregistry.SvcDiscoveryRegistry) *Chat {
	return &Chat{
		zk: zk,
	}
}

type Chat struct {
	zk discoveryregistry.SvcDiscoveryRegistry
}

func (o *Chat) getClient(ctx context.Context) (chat.ChatClient, error) {
	conn, err := o.zk.GetConn(ctx, config.Config.RpcRegisterName.OpenImChatName)
	if err != nil {
		return nil, err
	}
	return chat.NewChatClient(conn), nil
}

func (o *Chat) FindUserPublicInfo(ctx context.Context, userIDs []string) ([]*common.UserPublicInfo, error) {
	if len(userIDs) == 0 {
		return []*common.UserPublicInfo{}, nil
	}
	client, err := o.getClient(ctx)
	if err != nil {
		return nil, err
	}
	resp, err := client.FindUserPublicInfo(ctx, &chat.FindUserPublicInfoReq{UserIDs: userIDs})
	if err != nil {
		return nil, err
	}
	return resp.Users, nil
}

func (o *Chat) MapUserPublicInfo(ctx context.Context, userIDs []string) (map[string]*common.UserPublicInfo, error) {
	users, err := o.FindUserPublicInfo(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	return utils.SliceToMap(users, func(user *common.UserPublicInfo) string {
		return user.UserID
	}), nil
}

func (o *Chat) FindUserFullInfo(ctx context.Context, userIDs []string) ([]*common.UserFullInfo, error) {
	if len(userIDs) == 0 {
		return []*common.UserFullInfo{}, nil
	}
	client, err := o.getClient(ctx)
	if err != nil {
		return nil, err
	}
	resp, err := client.FindUserFullInfo(ctx, &chat.FindUserFullInfoReq{UserIDs: userIDs})
	if err != nil {
		return nil, err
	}
	return resp.Users, nil
}

func (o *Chat) MapUserFullInfo(ctx context.Context, userIDs []string) (map[string]*common.UserFullInfo, error) {
	users, err := o.FindUserFullInfo(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	userMap := make(map[string]*common.UserFullInfo)
	for i, user := range users {
		userMap[user.UserID] = users[i]
	}
	return userMap, nil
}

func (o *Chat) UpdateUser(ctx context.Context, req *chat.UpdateUserInfoReq) error {
	client, err := o.getClient(ctx)
	if err != nil {
		return err
	}
	_, err = client.UpdateUserInfo(ctx, req)
	return err
}
