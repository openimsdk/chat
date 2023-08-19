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
	"github.com/OpenIMSDK/tools/utils"

	"github.com/OpenIMSDK/chat/pkg/common/db/table/chat"
	"github.com/OpenIMSDK/chat/pkg/proto/common"
)

func DbToPbAttribute(attribute *chat.Attribute) *common.UserPublicInfo {
	if attribute == nil {
		return nil
	}
	return &common.UserPublicInfo{
		UserID:   attribute.UserID,
		Account:  attribute.Account,
		Email:    attribute.Email,
		Nickname: attribute.Nickname,
		FaceURL:  attribute.FaceURL,
		Gender:   attribute.Gender,
		Level:    attribute.Level,
	}
}

func DbToPbAttributes(attributes []*chat.Attribute) []*common.UserPublicInfo {
	return utils.Slice(attributes, DbToPbAttribute)
}

func DbToPbUserFullInfo(attribute *chat.Attribute) *common.UserFullInfo {
	return &common.UserFullInfo{
		UserID:           attribute.UserID,
		Password:         "",
		Account:          attribute.Account,
		PhoneNumber:      attribute.PhoneNumber,
		AreaCode:         attribute.AreaCode,
		Email:            attribute.Email,
		Nickname:         attribute.Nickname,
		FaceURL:          attribute.FaceURL,
		Gender:           attribute.Gender,
		Level:            attribute.Level,
		Birth:            attribute.BirthTime.UnixMilli(),
		AllowAddFriend:   attribute.AllowAddFriend,
		AllowBeep:        attribute.AllowBeep,
		AllowVibration:   attribute.AllowVibration,
		GlobalRecvMsgOpt: attribute.GlobalRecvMsgOpt,
	}
}

func DbToPbUserFullInfos(attributes []*chat.Attribute) []*common.UserFullInfo {
	return utils.Slice(attributes, DbToPbUserFullInfo)
}

func DbToPbLogInfo(log *chat.Log) *common.LogInfo {
	return &common.LogInfo{
		Filename:   log.FileName,
		UserID:     log.UserID,
		Platform:   utils.StringToInt32(log.Platform),
		Url:        log.Url,
		CreateTime: log.CreateTime.UnixMilli(),
		LogID:      log.LogID,
		SystemType: log.SystemType,
		Version:    log.Version,
		Ex:         log.Ex,
	}
}

func DbToPbLogInfos(logs []*chat.Log) []*common.LogInfo {
	return utils.Slice(logs, DbToPbLogInfo)
}
