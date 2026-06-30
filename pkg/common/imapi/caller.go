package imapi

import (
	"context"
	"sync"
	"time"

	"github.com/openimsdk/chat/pkg/botstruct"
	"github.com/openimsdk/chat/pkg/eerrs"
	"github.com/openimsdk/protocol/auth"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/protocol/group"
	"github.com/openimsdk/protocol/relation"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/protocol/user"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
)

type CallerInterface interface {
	ImAdminTokenWithDefaultAdmin(ctx context.Context) (string, error)
	ImportFriend(ctx context.Context, ownerUserID string, friendUserID []string) error
	GetUserToken(ctx context.Context, userID string, platform int32) (string, error)
	GetAdminTokenCache(ctx context.Context, userID string) (string, error)
	GetAdminTokenServer(ctx context.Context, userID string) (string, error)
	InviteToGroup(ctx context.Context, userID string, groupIDs []string) error

	UpdateUserInfo(ctx context.Context, userID string, nickName string, faceURL string) error
	GetUserInfo(ctx context.Context, userID string) (*sdkws.UserInfo, error)
	GetUsersInfo(ctx context.Context, userIDs []string) ([]*sdkws.UserInfo, error)
	AddNotificationAccount(ctx context.Context, req *user.AddNotificationAccountReq) error
	UpdateNotificationAccount(ctx context.Context, req *user.UpdateNotificationAccountInfoReq) error

	ForceOffLine(ctx context.Context, userID string) error
	RegisterUser(ctx context.Context, users []*sdkws.UserInfo) error
	FindGroupInfo(ctx context.Context, groupIDs []string) ([]*sdkws.GroupInfo, error)
	UserRegisterCount(ctx context.Context, start int64, end int64) (map[string]int64, int64, error)
	FriendUserIDs(ctx context.Context, userID string) ([]string, error)
	AccountCheckSingle(ctx context.Context, userID string) (bool, error)
	SendSimpleMsg(ctx context.Context, req *SendSingleMsgReq, key string) error
}

type authToken struct {
	token   string
	expired time.Time
}

type Caller struct {
	imApi           string
	imSecret        string
	defaultIMUserID string
	tokenCache      map[string]*authToken
	lock            sync.RWMutex
}

func New(imApi string, imSecret string, defaultIMUserID string) CallerInterface {
	return &Caller{
		imApi:           imApi,
		imSecret:        imSecret,
		defaultIMUserID: defaultIMUserID,
		tokenCache:      make(map[string]*authToken),
		lock:            sync.RWMutex{},
	}
}

func (c *Caller) ImportFriend(ctx context.Context, ownerUserID string, friendUserIDs []string) error {
	if len(friendUserIDs) == 0 {
		return nil
	}
	_, err := importFriend.Call(ctx, c.imApi, &relation.ImportFriendReq{
		OwnerUserID:   ownerUserID,
		FriendUserIDs: friendUserIDs,
	})
	return err
}

func (c *Caller) ImAdminTokenWithDefaultAdmin(ctx context.Context) (string, error) {
	return c.GetAdminTokenCache(ctx, c.defaultIMUserID)
}

func (c *Caller) GetAdminTokenCache(ctx context.Context, userID string) (string, error) {
	c.lock.RLock()
	t, ok := c.tokenCache[userID]
	c.lock.RUnlock()
	if ok && t.expired.After(time.Now()) {
		return t.token, nil
	}
	c.lock.Lock()
	defer c.lock.Unlock()
	t, ok = c.tokenCache[userID]
	if ok && t.expired.After(time.Now()) {
		return t.token, nil
	}
	token, err := c.GetAdminTokenServer(ctx, userID)
	if err != nil {
		return "", err
	}
	c.tokenCache[userID] = &authToken{token: token, expired: time.Now().Add(time.Minute * 4)}
	return token, nil
}

func (c *Caller) GetAdminTokenServer(ctx context.Context, userID string) (string, error) {
	resp, err := getAdminToken.Call(ctx, c.imApi, &auth.GetAdminTokenReq{
		Secret: c.imSecret,
		UserID: userID,
	})
	if err != nil {
		return "", err
	}
	log.ZDebug(ctx, "get im admin token from server", "userID", userID, "token", resp.Token)
	return resp.Token, nil
}

