package admin

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/openimsdk/chat/internal/api/util"
	"github.com/openimsdk/chat/pkg/common/apistruct"
	"github.com/openimsdk/chat/pkg/common/config"
	"github.com/openimsdk/chat/pkg/common/imapi"
	"github.com/openimsdk/chat/pkg/common/mctx"
	"github.com/openimsdk/chat/pkg/common/xlsx"
	"github.com/openimsdk/chat/pkg/common/xlsx/model"
	"github.com/openimsdk/chat/pkg/protocol/admin"
	"github.com/openimsdk/chat/pkg/protocol/chat"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/protocol/user"
	"github.com/openimsdk/tools/a2r"
	"github.com/openimsdk/tools/apiresp"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/utils/datautil"
	"github.com/openimsdk/tools/utils/encrypt"
)

func New(chatClient chat.ChatClient, adminClient admin.AdminClient, imApiCaller imapi.CallerInterface, api *util.Api) *Api {
	return &Api{
		Api:         api,
		chatClient:  chatClient,
		adminClient: adminClient,
		imApiCaller: imApiCaller,
	}
}

type Api struct {
	*util.Api
	chatClient  chat.ChatClient
	adminClient admin.AdminClient
	imApiCaller imapi.CallerInterface
}

func (o *Api) AdminLogin(c *gin.Context) {
	req, err := a2r.ParseRequest[admin.LoginReq](c)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	if req.Version == "" {
		apiresp.GinError(c, errs.New("openim-admin-front version too old, please use new version").Wrap())
		return
	}
	loginResp, err := o.adminClient.Login(c, req)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	imAdminUserID := o.GetDefaultIMAdminUserID()
	imToken, err := o.imApiCaller.ImAdminTokenWithDefaultAdmin(c)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	var resp apistruct.AdminLoginResp
	if err := datautil.CopyStructFields(&resp, loginResp); err != nil {
		apiresp.GinError(c, err)
		return
	}
	resp.ImToken = imToken
	resp.ImUserID = imAdminUserID
	apiresp.GinSuccess(c, resp)
}

func (o *Api) ResetUserPassword(c *gin.Context) {
	req, err := a2r.ParseRequest[chat.ChangePasswordReq](c)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	resp, err := o.chatClient.ChangePassword(c, req)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}

	imToken, err := o.imApiCaller.ImAdminTokenWithDefaultAdmin(c)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}

	err = o.imApiCaller.ForceOffLine(mctx.WithApiToken(c, imToken), req.UserID)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}

	apiresp.GinSuccess(c, resp)
}

func (o *Api) AdminUpdateInfo(c *gin.Context) {
	req, err := a2r.ParseRequest[admin.AdminUpdateInfoReq](c)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	resp, err := o.adminClient.AdminUpdateInfo(c, req)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}

	imAdminUserID := o.GetDefaultIMAdminUserID()
	imToken, err := o.imApiCaller.GetAdminTokenCache(c, imAdminUserID)
	if err != nil {
		log.ZError(c, "AdminUpdateInfo ImAdminTokenWithDefaultAdmin", err, "imAdminUserID", imAdminUserID)
		return
	}
	if err := o.imApiCaller.UpdateUserInfo(mctx.WithApiToken(c, imToken), imAdminUserID, resp.Nickname, resp.FaceURL); err != nil {
		log.ZError(c, "AdminUpdateInfo UpdateUserInfo", err, "userID", resp.UserID, "nickName", resp.Nickname, "faceURL", resp.FaceURL)
	}
	apiresp.GinSuccess(c, nil)
}

func (o *Api) AdminInfo(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.GetAdminInfo, o.adminClient)
}

func (o *Api) ChangeAdminPassword(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.ChangeAdminPassword, o.adminClient)
}

func (o *Api) AddAdminAccount(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.AddAdminAccount, o.adminClient)
}

