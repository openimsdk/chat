package chat

import (
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
	"github.com/OpenIMSDK/chat/pkg/common/db/table"
	"github.com/OpenIMSDK/chat/pkg/proto/common"
)

func DbToPbAttribute(attribute *table.Attribute) *common.UserPublicInfo {
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

func DbToPbAttributes(attributes []*table.Attribute) []*common.UserPublicInfo {
	return utils.Slice(attributes, DbToPbAttribute)
}

func DbToPbUserFullInfo(attribute *table.Attribute) *common.UserFullInfo {
	return &common.UserFullInfo{
		UserID:      attribute.UserID,
		Password:    "",
		Account:     attribute.Account,
		PhoneNumber: attribute.PhoneNumber,
		AreaCode:    attribute.AreaCode,
		Email:       attribute.Email,
		Nickname:    attribute.Nickname,
		FaceURL:     attribute.FaceURL,
		Gender:      attribute.Gender,
		Level:       attribute.Level,
		//Forbidden:      attribute.Forbidden,
		Birth:          attribute.BirthTime.UnixMilli(),
		AllowAddFriend: attribute.AllowAddFriend,
		AllowBeep:      attribute.AllowBeep,
		AllowVibration: attribute.AllowVibration,
	}
}

func DbToPbUserFullInfos(attributes []*table.Attribute) []*common.UserFullInfo {
	return utils.Slice(attributes, DbToPbUserFullInfo)
}
