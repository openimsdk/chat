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
	"strings"
	"time"

	"github.com/openimsdk/tools/utils/datautil"

	"github.com/google/uuid"
	"github.com/openimsdk/tools/errs"

	"github.com/openimsdk/chat/pkg/common/constant"
	admindb "github.com/openimsdk/chat/pkg/common/db/table/admin"
	"github.com/openimsdk/chat/pkg/common/mctx"
	"github.com/openimsdk/chat/pkg/protocol/admin"
	"github.com/openimsdk/chat/pkg/protocol/common"
)

func (o *adminServer) AddApplet(ctx context.Context, req *admin.AddAppletReq) (*admin.AddAppletResp, error) {
	if _, err := mctx.CheckAdmin(ctx); err != nil {
		return nil, err
	}
	if req.Name == "" {
		return nil, errs.ErrArgs.WrapMsg("name empty")
	}
	if req.AppID == "" {
		return nil, errs.ErrArgs.WrapMsg("appid empty")
	}
	if !(req.Status == constant.StatusOnShelf || req.Status == constant.StatusUnShelf) {
		return nil, errs.ErrArgs.WrapMsg("invalid status")
	}
	m := admindb.Applet{
		ID:         req.Id,
		Name:       req.Name,
		AppID:      req.AppID,
		Icon:       req.Icon,
		URL:        req.Url,
		MD5:        req.Md5,
		Size:       req.Size,
		Version:    req.Version,
		Priority:   req.Priority,
		Status:     uint8(req.Status),
		CreateTime: time.Now(),
	}
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	if err := o.Database.CreateApplet(ctx, []*admindb.Applet{&m}); err != nil {
		return nil, err
	}
	return &admin.AddAppletResp{}, nil
}

func (o *adminServer) DelApplet(ctx context.Context, req *admin.DelAppletReq) (*admin.DelAppletResp, error) {
	if _, err := mctx.CheckAdmin(ctx); err != nil {
		return nil, err
	}
	if len(req.AppletIds) == 0 {
		return nil, errs.ErrArgs.WrapMsg("AppletIds empty")
	}
	applets, err := o.Database.FindApplet(ctx, req.AppletIds)
	if err != nil {
		return nil, err
	}
	if ids := datautil.Single(req.AppletIds, datautil.Slice(applets, func(e *admindb.Applet) string { return e.ID })); len(ids) > 0 {
		return nil, errs.ErrArgs.WrapMsg("ids not found: " + strings.Join(ids, ", "))
	}
	if err := o.Database.DelApplet(ctx, req.AppletIds); err != nil {
		return nil, err
	}
	return &admin.DelAppletResp{}, nil
}

func (o *adminServer) UpdateApplet(ctx context.Context, req *admin.UpdateAppletReq) (*admin.UpdateAppletResp, error) {
	if _, err := mctx.CheckAdmin(ctx); err != nil {
		return nil, err
	}
	_, err := o.Database.GetApplet(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	update, err := ToDBAppletUpdate(req)
	if err != nil {
		return nil, err
	}
	if err := o.Database.UpdateApplet(ctx, req.Id, update); err != nil {
		return nil, err
	}
	return &admin.UpdateAppletResp{}, nil
}

func (o *adminServer) FindApplet(ctx context.Context, req *admin.FindAppletReq) (*admin.FindAppletResp, error) {
	if _, _, err := mctx.Check(ctx); err != nil {
		return nil, err
	}
	applets, err := o.Database.FindOnShelf(ctx)
	if err != nil {
		return nil, err
	}
	resp := &admin.FindAppletResp{Applets: make([]*common.AppletInfo, 0, len(applets))}
	for _, applet := range applets {
		resp.Applets = append(resp.Applets, &common.AppletInfo{
			Id:         applet.ID,
			Name:       applet.Name,
			AppID:      applet.AppID,
			Icon:       applet.Icon,
			Url:        applet.URL,
			Md5:        applet.MD5,
			Size:       applet.Size,
			Version:    applet.Version,
			Priority:   applet.Priority,
			Status:     uint32(applet.Status),
			CreateTime: applet.CreateTime.UnixMilli(),
		})
	}
	return resp, nil
}

func (o *adminServer) SearchApplet(ctx context.Context, req *admin.SearchAppletReq) (*admin.SearchAppletResp, error) {
	if _, err := mctx.CheckAdmin(ctx); err != nil {
		return nil, err
	}
	total, applets, err := o.Database.SearchApplet(ctx, req.Keyword, req.Pagination)
	if err != nil {
		return nil, err
	}
	resp := &admin.SearchAppletResp{Total: uint32(total), Applets: make([]*common.AppletInfo, 0, len(applets))}
	for _, applet := range applets {
		resp.Applets = append(resp.Applets, &common.AppletInfo{
			Id:         applet.ID,
			Name:       applet.Name,
			AppID:      applet.AppID,
			Icon:       applet.Icon,
			Url:        applet.URL,
			Md5:        applet.MD5,
			Size:       applet.Size,
			Version:    applet.Version,
			Priority:   applet.Priority,
			Status:     uint32(applet.Status),
			CreateTime: applet.CreateTime.UnixMilli(),
		})
	}
	return resp, nil
}
