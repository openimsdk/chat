package bot

import (
	"context"
	"crypto/rand"
	"time"

	"github.com/openimsdk/chat/pkg/common/constant"
	"github.com/openimsdk/chat/pkg/common/convert"
	"github.com/openimsdk/chat/pkg/common/mctx"
	"github.com/openimsdk/chat/pkg/protocol/bot"
	pbconstant "github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/protocol/user"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/utils/datautil"
)

func (b *botSvr) CreateAgent(ctx context.Context, req *bot.CreateAgentReq) (*bot.CreateAgentResp, error) {
	if req.Agent == nil {
		return nil, errs.ErrArgs.WrapMsg("req.Agent is nil")
	}

	now := time.Now()
	imToken, err := b.imCaller.ImAdminTokenWithDefaultAdmin(ctx)
	if err != nil {
		return nil, err
	}
	ctx = mctx.WithApiToken(ctx, imToken)
	if req.Agent.UserID != "" {
		req.Agent.UserID = constant.AgentUserIDPrefix + req.Agent.UserID
		users, err := b.imCaller.GetUsersInfo(ctx, []string{req.Agent.UserID})
		if err != nil {
			return nil, err
		}
		if len(users) > 0 {
			return nil, errs.ErrDuplicateKey.WrapMsg("agent userID already exists")
		}
	} else {
		randUserIDs := make([]string, 5)
		for i := range randUserIDs {
			randUserIDs[i] = constant.AgentUserIDPrefix + genID(10)
		}
		users, err := b.imCaller.GetUsersInfo(ctx, []string{req.Agent.UserID})
		if err != nil {
			return nil, err
		}
		if len(users) == len(randUserIDs) {
			return nil, errs.ErrDuplicateKey.WrapMsg("gen agent userID already exists, please try again")
		}
		for _, u := range users {
			if datautil.Contain(u.UserID, randUserIDs...) {
				continue
			}
			req.Agent.UserID = u.UserID
			break
		}
	}

	if err := b.imCaller.AddNotificationAccount(ctx, &user.AddNotificationAccountReq{
		UserID:         req.Agent.UserID,
		NickName:       req.Agent.NickName,
		FaceURL:        req.Agent.FaceURL,
		AppMangerLevel: pbconstant.AppRobotAdmin,
	}); err != nil {
		return nil, err
	}
	dbagent := convert.PB2DBAgent(req.Agent)
	dbagent.CreateTime = now
	err = b.database.CreateAgent(ctx, dbagent)
	if err != nil {
		return nil, err
	}
	return &bot.CreateAgentResp{}, nil
}

func (b *botSvr) UpdateAgent(ctx context.Context, req *bot.UpdateAgentReq) (*bot.UpdateAgentResp, error) {
	if _, err := b.database.TakeAgent(ctx, req.UserID); err != nil {
		return nil, errs.ErrArgs.Wrap()
	}

	if req.FaceURL != nil || req.NickName != nil {
		imReq := &user.UpdateNotificationAccountInfoReq{
			UserID: req.UserID,
		}
		if req.NickName != nil {
			imReq.NickName = *req.NickName
		}
		if req.FaceURL != nil {
			imReq.FaceURL = *req.FaceURL
		}
		err := b.imCaller.UpdateNotificationAccount(ctx, imReq)
		if err != nil {
			return nil, err
		}
	}

	update := ToDBAgentUpdate(req)
	err := b.database.UpdateAgent(ctx, req.UserID, update)
	if err != nil {
		return nil, err
	}

	return &bot.UpdateAgentResp{}, nil
}

func (b *botSvr) PageFindAgent(ctx context.Context, req *bot.PageFindAgentReq) (*bot.PageFindAgentResp, error) {
	total, agents, err := b.database.PageAgents(ctx, req.UserIDs, req.Pagination)
	if err != nil {
		return nil, err
	}
	_, userType, err := mctx.Check(ctx)
	if err != nil {
		return nil, err
	}
	if userType != constant.AdminUser {
		for i := range agents {
			agents[i].Key = ""
			agents[i].Url = ""
		}
	}
	return &bot.PageFindAgentResp{
		Total:  total,
		Agents: convert.BatchDB2PBAgent(agents),
	}, nil
}

func (b *botSvr) DeleteAgent(ctx context.Context, req *bot.DeleteAgentReq) (*bot.DeleteAgentResp, error) {
	err := b.database.DeleteAgents(ctx, req.UserIDs)
	if err != nil {
		return nil, err
	}
	return &bot.DeleteAgentResp{}, nil
}

func genID(l int) string {
	data := make([]byte, l)
	_, _ = rand.Read(data)
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
