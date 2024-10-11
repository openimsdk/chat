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
	"strconv"
	"strings"
	"time"

	"github.com/openimsdk/chat/pkg/common/db/dbutil"
	chatdb "github.com/openimsdk/chat/pkg/common/db/table/chat"
	constantpb "github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/mcontext"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/openimsdk/chat/pkg/common/constant"
	"github.com/openimsdk/chat/pkg/common/mctx"
	"github.com/openimsdk/chat/pkg/eerrs"
	"github.com/openimsdk/chat/pkg/protocol/chat"
	"github.com/openimsdk/tools/errs"
)

func (o *chatSvr) checkUpdateInfo(ctx context.Context, req *chat.UpdateUserInfoReq) error {
	attribute, err := o.Database.TakeAttributeByUserID(ctx, req.UserID)
	if err != nil {
		return err
	}
	checkEmail := func() error {
		if req.Email == nil {
			return nil
		}
		if req.Email.Value == attribute.Email {
			req.Email = nil
			return nil
		}
		if req.Email.Value == "" {
			if !(attribute.Account != "" || (attribute.AreaCode != "" && attribute.PhoneNumber != "")) {
				return errs.ErrArgs.WrapMsg("a login method must exist")
			}
			return nil
		} else {
			if _, err := o.Database.GetAttributeByEmail(ctx, req.Email.Value); err == nil {
				return errs.ErrDuplicateKey.WrapMsg("email already exists")
			} else if !dbutil.IsDBNotFound(err) {
				return err
			}
		}
		return nil
	}
	checkPhone := func() error {
		if req.AreaCode == nil {
			return nil
		}
		if req.AreaCode.Value == attribute.AreaCode && req.PhoneNumber.Value == attribute.PhoneNumber {
			req.AreaCode = nil
			req.PhoneNumber = nil
			return nil
		}
		if req.AreaCode.Value == "" || req.PhoneNumber.Value == "" {
			if attribute.Email == "" || attribute.Account == "" {
				return errs.ErrArgs.WrapMsg("a login method must exist")
			}
		} else {
			if _, err := o.Database.GetAttributeByPhone(ctx, req.AreaCode.Value, req.PhoneNumber.Value); err == nil {
				return errs.ErrDuplicateKey.WrapMsg("phone number already exists")
			} else if !dbutil.IsDBNotFound(err) {
				return err
			}
		}
		return nil
	}
	checkAccount := func() error {
		if req.Account == nil {
			return nil
		}
		if req.Account.Value == attribute.Account {
			req.Account = nil
			return nil
		}
		if req.Account.Value == "" {
			if !(attribute.Email == "" && (attribute.AreaCode == "" || attribute.PhoneNumber == "")) {
				return errs.ErrArgs.WrapMsg("a login method must exist")
			}
		} else {
			if _, err := o.Database.GetAttributeByAccount(ctx, req.Account.Value); err == nil {
				return errs.ErrDuplicateKey.WrapMsg("account already exists")
			} else if !dbutil.IsDBNotFound(err) {
				return err
			}
		}
		return nil
	}
	for _, fn := range []func() error{checkEmail, checkPhone, checkAccount} {
		if err := fn(); err != nil {
			return err
		}
	}
	return nil
}

