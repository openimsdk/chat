package chat

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	constantpb "github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/tools/utils/datautil"

	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/mcontext"

	"github.com/openimsdk/chat/pkg/common/constant"
	"github.com/openimsdk/chat/pkg/common/db/dbutil"
	chatdb "github.com/openimsdk/chat/pkg/common/db/table/chat"
	"github.com/openimsdk/chat/pkg/eerrs"
	"github.com/openimsdk/chat/pkg/protocol/chat"
)

type verifyType int

const (
	phone verifyType = iota
	mail
)

func (o *chatSvr) verifyCodeJoin(areaCode, phoneNumber string) string {
	return areaCode + " " + phoneNumber
}

func (o *chatSvr) SendVerifyCode(ctx context.Context, req *chat.SendVerifyCodeReq) (*chat.SendVerifyCodeResp, error) {
	switch int(req.UsedFor) {
	case constant.VerificationCodeForRegister:
		if err := o.Admin.CheckRegister(ctx, req.Ip); err != nil {
			return nil, err
		}
		if req.Email == "" {
			if req.AreaCode == "" || req.PhoneNumber == "" {
				return nil, errs.ErrArgs.WrapMsg("area code or phone number is empty")
			}
			if !strings.HasPrefix(req.AreaCode, "+") {
				req.AreaCode = "+" + req.AreaCode
			}
			if _, err := strconv.ParseUint(req.AreaCode[1:], 10, 64); err != nil {
				return nil, errs.ErrArgs.WrapMsg("area code must be number")
			}
			if _, err := strconv.ParseUint(req.PhoneNumber, 10, 64); err != nil {
				return nil, errs.ErrArgs.WrapMsg("phone number must be number")
			}
		} else {
			if err := chat.EmailCheck(req.Email); err != nil {
				return nil, errs.ErrArgs.WrapMsg("email must be right")
			}
		}
		conf, err := o.Admin.GetConfig(ctx)
		if err != nil {
			return nil, err
		}
		if val := conf[constant.NeedInvitationCodeRegisterConfigKey]; datautil.Contain(strings.ToLower(val), "1", "true", "yes") {
			if req.InvitationCode == "" {
				return nil, errs.ErrArgs.WrapMsg("invitation code is empty")
			}
			if err := o.Admin.CheckInvitationCode(ctx, req.InvitationCode); err != nil {
				return nil, err
			}
		}
	case constant.VerificationCodeForLogin, constant.VerificationCodeForResetPassword:
		if req.Email == "" {
			_, err := o.Database.TakeAttributeByPhone(ctx, req.AreaCode, req.PhoneNumber)
			if dbutil.IsDBNotFound(err) {
				return nil, eerrs.ErrAccountNotFound.WrapMsg("phone unregistered")
			} else if err != nil {
				return nil, err
			}
		} else {
			_, err := o.Database.TakeAttributeByEmail(ctx, req.Email)
			if dbutil.IsDBNotFound(err) {
				return nil, eerrs.ErrAccountNotFound.WrapMsg("email unregistered")
			} else if err != nil {
				return nil, err
			}
		}

	default:
		return nil, errs.ErrArgs.WrapMsg("used unknown")
	}
	if o.SMS == nil && o.Mail == nil {
		return &chat.SendVerifyCodeResp{}, nil // super code
	}
	if req.Email != "" {
		switch o.conf.Mail.Use {
		case constant.VerifySuperCode:
			return &chat.SendVerifyCodeResp{}, nil // super code
		case constant.VerifyMail:
		default:
			return nil, errs.ErrInternalServer.WrapMsg("email verification code is not enabled")
		}
	}

	if req.AreaCode != "" {
		switch o.conf.Phone.Use {
		case constant.VerifySuperCode:
			return &chat.SendVerifyCodeResp{}, nil // super code
		case constant.VerifyALi:
		default:
			return nil, errs.ErrInternalServer.WrapMsg("phone verification code is not enabled")
		}
	}

	isEmail := req.Email != ""
	var (
		code     = o.genVerifyCode()
		account  string
		sendCode func() error
	)
	if isEmail {
		sendCode = func() error {
			return o.Mail.SendMail(ctx, req.Email, code)
		}
		account = req.Email
	} else {
		sendCode = func() error {
			return o.SMS.SendCode(ctx, req.AreaCode, req.PhoneNumber, code)
		}
		account = o.verifyCodeJoin(req.AreaCode, req.PhoneNumber)
	}
	now := time.Now()
	count, err := o.Database.CountVerifyCodeRange(ctx, account, now.Add(-o.Code.UintTime), now)
	if err != nil {
		return nil, err
	}
	if o.Code.MaxCount < int(count) {
		return nil, eerrs.ErrVerifyCodeSendFrequently.Wrap()
	}
	platformName := constantpb.PlatformIDToName(int(req.Platform))
	if platformName == "" {
		platformName = fmt.Sprintf("platform:%d", req.Platform)
	}
	vc := &chatdb.VerifyCode{
		Account:    account,
		Code:       code,
		Platform:   platformName,
		Duration:   uint(o.Code.ValidTime / time.Second),
		Count:      0,
		Used:       false,
		CreateTime: now,
	}
	if err := o.Database.AddVerifyCode(ctx, vc, sendCode); err != nil {
		return nil, err
	}
	log.ZDebug(ctx, "send code success", "account", account, "code", code, "platform", platformName)
	return &chat.SendVerifyCodeResp{}, nil
}

