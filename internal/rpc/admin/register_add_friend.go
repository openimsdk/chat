package admin

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
	admin2 "github.com/OpenIMSDK/chat/pkg/common/db/table/admin"
	"github.com/OpenIMSDK/chat/pkg/common/mctx"
	"github.com/OpenIMSDK/chat/pkg/proto/admin"
	"github.com/OpenIMSDK/chat/pkg/proto/common"
	"strings"
	"time"
)

func (o *adminServer) AddDefaultFriend(ctx context.Context, req *admin.AddDefaultFriendReq) (*admin.AddDefaultFriendResp, error) {
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
		return nil, errs.ErrUserIDExisted.Wrap(strings.Join(exists, ", "))
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