func (o *chatSvr) UpdateUserInfo(ctx context.Context, req *chat.UpdateUserInfoReq) (*chat.UpdateUserInfoResp, error) {
	if req.AreaCode != nil || req.PhoneNumber != nil {
		if !(req.AreaCode != nil && req.PhoneNumber != nil) {
			return nil, errs.ErrArgs.WrapMsg("areaCode and phoneNumber must be set together")
		}
		if req.AreaCode.Value == "" || req.PhoneNumber.Value == "" {
			if req.AreaCode.Value != req.PhoneNumber.Value {
				return nil, errs.ErrArgs.WrapMsg("areaCode and phoneNumber must be set together")
			}
		}
	}
	opUserID, userType, err := mctx.Check(ctx)
	if err != nil {
		return nil, err
	}
	if req.UserID == "" {
		return nil, errs.ErrArgs.WrapMsg("user id is empty")
	}
	switch userType {
	case constant.NormalUser:
		//if req.UserID == "" {
		//	req.UserID = opUserID
		//}
		if req.UserID != opUserID {
			return nil, errs.ErrNoPermission.WrapMsg("only admin can update other user info")
		}
		if req.AreaCode != nil {
			return nil, errs.ErrNoPermission.WrapMsg("areaCode can not be updated")
		}
		if req.PhoneNumber != nil {
			return nil, errs.ErrNoPermission.WrapMsg("phoneNumber can not be updated")
		}
		if req.Account != nil {
			return nil, errs.ErrNoPermission.WrapMsg("account can not be updated")
		}
		if req.Level != nil {
			return nil, errs.ErrNoPermission.WrapMsg("level can not be updated")
		}
	case constant.AdminUser:
	default:
		return nil, errs.ErrNoPermission.WrapMsg("user type error")
	}
	if err := o.checkUpdateInfo(ctx, req); err != nil {
		return nil, err
	}
	update, err := ToDBAttributeUpdate(req)
	if err != nil {
		return nil, err
	}
	if len(update) > 0 {
		if err := o.Database.UpdateUseInfo(ctx, req.UserID, update); err != nil {
			return nil, err
		}
	}
	return &chat.UpdateUserInfoResp{}, nil
}

func (o *chatSvr) FindUserPublicInfo(ctx context.Context, req *chat.FindUserPublicInfoReq) (*chat.FindUserPublicInfoResp, error) {
	if len(req.UserIDs) == 0 {
		return nil, errs.ErrArgs.WrapMsg("UserIDs is empty")
	}
	attributes, err := o.Database.FindAttribute(ctx, req.UserIDs)
	if err != nil {
		return nil, err
	}
	return &chat.FindUserPublicInfoResp{
		Users: DbToPbAttributes(attributes),
	}, nil
}