func (o *chatSvr) verifyCode(ctx context.Context, account string, verifyCode string, type_ verifyType) (string, error) {
	if verifyCode == "" {
		return "", errs.ErrArgs.WrapMsg("verify code is empty")
	}
	switch type_ {
	case phone:
		switch o.conf.Phone.Use {
		case constant.VerifySuperCode:
			if o.Code.SuperCode != verifyCode {
				return "", eerrs.ErrVerifyCodeNotMatch.Wrap()
			}
			return "", nil
		case constant.VerifyALi:
		default:
			return "", errs.ErrInternalServer.WrapMsg("phone verification code is not enabled", "use", o.conf.Phone.Use)
		}
	case mail:
		switch o.conf.Mail.Use {
		case constant.VerifySuperCode:
			if o.Code.SuperCode != verifyCode {
				return "", eerrs.ErrVerifyCodeNotMatch.Wrap()
			}
			return "", nil
		case constant.VerifyMail:
		default:
			return "", errs.ErrInternalServer.WrapMsg("email verification code is not enabled")
		}
	}

	last, err := o.Database.TakeLastVerifyCode(ctx, account)
	if err != nil {
		if dbutil.IsDBNotFound(err) {
			return "", eerrs.ErrVerifyCodeExpired.Wrap()
		}
		return "", err
	}
	if last.CreateTime.Unix()+int64(last.Duration) < time.Now().Unix() {
		return last.ID, eerrs.ErrVerifyCodeExpired.Wrap()
	}
	if last.Used {
		return last.ID, eerrs.ErrVerifyCodeUsed.Wrap()
	}
	if n := o.Code.ValidCount; n > 0 {
		if last.Count >= n {
			return last.ID, eerrs.ErrVerifyCodeMaxCount.Wrap()
		}
		if last.Code != verifyCode {
			if err := o.Database.UpdateVerifyCodeIncrCount(ctx, last.ID); err != nil {
				return last.ID, err
			}
		}
	}
	if last.Code != verifyCode {
		return last.ID, eerrs.ErrVerifyCodeNotMatch.Wrap()
	}
	return last.ID, nil
}

