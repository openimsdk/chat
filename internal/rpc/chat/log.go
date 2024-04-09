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

package chat

import (
	"context"
	"github.com/openimsdk/tools/utils/datautil"
	"math/rand"
	"time"

	"github.com/openimsdk/chat/pkg/common/constant"
	table "github.com/openimsdk/chat/pkg/common/db/table/chat"
	"github.com/openimsdk/chat/pkg/common/mctx"
	"github.com/openimsdk/chat/pkg/proto/chat"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/mw/specialerror"
	"gorm.io/gorm/utils"
)

func (o *chatSvr) genLogID() string {
	const l = 10
	data := make([]byte, l)
	rand.Read(data)
	chars := []byte("0123456789")
	for i := 0; i < len(data); i++ {
		if i == 0 {
			data[i] = chars[1:][data[i]%9]
		} else {
			data[i] = chars[data[i]%10]
		}
	}
	return string(data)
}

func (o *chatSvr) IsNotFound(err error) bool {
	return errs.ErrRecordNotFound.Is(specialerror.ErrCode(errs.Unwrap(err)))
}

func (o *chatSvr) UploadLogs(ctx context.Context, req *chat.UploadLogsReq) (*chat.UploadLogsResp, error) {
	var DBlogs []*table.Log
	userID, _, err := mctx.Check(ctx)
	if _, err := o.Database.GetUser(ctx, userID); err != nil {
		return nil, err
	}
	for _, fileURL := range req.FileURLs {
		log := table.Log{
			Version:    req.Version,
			SystemType: req.SystemType,
			Platform:   utils.ToString(req.Platform),
			UserID:     userID,
			CreateTime: time.Now(),
			Url:        fileURL.URL,
			FileName:   fileURL.Filename,
		}
		for i := 0; i < 20; i++ {
			id := o.genLogID()
			logs, err := o.Database.GetLogs(ctx, []string{id}, "")
			if err != nil {
				return nil, err
			}
			if len(logs) == 0 {
				log.LogID = id
				break
			}
		}
		if log.LogID == "" {
			return nil, errs.ErrInternalServer.WrapMsg("log id gen error")
		}
		DBlogs = append(DBlogs, &log)
	}
	err = o.Database.UploadLogs(ctx, DBlogs)
	if err != nil {
		return nil, err
	}
	return &chat.UploadLogsResp{}, nil
}

func (o *chatSvr) DeleteLogs(ctx context.Context, req *chat.DeleteLogsReq) (*chat.DeleteLogsResp, error) {
	userID, userType, err := mctx.Check(ctx)
	if err != nil {
		return nil, err
	}
	if userType == constant.AdminUser {
		userID = ""
	}
	logs, err := o.Database.GetLogs(ctx, req.LogIDs, userID)
	if err != nil {
		return nil, err
	}
	var logIDs []string
	for _, log := range logs {
		logIDs = append(logIDs, log.LogID)
	}
	if ids := datautil.Single(req.LogIDs, logIDs); len(ids) > 0 {
		return nil, errs.ErrRecordNotFound.WrapMsg("logIDs not found", "logIDs", ids)
	}
	err = o.Database.DeleteLogs(ctx, req.LogIDs, userID)
	if err != nil {
		return nil, err
	}
	return &chat.DeleteLogsResp{}, nil
}

func (o *chatSvr) SearchLogs(ctx context.Context, req *chat.SearchLogsReq) (*chat.SearchLogsResp, error) {
	var (
		resp    chat.SearchLogsResp
		userIDs []string
	)
	if req.StartTime > req.EndTime {
		return nil, errs.ErrArgs.WrapMsg("startTime>endTime")
	}
	total, logs, err := o.Database.SearchLogs(ctx, req.Keyword, time.UnixMilli(req.StartTime), time.UnixMilli(req.EndTime), req.Pagination)
	if err != nil {
		return nil, err
	}
	pbLogs := DbToPbLogInfos(logs)
	for _, log := range logs {
		userIDs = append(userIDs, log.UserID)
	}
	attributes, err := o.Database.FindAttribute(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	IDtoName := make(map[string]string)
	for _, attribute := range attributes {
		IDtoName[attribute.UserID] = attribute.Nickname
	}
	for _, pbLog := range pbLogs {
		pbLog.Nickname = IDtoName[pbLog.UserID]
	}
	resp.LogsInfos = pbLogs
	resp.Total = uint32(total)
	return &resp, nil
}
