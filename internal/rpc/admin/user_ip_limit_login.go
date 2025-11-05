package admin

import (
	"context"
	"time"

	"github.com/openimsdk/tools/utils/datautil"

	admindb "github.com/openimsdk/chat/pkg/common/db/table/admin"
	"github.com/openimsdk/chat/pkg/common/mctx"
	"github.com/openimsdk/chat/pkg/protocol/admin"
	"github.com/openimsdk/tools/errs"
)

func (o *adminServer) SearchUserIPLimitLogin(ctx context.Context, req *admin.SearchUserIPLimitLoginReq) (*admin.SearchUserIPLimitLoginResp, error) {
	if _, err := mctx.CheckAdmin(ctx); err != nil {
		return nil, err
	}
	total, list, err := o.Database.SearchUserLimitLogin(ctx, req.Keyword, req.Pagination)
	if err != nil {
		return nil, err
	}
	userIDs := datautil.Slice(list, func(info *admindb.LimitUserLoginIP) string { return info.UserID })
	userMap, err := o.Chat.MapUserPublicInfo(ctx, datautil.Distinct(userIDs))
	if err != nil {
		return nil, err
	}
	limits := make([]*admin.LimitUserLoginIP, 0, len(list))
	for _, info := range list {
		limits = append(limits, &admin.LimitUserLoginIP{
			UserID:     info.UserID,
			Ip:         info.IP,
			CreateTime: info.CreateTime.UnixMilli(),
			User:       userMap[info.UserID],
		})
	}
	return &admin.SearchUserIPLimitLoginResp{Total: uint32(total), Limits: limits}, nil
}

func (o *adminServer) AddUserIPLimitLogin(ctx context.Context, req *admin.AddUserIPLimitLoginReq) (*admin.AddUserIPLimitLoginResp, error) {
	if _, err := mctx.CheckAdmin(ctx); err != nil {
		return nil, err
	}
	if len(req.Limits) == 0 {
		return nil, errs.ErrArgs.WrapMsg("limits is empty")
	}
	now := time.Now()
	ts := make([]*admindb.LimitUserLoginIP, 0, len(req.Limits))
	for _, limit := range req.Limits {
		ts = append(ts, &admindb.LimitUserLoginIP{
			UserID:     limit.UserID,
			IP:         limit.Ip,
			CreateTime: now,
		})
	}
	if err := o.Database.AddUserLimitLogin(ctx, ts); err != nil {
		return nil, err
	}
	return &admin.AddUserIPLimitLoginResp{}, nil
}

func (o *adminServer) DelUserIPLimitLogin(ctx context.Context, req *admin.DelUserIPLimitLoginReq) (*admin.DelUserIPLimitLoginResp, error) {
	if _, err := mctx.CheckAdmin(ctx); err != nil {
		return nil, err
	}
	if len(req.Limits) == 0 {
		return nil, errs.ErrArgs.WrapMsg("limits is empty")
	}
	ts := make([]*admindb.LimitUserLoginIP, 0, len(req.Limits))
	for _, limit := range req.Limits {
		if limit.UserID == "" || limit.Ip == "" {
			return nil, errs.ErrArgs.WrapMsg("user_id or ip is empty")
		}
		ts = append(ts, &admindb.LimitUserLoginIP{
			UserID: limit.UserID,
			IP:     limit.Ip,
		})
	}
	if err := o.Database.DelUserLimitLogin(ctx, ts); err != nil {
		return nil, err
	}
	return &admin.DelUserIPLimitLoginResp{}, nil
}
