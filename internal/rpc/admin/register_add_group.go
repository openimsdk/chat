package admin

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
	admin2 "github.com/OpenIMSDK/chat/pkg/common/db/table/admin"
	"github.com/OpenIMSDK/chat/pkg/common/mctx"
	"github.com/OpenIMSDK/chat/pkg/proto/admin"
	"strings"
	"time"
)

func (o *adminServer) AddDefaultGroup(ctx context.Context, req *admin.AddDefaultGroupReq) (*admin.AddDefaultGroupResp, error) {
	if _, err := mctx.CheckAdmin(ctx); err != nil {
		return nil, err
	}
	if len(req.GroupIDs) == 0 {
		return nil, errs.ErrArgs.Wrap("group ids is empty")
	}
	if utils.Duplicate(req.GroupIDs) {
		return nil, errs.ErrArgs.Wrap("group ids is duplicate")
	}
	groups, err := o.OpenIM.FindGroup(ctx, req.GroupIDs)
	if err != nil {
		return nil, err
	}
	if ids := utils.Single(req.GroupIDs, utils.Slice(groups, func(group *sdkws.GroupInfo) string { return group.GroupID })); len(ids) > 0 {
		return nil, errs.ErrGroupIDNotFound.Wrap(strings.Join(ids, ", "))
	}
	exists, err := o.Database.FindDefaultGroup(ctx, req.GroupIDs)
	if err != nil {
		return nil, err
	}
	if len(exists) > 0 {
		return nil, errs.ErrGroupIDExisted.Wrap(strings.Join(exists, ", "))
	}
	now := time.Now()
	ms := make([]*admin2.RegisterAddGroup, 0, len(req.GroupIDs))
	for _, groupID := range req.GroupIDs {
		ms = append(ms, &admin2.RegisterAddGroup{
			GroupID:    groupID,
			CreateTime: now,
		})
	}
	if err := o.Database.AddDefaultGroup(ctx, ms); err != nil {
		return nil, err
	}
	return &admin.AddDefaultGroupResp{}, nil
}

func (o *adminServer) DelDefaultGroup(ctx context.Context, req *admin.DelDefaultGroupReq) (*admin.DelDefaultGroupResp, error) {
	if _, err := mctx.CheckAdmin(ctx); err != nil {
		return nil, err
	}
	if len(req.GroupIDs) == 0 {
		return nil, errs.ErrArgs.Wrap("group ids is empty")
	}
	if utils.Duplicate(req.GroupIDs) {
		return nil, errs.ErrArgs.Wrap("group ids is duplicate")
	}
	exists, err := o.Database.FindDefaultGroup(ctx, req.GroupIDs)
	if err != nil {
		return nil, err
	}
	if ids := utils.Single(req.GroupIDs, exists); len(ids) > 0 {
		return nil, errs.ErrGroupIDNotFound.Wrap(strings.Join(ids, ", "))
	}
	now := time.Now()
	ms := make([]*admin2.RegisterAddGroup, 0, len(req.GroupIDs))
	for _, groupID := range req.GroupIDs {
		ms = append(ms, &admin2.RegisterAddGroup{
			GroupID:    groupID,
			CreateTime: now,
		})
	}
	if err := o.Database.DelDefaultGroup(ctx, req.GroupIDs); err != nil {
		return nil, err
	}
	return &admin.DelDefaultGroupResp{}, nil
}

func (o *adminServer) FindDefaultGroup(ctx context.Context, req *admin.FindDefaultGroupReq) (*admin.FindDefaultGroupResp, error) {
	if _, _, err := mctx.Check(ctx); err != nil {
		return nil, err
	}
	groupIDs, err := o.Database.FindDefaultGroup(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &admin.FindDefaultGroupResp{GroupIDs: groupIDs}, nil
}

func (o *adminServer) SearchDefaultGroup(ctx context.Context, req *admin.SearchDefaultGroupReq) (*admin.SearchDefaultGroupResp, error) {
	if _, err := mctx.CheckAdmin(ctx); err != nil {
		return nil, err
	}
	total, infos, err := o.Database.SearchDefaultGroup(ctx, req.Keyword, req.Pagination.PageNumber, req.Pagination.ShowNumber)
	if err != nil {
		return nil, err
	}
	groupIDs := utils.Slice(infos, func(info *admin2.RegisterAddGroup) string { return info.GroupID })
	groupMap, err := o.OpenIM.MapGroup(ctx, groupIDs)
	if err != nil {
		return nil, err
	}
	attributes := make([]*admin.GroupAttribute, 0, len(infos))
	for _, info := range infos {
		attribute := &admin.GroupAttribute{
			GroupID:    info.GroupID,
			CreateTime: info.CreateTime.UnixMilli(),
			Group:      groupMap[info.GroupID],
		}
		attributes = append(attributes, attribute)
	}
	return &admin.SearchDefaultGroupResp{Total: total, Groups: attributes}, nil
}
