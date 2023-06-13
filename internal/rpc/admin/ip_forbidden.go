package admin

import (
	"context"
	admin2 "github.com/OpenIMSDK/chat/pkg/common/db/table/admin"
	"github.com/OpenIMSDK/chat/pkg/common/mctx"
	"github.com/OpenIMSDK/chat/pkg/proto/admin"
	"time"
)

func (o *adminServer) SearchIPForbidden(ctx context.Context, req *admin.SearchIPForbiddenReq) (*admin.SearchIPForbiddenResp, error) {
	if _, err := mctx.CheckAdmin(ctx); err != nil {
		return nil, err
	}
	total, forbiddens, err := o.Database.SearchIPForbidden(ctx, req.Keyword, req.Status, req.Pagination.PageNumber, req.Pagination.ShowNumber)
	if err != nil {
		return nil, err
	}
	resp := &admin.SearchIPForbiddenResp{
		Forbiddens: make([]*admin.IPForbidden, 0, len(forbiddens)),
		Total:      total,
	}
	for _, forbidden := range forbiddens {
		resp.Forbiddens = append(resp.Forbiddens, &admin.IPForbidden{
			Ip:            forbidden.IP,
			LimitLogin:    forbidden.LimitLogin,
			LimitRegister: forbidden.LimitRegister,
			CreateTime:    forbidden.CreateTime.UnixMilli(),
		})
	}
	return resp, nil
}

func (o *adminServer) AddIPForbidden(ctx context.Context, req *admin.AddIPForbiddenReq) (*admin.AddIPForbiddenResp, error) {
	if _, err := mctx.CheckAdmin(ctx); err != nil {
		return nil, err
	}
	now := time.Now()
	tables := make([]*admin2.IPForbidden, 0, len(req.Forbiddens))
	for _, forbidden := range req.Forbiddens {
		tables = append(tables, &admin2.IPForbidden{
			IP:            forbidden.Ip,
			LimitLogin:    forbidden.LimitLogin,
			LimitRegister: forbidden.LimitRegister,
			CreateTime:    now,
		})
	}
	if err := o.Database.AddIPForbidden(ctx, tables); err != nil {
		return nil, err
	}
	return &admin.AddIPForbiddenResp{}, nil
}

func (o *adminServer) DelIPForbidden(ctx context.Context, req *admin.DelIPForbiddenReq) (*admin.DelIPForbiddenResp, error) {
	if _, err := mctx.CheckAdmin(ctx); err != nil {
		return nil, err
	}
	if err := o.Database.DelIPForbidden(ctx, req.Ips); err != nil {
		return nil, err
	}
	return &admin.DelIPForbiddenResp{}, nil
}