func (o *Api) AddUserAccount(c *gin.Context) {
	req, err := a2r.ParseRequest[chat.AddUserAccountReq](c)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	ip, err := o.GetClientIP(c)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	imToken, err := o.imApiCaller.ImAdminTokenWithDefaultAdmin(c)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}

	ctx := o.WithAdminUser(mctx.WithApiToken(c, imToken))

	err = o.registerChatUser(ctx, ip, []*chat.RegisterUserInfo{req.User})
	if err != nil {
		return
	}

	apiresp.GinSuccess(c, nil)
}

func (o *Api) DelAdminAccount(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.DelAdminAccount, o.adminClient)
}

func (o *Api) SearchAdminAccount(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.SearchAdminAccount, o.adminClient)
}

func (o *Api) AddDefaultFriend(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.AddDefaultFriend, o.adminClient)
}

func (o *Api) DelDefaultFriend(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.DelDefaultFriend, o.adminClient)
}

func (o *Api) SearchDefaultFriend(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.SearchDefaultFriend, o.adminClient)
}

func (o *Api) FindDefaultFriend(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.FindDefaultFriend, o.adminClient)
}

func (o *Api) AddDefaultGroup(c *gin.Context) {
	req, err := a2r.ParseRequest[admin.AddDefaultGroupReq](c)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	imToken, err := o.imApiCaller.ImAdminTokenWithDefaultAdmin(c)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	groups, err := o.imApiCaller.FindGroupInfo(mctx.WithApiToken(c, imToken), req.GroupIDs)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	if len(req.GroupIDs) != len(groups) {
		apiresp.GinError(c, errs.ErrArgs.WrapMsg("group id not found"))
		return
	}
	resp, err := o.adminClient.AddDefaultGroup(c, &admin.AddDefaultGroupReq{
		GroupIDs: req.GroupIDs,
	})
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiresp.GinSuccess(c, resp)
}

func (o *Api) DelDefaultGroup(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.DelDefaultGroup, o.adminClient)
}

func (o *Api) FindDefaultGroup(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.FindDefaultGroup, o.adminClient)
}

func (o *Api) SearchDefaultGroup(c *gin.Context) {
	req, err := a2r.ParseRequest[admin.SearchDefaultGroupReq](c)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	searchResp, err := o.adminClient.SearchDefaultGroup(c, req)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	resp := apistruct.SearchDefaultGroupResp{
		Total:  searchResp.Total,
		Groups: make([]*sdkws.GroupInfo, 0, len(searchResp.GroupIDs)),
	}
	if len(searchResp.GroupIDs) > 0 {
		imToken, err := o.imApiCaller.ImAdminTokenWithDefaultAdmin(c)
		if err != nil {
			apiresp.GinError(c, err)
			return
		}
		groups, err := o.imApiCaller.FindGroupInfo(mctx.WithApiToken(c, imToken), searchResp.GroupIDs)
		if err != nil {
			apiresp.GinError(c, err)
			return
		}
		groupMap := make(map[string]*sdkws.GroupInfo)
		for _, group := range groups {
			groupMap[group.GroupID] = group
		}
		for _, groupID := range searchResp.GroupIDs {
			if group, ok := groupMap[groupID]; ok {
				resp.Groups = append(resp.Groups, group)
			} else {
				resp.Groups = append(resp.Groups, &sdkws.GroupInfo{
					GroupID: groupID,
				})
			}
		}
	}
	apiresp.GinSuccess(c, resp)
}

func (o *Api) AddInvitationCode(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.AddInvitationCode, o.adminClient)
}

func (o *Api) GenInvitationCode(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.GenInvitationCode, o.adminClient)
}

func (o *Api) DelInvitationCode(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.DelInvitationCode, o.adminClient)
}

func (o *Api) SearchInvitationCode(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.SearchInvitationCode, o.adminClient)
}

func (o *Api) AddUserIPLimitLogin(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.AddUserIPLimitLogin, o.adminClient)
}

func (o *Api) SearchUserIPLimitLogin(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.SearchUserIPLimitLogin, o.adminClient)
}

func (o *Api) DelUserIPLimitLogin(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.DelUserIPLimitLogin, o.adminClient)
}

func (o *Api) SearchIPForbidden(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.SearchIPForbidden, o.adminClient)
}

func (o *Api) AddIPForbidden(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.AddIPForbidden, o.adminClient)
}

