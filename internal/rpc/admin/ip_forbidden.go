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
	"time"

	admindb "github.com/openimsdk/chat/pkg/common/db/table/admin"
	"github.com/openimsdk/chat/pkg/common/mctx"
	"github.com/openimsdk/chat/pkg/protocol/admin"
)

func (o *adminServer) SearchIPForbidden(ctx context.Context, req *admin.SearchIPForbiddenReq) (*admin.SearchIPForbiddenResp, error) {
	if _, err := mctx.CheckAdmin(ctx); err != nil {
		return nil, err
	}
	total, forbiddens, err := o.Database.SearchIPForbidden(ctx, req.Keyword, req.Status, req.Pagination)
	if err != nil {
		return nil, err
	}
	resp := &admin.SearchIPForbiddenResp{
		Forbiddens: make([]*admin.IPForbidden, 0, len(forbiddens)),
		Total:      uint32(total),
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
	tables := make([]*admindb.IPForbidden, 0, len(req.Forbiddens))
	for _, forbidden := range req.Forbiddens {
		tables = append(tables, &admindb.IPForbidden{
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