func (o *chatSvr) AddUserAccount(ctx context.Context, req *chat.AddUserAccountReq) (*chat.AddUserAccountResp, error) {
	if _, _, err := mctx.Check(ctx); err != nil {
		return nil, err
	}

	if req.User.Email == "" && !(req.User.PhoneNumber != "" && req.User.AreaCode != "") && req.User.Account == "" {
		return nil, errs.ErrArgs.WrapMsg("at least one valid account is required")
	}

	if req.User.UserID == "" {
		for i := 0; i < 20; i++ {
			userID := o.genUserID()
			_, err := o.Database.GetUser(ctx, userID)
			if err == nil {
				continue
			} else if dbutil.IsDBNotFound(err) {
				req.User.UserID = userID
				break
			} else {
				return nil, err
			}
		}
		if req.User.UserID == "" {
			return nil, errs.ErrInternalServer.WrapMsg("gen user id failed")
		}
	}

	var (
		credentials     []*chatdb.Credential
		allowChangeRule = datautil.If(req.User.UserType == constant.CommonUser, true, false)
	)

	if req.User.PhoneNumber != "" {
		if !strings.HasPrefix(req.User.AreaCode, "+") {
			req.User.AreaCode = "+" + req.User.AreaCode
		}
		if _, err := strconv.ParseUint(req.User.AreaCode[1:], 10, 64); err != nil {
			return nil, errs.ErrArgs.WrapMsg("area code must be number")
		}
		if _, err := strconv.ParseUint(req.User.PhoneNumber, 10, 64); err != nil {
			return nil, errs.ErrArgs.WrapMsg("phone number must be number")
		}
		_, err := o.Database.TakeAttributeByPhone(ctx, req.User.AreaCode, req.User.PhoneNumber)
		if err == nil {
			return nil, eerrs.ErrPhoneAlreadyRegister.Wrap()
		} else if !dbutil.IsDBNotFound(err) {
			return nil, err
		}
		credentials = append(credentials, &chatdb.Credential{
			UserID:      req.User.UserID,
			Account:     req.User.AreaCode + " " + req.User.PhoneNumber,
			Type:        constant.CredentialPhone,
			AllowChange: allowChangeRule,
		})
	}

	if req.User.Account != "" {
		_, err := o.Database.TakeAttributeByAccount(ctx, req.User.Account)
		if err == nil {
			return nil, eerrs.ErrAccountAlreadyRegister.Wrap()
		} else if !dbutil.IsDBNotFound(err) {
			return nil, err
		}
		credentials = append(credentials, &chatdb.Credential{
			UserID:      req.User.UserID,
			Account:     req.User.Account,
			Type:        constant.CredentialAccount,
			AllowChange: allowChangeRule,
		})
	}

	if req.User.Email != "" {
		_, err := o.Database.TakeAttributeByEmail(ctx, req.User.Email)
		if err == nil {
			return nil, eerrs.ErrEmailAlreadyRegister.Wrap()
		} else if !dbutil.IsDBNotFound(err) {
			return nil, err
		}
		credentials = append(credentials, &chatdb.Credential{
			UserID:      req.User.UserID,
			Account:     req.User.Email,
			Type:        constant.CredentialEmail,
			AllowChange: allowChangeRule,
		})
	}

	register := &chatdb.Register{
		UserID:      req.User.UserID,
		DeviceID:    req.DeviceID,
		IP:          req.Ip,
		Platform:    constantpb.PlatformID2Name[int(req.Platform)],
		AccountType: "",
		Mode:        constant.UserMode,
		CreateTime:  time.Now(),
	}
	account := &chatdb.Account{
		UserID:         req.User.UserID,
		Password:       req.User.Password,
		OperatorUserID: mcontext.GetOpUserID(ctx),
		ChangeTime:     register.CreateTime,
		CreateTime:     register.CreateTime,
	}
	attribute := &chatdb.Attribute{
		UserID:         req.User.UserID,
		Account:        req.User.Account,
		PhoneNumber:    req.User.PhoneNumber,
		AreaCode:       req.User.AreaCode,
		Email:          req.User.Email,
		Nickname:       req.User.Nickname,
		FaceURL:        req.User.FaceURL,
		Gender:         req.User.Gender,
		BirthTime:      time.UnixMilli(req.User.Birth),
		ChangeTime:     register.CreateTime,
		CreateTime:     register.CreateTime,
		AllowVibration: constant.DefaultAllowVibration,
		AllowBeep:      constant.DefaultAllowBeep,
		AllowAddFriend: constant.DefaultAllowAddFriend,
	}

	if req.User.UserType == constant.OrgUser {
		attribute.EnglishName = req.User.EnglishName.GetValuePtr()
		attribute.Station = req.User.Station.GetValuePtr()
		attribute.Telephone = req.User.Telephone.GetValuePtr()
	}
	if err := o.Database.RegisterUser(ctx, register, account, attribute, credentials); err != nil {
		return nil, err
	}
	return &chat.AddUserAccountResp{}, nil
}

func (o *chatSvr) SearchUserPublicInfo(ctx context.Context, req *chat.SearchUserPublicInfoReq) (*chat.SearchUserPublicInfoResp, error) {
	if _, _, err := mctx.Check(ctx); err != nil {
		return nil, err
	}
	total, list, err := o.Database.Search(ctx, constant.FinDAllUser, req.Keyword, req.Genders, req.Pagination)
	if err != nil {
		return nil, err
	}
	return &chat.SearchUserPublicInfoResp{
		Total: uint32(total),
		Users: DbToPbAttributes(list),
	}, nil
}

func (o *chatSvr) FindUserFullInfo(ctx context.Context, req *chat.FindUserFullInfoReq) (*chat.FindUserFullInfoResp, error) {
	if _, _, err := mctx.Check(ctx); err != nil {
		return nil, err
	}
	if len(req.UserIDs) == 0 {
		return nil, errs.ErrArgs.WrapMsg("UserIDs is empty")
	}
	attributes, err := o.Database.FindAttribute(ctx, req.UserIDs)
	if err != nil {
		return nil, err
	}
	return &chat.FindUserFullInfoResp{Users: DbToPbUserFullInfos(attributes)}, nil
}

func (o *chatSvr) SearchUserFullInfo(ctx context.Context, req *chat.SearchUserFullInfoReq) (*chat.SearchUserFullInfoResp, error) {
	if _, _, err := mctx.Check(ctx); err != nil {
		return nil, err
	}
	total, list, err := o.Database.Search(ctx, req.Normal, req.Keyword, req.Genders, req.Pagination)
	if err != nil {
		return nil, err
	}
	return &chat.SearchUserFullInfoResp{
		Total: uint32(total),
		Users: DbToPbUserFullInfos(list),
	}, nil
}