func (o *Api) DelIPForbidden(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.DelIPForbidden, o.adminClient)
}

func (o *Api) ParseToken(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.ParseToken, o.adminClient)
}

func (o *Api) BlockUser(c *gin.Context) {
	req, err := a2r.ParseRequest[admin.BlockUserReq](c)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	resp, err := o.adminClient.BlockUser(c, req)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	imToken, err := o.imApiCaller.ImAdminTokenWithDefaultAdmin(c)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	err = o.imApiCaller.ForceOffLine(mctx.WithApiToken(c, imToken), req.UserID)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiresp.GinSuccess(c, resp)
}

func (o *Api) UnblockUser(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.UnblockUser, o.adminClient)
}

func (o *Api) SearchBlockUser(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.SearchBlockUser, o.adminClient)
}

func (o *Api) SetClientConfig(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.SetClientConfig, o.adminClient)
}

func (o *Api) DelClientConfig(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.DelClientConfig, o.adminClient)
}

func (o *Api) GetClientConfig(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.GetClientConfig, o.adminClient)
}

func (o *Api) AddApplet(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.AddApplet, o.adminClient)
}

func (o *Api) DelApplet(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.DelApplet, o.adminClient)
}

func (o *Api) UpdateApplet(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.UpdateApplet, o.adminClient)
}

func (o *Api) SearchApplet(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.SearchApplet, o.adminClient)
}

func (o *Api) LoginUserCount(c *gin.Context) {
	a2r.Call(c, chat.ChatClient.UserLoginCount, o.chatClient)
}

func (o *Api) NewUserCount(c *gin.Context) {
	req, err := a2r.ParseRequest[user.UserRegisterCountReq](c)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	imToken, err := o.imApiCaller.ImAdminTokenWithDefaultAdmin(c)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	dateCount, total, err := o.imApiCaller.UserRegisterCount(mctx.WithApiToken(c, imToken), req.Start, req.End)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiresp.GinSuccess(c, &apistruct.NewUserCountResp{
		DateCount: dateCount,
		Total:     total,
	})
}

func (o *Api) ImportUserByXlsx(c *gin.Context) {
	formFile, err := c.FormFile("data")
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	ip, err := o.GetClientIP(c)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	file, err := formFile.Open()
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	defer file.Close()
	var users []model.User
	if err := xlsx.ParseAll(file, &users); err != nil {
		apiresp.GinError(c, errs.ErrArgs.WrapMsg("xlsx file parse error "+err.Error()))
		return
	}
	us, err := o.xlsx2user(users)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	imToken, err := o.imApiCaller.ImAdminTokenWithDefaultAdmin(c)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}

	ctx := o.WithAdminUser(mctx.WithApiToken(c, imToken))
	apiresp.GinError(c, o.registerChatUser(ctx, ip, us))
}

func (o *Api) ImportUserByJson(c *gin.Context) {
	req, err := a2r.ParseRequest[struct {
		Users []*chat.RegisterUserInfo `json:"users"`
	}](c)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	ip, err := o.GetClientIP(c)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	imToken, err := o.imApiCaller.ImAdminTokenWithDefaultAdmin(c)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	ctx := o.WithAdminUser(mctx.WithApiToken(c, imToken))
	apiresp.GinError(c, o.registerChatUser(ctx, ip, req.Users))
}

func (o *Api) xlsx2user(users []model.User) ([]*chat.RegisterUserInfo, error) {
	chatUsers := make([]*chat.RegisterUserInfo, len(users))
	for i, info := range users {
		if info.Nickname == "" {
			return nil, errs.ErrArgs.WrapMsg("nickname is empty")
		}
		if info.AreaCode == "" || info.PhoneNumber == "" {
			return nil, errs.ErrArgs.WrapMsg("areaCode or phoneNumber is empty")
		}
		if info.Password == "" {
			return nil, errs.ErrArgs.WrapMsg("password is empty")
		}
		if !strings.HasPrefix(info.AreaCode, "+") {
			return nil, errs.ErrArgs.WrapMsg("areaCode format error")
		}
		if _, err := strconv.ParseUint(info.AreaCode[1:], 10, 16); err != nil {
			return nil, errs.ErrArgs.WrapMsg("areaCode format error")
		}
		gender, _ := strconv.Atoi(info.Gender)
		chatUsers[i] = &chat.RegisterUserInfo{
			UserID:      info.UserID,
			Nickname:    info.Nickname,
			FaceURL:     info.FaceURL,
			Birth:       o.xlsxBirth(info.Birth).UnixMilli(),
			Gender:      int32(gender),
			AreaCode:    info.AreaCode,
			PhoneNumber: info.PhoneNumber,
			Email:       info.Email,
			Account:     info.Account,
			Password:    encrypt.Md5(info.Password),
		}
	}
	return chatUsers, nil
}