func (c *Caller) GetUserToken(ctx context.Context, userID string, platformID int32) (string, error) {
	resp, err := getuserToken.Call(ctx, c.imApi, &auth.GetUserTokenReq{
		PlatformID: platformID,
		UserID:     userID,
	})
	if err != nil {
		return "", err
	}
	return resp.Token, nil
}

func (c *Caller) InviteToGroup(ctx context.Context, userID string, groupIDs []string) error {
	for _, groupID := range groupIDs {
		_, _ = inviteToGroup.Call(ctx, c.imApi, &group.InviteUserToGroupReq{
			GroupID:        groupID,
			Reason:         "",
			InvitedUserIDs: []string{userID},
		})
	}
	return nil
}

func (c *Caller) UpdateUserInfo(ctx context.Context, userID string, nickName string, faceURL string) error {
	_, err := updateUserInfo.Call(ctx, c.imApi, &user.UpdateUserInfoReq{UserInfo: &sdkws.UserInfo{
		UserID:   userID,
		Nickname: nickName,
		FaceURL:  faceURL,
	}})
	return err
}

func (c *Caller) GetUserInfo(ctx context.Context, userID string) (*sdkws.UserInfo, error) {
	resp, err := c.GetUsersInfo(ctx, []string{userID})
	if err != nil {
		return nil, err
	}
	if len(resp) == 0 {
		return nil, errs.ErrRecordNotFound.WrapMsg("record not found")
	}
	return resp[0], nil
}

func (c *Caller) GetUsersInfo(ctx context.Context, userIDs []string) ([]*sdkws.UserInfo, error) {
	resp, err := getUserInfo.Call(ctx, c.imApi, &user.GetDesignateUsersReq{
		UserIDs: userIDs,
	})
	if err != nil {
		return nil, err
	}
	return resp.UsersInfo, nil
}

func (c *Caller) RegisterUser(ctx context.Context, users []*sdkws.UserInfo) error {
	_, err := registerUser.Call(ctx, c.imApi, &user.UserRegisterReq{
		Users: users,
	})
	return err
}

func (c *Caller) ForceOffLine(ctx context.Context, userID string) error {
	for id := range constant.PlatformID2Name {
		_, _ = forceOffLine.Call(ctx, c.imApi, &auth.ForceLogoutReq{
			PlatformID: int32(id),
			UserID:     userID,
		})
	}
	return nil
}

func (c *Caller) FindGroupInfo(ctx context.Context, groupIDs []string) ([]*sdkws.GroupInfo, error) {
	resp, err := getGroupsInfo.Call(ctx, c.imApi, &group.GetGroupsInfoReq{
		GroupIDs: groupIDs,
	})
	if err != nil {
		return nil, err
	}
	return resp.GroupInfos, nil
}

func (c *Caller) UserRegisterCount(ctx context.Context, start int64, end int64) (map[string]int64, int64, error) {
	resp, err := registerUserCount.Call(ctx, c.imApi, &user.UserRegisterCountReq{
		Start: start,
		End:   end,
	})
	if err != nil {
		return nil, 0, err
	}
	return resp.Count, resp.Total, nil
}

func (c *Caller) FriendUserIDs(ctx context.Context, userID string) ([]string, error) {
	resp, err := friendUserIDs.Call(ctx, c.imApi, &relation.GetFriendIDsReq{UserID: userID})
	if err != nil {
		return nil, err
	}
	return resp.FriendIDs, nil
}

// return true when isUserNotExist.
func (c *Caller) AccountCheckSingle(ctx context.Context, userID string) (bool, error) {
	resp, err := accountCheck.Call(ctx, c.imApi, &user.AccountCheckReq{CheckUserIDs: []string{userID}})
	if err != nil {
		return false, err
	}
	if resp.Results[0].AccountStatus == constant.Registered {
		return false, eerrs.ErrAccountAlreadyRegister.Wrap()
	}
	return true, nil
}

func (c *Caller) SendSimpleMsg(ctx context.Context, req *SendSingleMsgReq, key string) error {
	_, err := sendSimpleMsg.CallWithQuery(ctx, c.imApi, req, map[string]string{botstruct.Key: key})
	return err
}

func (c *Caller) AddNotificationAccount(ctx context.Context, req *user.AddNotificationAccountReq) error {
	_, err := addNotificationAccount.Call(ctx, c.imApi, req)
	return err
}

func (c *Caller) UpdateNotificationAccount(ctx context.Context, req *user.UpdateNotificationAccountInfoReq) error {
	_, err := updateNotificationAccount.Call(ctx, c.imApi, req)
	return err
}