func (o *chatSvr) FindUserAccount(ctx context.Context, req *chat.FindUserAccountReq) (*chat.FindUserAccountResp, error) {
	if len(req.UserIDs) == 0 {
		return nil, errs.ErrArgs.WrapMsg("user id list must be set")
	}
	if _, _, err := mctx.CheckAdminOrUser(ctx); err != nil {
		return nil, err
	}
	attributes, err := o.Database.FindAttribute(ctx, req.UserIDs)
	if err != nil {
		return nil, err
	}
	userAccountMap := make(map[string]string)
	for _, attribute := range attributes {
		userAccountMap[attribute.UserID] = attribute.Account
	}
	return &chat.FindUserAccountResp{UserAccountMap: userAccountMap}, nil
}

func (o *chatSvr) FindAccountUser(ctx context.Context, req *chat.FindAccountUserReq) (*chat.FindAccountUserResp, error) {
	if len(req.Accounts) == 0 {
		return nil, errs.ErrArgs.WrapMsg("account list must be set")
	}
	if _, _, err := mctx.CheckAdminOrUser(ctx); err != nil {
		return nil, err
	}
	attributes, err := o.Database.FindAttribute(ctx, req.Accounts)
	if err != nil {
		return nil, err
	}
	accountUserMap := make(map[string]string)
	for _, attribute := range attributes {
		accountUserMap[attribute.Account] = attribute.UserID
	}
	return &chat.FindAccountUserResp{AccountUserMap: accountUserMap}, nil
}

func (o *chatSvr) SearchUserInfo(ctx context.Context, req *chat.SearchUserInfoReq) (*chat.SearchUserInfoResp, error) {
	if _, _, err := mctx.Check(ctx); err != nil {
		return nil, err
	}
	total, list, err := o.Database.SearchUser(ctx, req.Keyword, req.UserIDs, req.Genders, req.Pagination)
	if err != nil {
		return nil, err
	}
	return &chat.SearchUserInfoResp{
		Total: uint32(total),
		Users: DbToPbUserFullInfos(list),
	}, nil
}

func (o *chatSvr) CheckUserExist(ctx context.Context, req *chat.CheckUserExistReq) (resp *chat.CheckUserExistResp, err error) {
	if req.User.PhoneNumber != "" {
		attributeByPhone, err := o.Database.TakeAttributeByPhone(ctx, req.User.AreaCode, req.User.PhoneNumber)
		// err != nil is not found User
		if err != nil && errs.Unwrap(err) != mongo.ErrNoDocuments {
			return nil, err
		}
		if attributeByPhone != nil {
			log.ZDebug(ctx, "Check Number is ", attributeByPhone.PhoneNumber)
			log.ZDebug(ctx, "Check userID is ", attributeByPhone.UserID)
			if attributeByPhone.PhoneNumber == req.User.PhoneNumber {
				return &chat.CheckUserExistResp{Userid: attributeByPhone.UserID, IsRegistered: true}, nil
			}
		}
	} else {
		if req.User.Email != "" {
			attributeByEmail, err := o.Database.TakeAttributeByEmail(ctx, req.User.Email)
			if err != nil && errs.Unwrap(err) != mongo.ErrNoDocuments {
				return nil, err
			}
			if attributeByEmail != nil {
				log.ZDebug(ctx, "Check email is ", attributeByEmail.Email)
				log.ZDebug(ctx, "Check userID is ", attributeByEmail.UserID)
				if attributeByEmail.Email == req.User.Email {
					return &chat.CheckUserExistResp{Userid: attributeByEmail.UserID, IsRegistered: true}, nil
				}
			}
		}
	}
	return nil, nil
}

func (o *chatSvr) DelUserAccount(ctx context.Context, req *chat.DelUserAccountReq) (resp *chat.DelUserAccountResp, err error) {
	if err := o.Database.DelUserAccount(ctx, req.UserIDs); err != nil && errs.Unwrap(err) != mongo.ErrNoDocuments {
		return nil, err
	}
	return nil, nil
}