func (o *Api) xlsxBirth(s string) time.Time {
	if s == "" {
		return time.Now()
	}
	var separator byte
	for _, b := range []byte(s) {
		if b < '0' || b > '9' {
			separator = b
		}
	}

	arr := strings.Split(s, string([]byte{separator}))
	if len(arr) != 3 {
		return time.Now()
	}
	year, _ := strconv.Atoi(arr[0])
	month, _ := strconv.Atoi(arr[1])
	day, _ := strconv.Atoi(arr[2])
	t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
	if t.Before(time.Date(1900, 0, 0, 0, 0, 0, 0, time.Local)) {
		return time.Now()
	}
	return t
}

func (o *Api) registerChatUser(ctx context.Context, ip string, users []*chat.RegisterUserInfo) error {
	if len(users) == 0 {
		return errs.ErrArgs.WrapMsg("users is empty")
	}
	for _, info := range users {
		respRegisterUser, err := o.chatClient.RegisterUser(ctx, &chat.RegisterUserReq{Ip: ip, User: info, Platform: constant.AdminPlatformID})
		if err != nil {
			return err
		}
		userInfo := &sdkws.UserInfo{
			UserID:   respRegisterUser.UserID,
			Nickname: info.Nickname,
			FaceURL:  info.FaceURL,
		}
		if err = o.imApiCaller.RegisterUser(ctx, []*sdkws.UserInfo{userInfo}); err != nil {
			return err
		}

		if resp, err := o.adminClient.FindDefaultFriend(ctx, &admin.FindDefaultFriendReq{}); err == nil {
			_ = o.imApiCaller.ImportFriend(ctx, respRegisterUser.UserID, resp.UserIDs)
		}
		if resp, err := o.adminClient.FindDefaultGroup(ctx, &admin.FindDefaultGroupReq{}); err == nil {
			_ = o.imApiCaller.InviteToGroup(ctx, respRegisterUser.UserID, resp.GroupIDs)
		}
	}
	return nil
}

func (o *Api) BatchImportTemplate(c *gin.Context) {
	md5Sum := md5.Sum(config.ImportTemplate)
	md5Val := hex.EncodeToString(md5Sum[:])
	if c.GetHeader("If-None-Match") == md5Val {
		c.Status(http.StatusNotModified)
		return
	}
	c.Header("Content-Disposition", "attachment; filename=template.xlsx")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Length", strconv.Itoa(len(config.ImportTemplate)))
	c.Header("ETag", md5Val)
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", config.ImportTemplate)
}

func (o *Api) SetAllowRegister(c *gin.Context) {
	a2r.Call(c, chat.ChatClient.SetAllowRegister, o.chatClient)
}

func (o *Api) GetAllowRegister(c *gin.Context) {
	a2r.Call(c, chat.ChatClient.GetAllowRegister, o.chatClient)
}

func (o *Api) LatestApplicationVersion(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.LatestApplicationVersion, o.adminClient)
}

func (o *Api) PageApplicationVersion(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.PageApplicationVersion, o.adminClient)
}

func (o *Api) AddApplicationVersion(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.AddApplicationVersion, o.adminClient)
}

func (o *Api) UpdateApplicationVersion(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.UpdateApplicationVersion, o.adminClient)
}

func (o *Api) DeleteApplicationVersion(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.DeleteApplicationVersion, o.adminClient)
}