func (o *chatSvr) VerifyCode(ctx context.Context, req *chat.VerifyCodeReq) (*chat.VerifyCodeResp, error) {
	var account string
	if req.PhoneNumber != "" {
		account = o.verifyCodeJoin(req.AreaCode, req.PhoneNumber)
		if _, err := o.verifyCode(ctx, account, req.VerifyCode, phone); err != nil {
			return nil, err
		}
	} else {
		account = req.Email
		if _, err := o.verifyCode(ctx, account, req.VerifyCode, mail); err != nil {
			return nil, err
		}
	}

	return &chat.VerifyCodeResp{}, nil
}

func (o *chatSvr) genUserID() string {
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

func (o *chatSvr) genVerifyCode() string {
	data := make([]byte, o.Code.Len)
	rand.Read(data)
	chars := []byte("0123456789")
	for i := 0; i < len(data); i++ {
		data[i] = chars[data[i]%10]
	}
	return string(data)
}

func (o *chatSvr) RegisterUser(ctx context.Context, req *chat.RegisterUserReq) (*chat.RegisterUserResp, error) {
	isAdmin, err := o.Admin.CheckNilOrAdmin(ctx)
	ctx = o.WithAdminUser(ctx)
	if err != nil {
		return nil, err
	}
	if err = o.checkRegisterInfo(ctx, req.User, isAdmin); err != nil {
		return nil, err
	}
	var usedInvitationCode bool
	if !isAdmin {
		if !o.AllowRegister {
			return nil, errs.ErrNoPermission.WrapMsg("register user is disabled")
		}
		if req.User.UserID != "" {
			return nil, errs.ErrNoPermission.WrapMsg("only admin can set user id")
		}
		if err := o.Admin.CheckRegister(ctx, req.Ip); err != nil {
			return nil, err
		}
		conf, err := o.Admin.GetConfig(ctx)
		if err != nil {
			return nil, err
		}
		if val := conf[constant.NeedInvitationCodeRegisterConfigKey]; datautil.Contain(strings.ToLower(val), "1", "true", "yes") {
			usedInvitationCode = true
			if req.InvitationCode == "" {
				return nil, errs.ErrArgs.WrapMsg("invitation code is empty")
			}
			if err := o.Admin.CheckInvitationCode(ctx, req.InvitationCode); err != nil {
				return nil, err
			}
		}
		if req.User.Email == "" {
			if _, err := o.verifyCode(ctx, o.verifyCodeJoin(req.User.AreaCode, req.User.PhoneNumber), req.VerifyCode, phone); err != nil {
				return nil, err
			}
		} else {
			if _, err := o.verifyCode(ctx, req.User.Email, req.VerifyCode, mail); err != nil {
				return nil, err
			}
		}
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
	} else {
		_, err := o.Database.GetUser(ctx, req.User.UserID)
		if err == nil {
			return nil, errs.ErrArgs.WrapMsg("appoint user id already register")
		} else if !dbutil.IsDBNotFound(err) {
			return nil, err
		}
	}
	var (
		credentials  []*chatdb.Credential
		registerType int32
	)

	if req.User.PhoneNumber != "" {
		registerType = constant.PhoneRegister
		credentials = append(credentials, &chatdb.Credential{
			UserID:      req.User.UserID,
			Account:     BuildCredentialPhone(req.User.AreaCode, req.User.PhoneNumber),
			Type:        constant.CredentialPhone,
			AllowChange: true,
		})
	}

	if req.User.Account != "" {
		credentials = append(credentials, &chatdb.Credential{
			UserID:      req.User.UserID,
			Account:     req.User.Account,
			Type:        constant.CredentialAccount,
			AllowChange: true,
		})
		registerType = constant.AccountRegister
	}

	if req.User.Email != "" {
		registerType = constant.EmailRegister
		credentials = append(credentials, &chatdb.Credential{
			UserID:      req.User.UserID,
			Account:     req.User.Email,
			Type:        constant.CredentialEmail,
			AllowChange: true,
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
		RegisterType:   registerType,
	}
	if err := o.Database.RegisterUser(ctx, register, account, attribute, credentials); err != nil {
		return nil, err
	}
	if usedInvitationCode {
		if err := o.Admin.UseInvitationCode(ctx, req.User.UserID, req.InvitationCode); err != nil {
			log.ZError(ctx, "UseInvitationCode", err, "userID", req.User.UserID, "invitationCode", req.InvitationCode)
		}
	}
	var resp chat.RegisterUserResp
	if req.AutoLogin {
		chatToken, err := o.Admin.CreateToken(ctx, req.User.UserID, constant.NormalUser)
		if err == nil {
			resp.ChatToken = chatToken.Token
		} else {
			log.ZError(ctx, "Admin CreateToken Failed", err, "userID", req.User.UserID, "platform", req.Platform)
		}
	}
	resp.UserID = req.User.UserID
	return &resp, nil
}

func (o *chatSvr) Login(ctx context.Context, req *chat.LoginReq) (*chat.LoginResp, error) {
	resp := &chat.LoginResp{}
	if req.Password == "" && req.VerifyCode == "" {
		return nil, errs.ErrArgs.WrapMsg("password or code must be set")
	}
	var (
		err        error
		credential *chatdb.Credential
		acc        string
	)

	switch {
	case req.Account != "":
		acc = req.Account
	case req.PhoneNumber != "":
		if req.AreaCode == "" {
			return nil, errs.ErrArgs.WrapMsg("area code must")
		}
		if !strings.HasPrefix(req.AreaCode, "+") {
			req.AreaCode = "+" + req.AreaCode
		}
		if _, err := strconv.ParseUint(req.AreaCode[1:], 10, 64); err != nil {
			return nil, errs.ErrArgs.WrapMsg("area code must be number")
		}
		acc = BuildCredentialPhone(req.AreaCode, req.PhoneNumber)
	case req.Email != "":
		acc = req.Email
	default:
		return nil, errs.ErrArgs.WrapMsg("account or phone number or email must be set")
	}
	credential, err = o.Database.TakeCredentialByAccount(ctx, acc)
	if err != nil {
		if dbutil.IsDBNotFound(err) {
			return nil, eerrs.ErrAccountNotFound.WrapMsg("user unregistered")
		}
		return nil, err
	}
	if err := o.Admin.CheckLogin(ctx, credential.UserID, req.Ip); err != nil {
		return nil, err
	}
	var verifyCodeID *string
	if req.Password == "" {
		var (
			id string
		)

		if req.Email == "" {
			account := o.verifyCodeJoin(req.AreaCode, req.PhoneNumber)
			id, err = o.verifyCode(ctx, account, req.VerifyCode, phone)
			if err != nil {
				return nil, err
			}
		} else {
			account := req.Email
			id, err = o.verifyCode(ctx, account, req.VerifyCode, mail)
			if err != nil {
				return nil, err
			}
		}

		if id != "" {
			verifyCodeID = &id
		}
	} else {
		account, err := o.Database.TakeAccount(ctx, credential.UserID)
		if err != nil {
			return nil, err
		}
		if account.Password != req.Password {
			return nil, eerrs.ErrPassword.Wrap()
		}
	}
	chatToken, err := o.Admin.CreateToken(ctx, credential.UserID, constant.NormalUser)
	if err != nil {
		return nil, err
	}
	record := &chatdb.UserLoginRecord{
		UserID:    credential.UserID,
		LoginTime: time.Now(),
		IP:        req.Ip,
		DeviceID:  req.DeviceID,
		Platform:  constantpb.PlatformIDToName(int(req.Platform)),
	}
	if err := o.Database.LoginRecord(ctx, record, verifyCodeID); err != nil {
		return nil, err
	}
	if verifyCodeID != nil {
		if err := o.Database.DelVerifyCode(ctx, *verifyCodeID); err != nil {
			return nil, err
		}
	}
	resp.UserID = credential.UserID
	resp.ChatToken = chatToken.Token
	return resp, nil
}
