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

	"github.com/OpenIMSDK/tools/log"

	"github.com/OpenIMSDK/tools/errs"
	"github.com/OpenIMSDK/tools/utils"
	"github.com/google/uuid"

	"github.com/OpenIMSDK/chat/pkg/common/constant"
	admin2 "github.com/OpenIMSDK/chat/pkg/common/db/table/admin"
	"github.com/OpenIMSDK/chat/pkg/common/mctx"
	"github.com/OpenIMSDK/chat/pkg/proto/admin"
	"github.com/OpenIMSDK/chat/pkg/proto/common"
)

func (o *adminServer) AddApplet(ctx context.Context, req *admin.AddAppletReq) (*admin.AddAppletResp, error) {
	defer log.ZDebug(ctx, "return")
	if _, err := mctx.CheckAdmin(ctx); err != nil {
		return nil, err
	}
	if req.Name == "" {
		return nil, errs.ErrArgs.Wrap("name empty")
	}
	if req.AppID == "" {
		return nil, errs.ErrArgs.Wrap("appid empty")
	}
	if !(req.Status == constant.StatusOnShelf || req.Status == constant.StatusUnShelf) {
		return nil, errs.ErrArgs.Wrap("invalid status")
	}
	m := admin2.Applet{
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
	if err := o.Database.CreateApplet(ctx, &m); err != nil {
		return nil, err
	}
	return &admin.AddAppletResp{}, nil
}

func (o *adminServer) DelApplet(ctx context.Context, req *admin.DelAppletReq) (*admin.DelAppletResp, error) {
	defer log.ZDebug(ctx, "return")
	if _, err := mctx.CheckAdmin(ctx); err != nil {
		return nil, err
	}
	if len(req.AppletIds) == 0 {
		return nil, errs.ErrArgs.Wrap("AppletIds empty")
	}
	applets, err := o.Database.FindApplet(ctx, req.AppletIds)
	if err != nil {
		return nil, err
	}
	if ids := utils.Single(req.AppletIds, utils.Slice(applets, func(e *admin2.Applet) string { return e.ID })); len(ids) > 0 {
		return nil, errs.ErrArgs.Wrap("ids not found: " + strings.Join(ids, ", "))
	}
	if err := o.Database.DelApplet(ctx, req.AppletIds); err != nil {
		return nil, err
	}
	return &admin.DelAppletResp{}, nil
}

func (o *adminServer) UpdateApplet(ctx context.Context, req *admin.UpdateAppletReq) (*admin.UpdateAppletResp, error) {
	defer log.ZDebug(ctx, "return")
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
	defer log.ZDebug(ctx, "return")
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
	defer log.ZDebug(ctx, "return")
	if _, err := mctx.CheckAdmin(ctx); err != nil {
		return nil, err
	}
	total, applets, err := o.Database.SearchApplet(ctx, req.Keyword, req.Pagination.PageNumber, req.Pagination.ShowNumber)
	if err != nil {
		return nil, err
	}
	resp := &admin.SearchAppletResp{Total: total, Applets: make([]*common.AppletInfo, 0, len(applets))}
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
