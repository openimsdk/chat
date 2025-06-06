package chat

import (
	"context"
	"strconv"
	"strings"

	"github.com/openimsdk/chat/pkg/common/db/dbutil"
	table "github.com/openimsdk/chat/pkg/common/db/table/chat"
	"github.com/openimsdk/chat/pkg/eerrs"
	"github.com/openimsdk/chat/pkg/protocol/chat"
	"github.com/openimsdk/chat/pkg/protocol/common"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/utils/datautil"
	"github.com/openimsdk/tools/utils/stringutil"
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
	return datautil.Slice(attributes, DbToPbAttribute)
}

func DbToPbUserFullInfo(attribute *table.Attribute) *common.UserFullInfo {
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
		RegisterType:     attribute.RegisterType,
	}
}

func DbToPbUserFullInfos(attributes []*table.Attribute) []*common.UserFullInfo {
	return datautil.Slice(attributes, DbToPbUserFullInfo)
}

func BuildCredentialPhone(areaCode, phone string) string {
	return areaCode + " " + phone
}

func (o *chatSvr) checkRegisterInfo(ctx context.Context, user *chat.RegisterUserInfo, isAdmin bool) error {
	if user == nil {
		return errs.ErrArgs.WrapMsg("user is nil")
	}
	user.Account = strings.TrimSpace(user.Account)
	if user.Email == "" && !(user.PhoneNumber != "" && user.AreaCode != "") && (!isAdmin || user.Account == "") {
		return errs.ErrArgs.WrapMsg("at least one valid account is required")
	}
	if user.PhoneNumber != "" {
		if !strings.HasPrefix(user.AreaCode, "+") {
			user.AreaCode = "+" + user.AreaCode
		}
		if _, err := strconv.ParseUint(user.AreaCode[1:], 10, 64); err != nil {
			return errs.ErrArgs.WrapMsg("area code must be number")
		}
		if _, err := strconv.ParseUint(user.PhoneNumber, 10, 64); err != nil {
			return errs.ErrArgs.WrapMsg("phone number must be number")
		}
		_, err := o.Database.TakeAttributeByPhone(ctx, user.AreaCode, user.PhoneNumber)
		if err == nil {
			return eerrs.ErrPhoneAlreadyRegister.Wrap()
		} else if !dbutil.IsDBNotFound(err) {
			return err
		}
	}
	if user.Account != "" {
		if !stringutil.IsAlphanumeric(user.Account) {
			return errs.ErrArgs.WrapMsg("account must be alphanumeric")
		}
		_, err := o.Database.TakeAttributeByAccount(ctx, user.Account)
		if err == nil {
			return eerrs.ErrAccountAlreadyRegister.Wrap()
		} else if !dbutil.IsDBNotFound(err) {
			return err
		}
	}
	if user.Email != "" {
		if !stringutil.IsValidEmail(user.Email) {
			return errs.ErrArgs.WrapMsg("invalid email")
		}
		_, err := o.Database.TakeAttributeByAccount(ctx, user.Email)
		if err == nil {
			return eerrs.ErrEmailAlreadyRegister.Wrap()
		} else if !dbutil.IsDBNotFound(err) {
			return err
		}
	}
	return nil
}
