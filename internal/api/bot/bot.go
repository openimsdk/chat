package bot

import (
	"encoding/json"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/openimsdk/chat/internal/api/util"
	"github.com/openimsdk/chat/pkg/botstruct"
	"github.com/openimsdk/chat/pkg/common/imwebhook"
	"github.com/openimsdk/chat/pkg/protocol/bot"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/tools/a2r"
	"github.com/openimsdk/tools/apiresp"
	"github.com/openimsdk/tools/errs"
	"golang.org/x/sync/errgroup"
)

func New(botClient bot.BotClient, api *util.Api) *Api {
	return &Api{
		Api:       api,
		botClient: botClient,
	}
}

type Api struct {
	*util.Api
	botClient bot.BotClient
}

func (o *Api) CreateAgent(c *gin.Context) {
	a2r.Call(c, bot.BotClient.CreateAgent, o.botClient)
}

func (o *Api) DeleteAgent(c *gin.Context) {
	a2r.Call(c, bot.BotClient.DeleteAgent, o.botClient)
}

func (o *Api) UpdateAgent(c *gin.Context) {
	a2r.Call(c, bot.BotClient.UpdateAgent, o.botClient)
}

func (o *Api) PageFindAgent(c *gin.Context) {
	a2r.Call(c, bot.BotClient.PageFindAgent, o.botClient)
}

func (o *Api) AfterSendSingleMsg(c *gin.Context) {
	var (
		req = imwebhook.CallbackAfterSendSingleMsgReq{}
	)

	if err := c.BindJSON(&req); err != nil {
		apiresp.GinError(c, errs.ErrArgs.WithDetail(err.Error()).Wrap())
		return
	}
	if req.ContentType != constant.Text {
		apiresp.GinSuccess(c, nil)
		return
	}
	isAgent := botstruct.IsAgentUserID(req.RecvID)
	if !isAgent {
		apiresp.GinSuccess(c, nil)
		return
	}

	var elem botstruct.TextElem
	err := json.Unmarshal([]byte(req.Content), &elem)
	if err != nil {
		apiresp.GinError(c, errs.ErrArgs.WrapMsg("json unmarshal error: "+err.Error()))
		return
	}
	convID := getConversationIDByMsg(req.SessionType, req.SendID, req.RecvID, "")

	key, ok := c.GetQuery(botstruct.Key)
	if !ok {
		apiresp.GinError(c, errs.ErrArgs.WithDetail("missing key in query").Wrap())
		return
	}
	res, err := o.botClient.SendBotMessage(c, &bot.SendBotMessageReq{
		AgentID:        req.RecvID,
		ConversationID: convID,
		ContentType:    req.ContentType,
		Content:        elem.Content,
		Ex:             req.Ex,
		Key:            key,
	})
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiresp.GinSuccess(c, res)
}

func (o *Api) AfterSendGroupMsg(c *gin.Context) {
	var (
		req = imwebhook.CallbackAfterSendGroupMsgReq{}
	)
	if err := c.BindJSON(&req); err != nil {
		apiresp.GinError(c, errs.ErrArgs.WithDetail(err.Error()).Wrap())
		return
	}

	if req.ContentType != constant.AtText {
		apiresp.GinSuccess(c, nil)
	}
	key, ok := c.GetQuery(botstruct.Key)
	if !ok {
		apiresp.GinError(c, errs.ErrArgs.WithDetail("missing key in query").Wrap())
		return
	}

	var (
		elem botstruct.AtElem
		reqs []*bot.SendBotMessageReq
	)

	convID := getConversationIDByMsg(req.SessionType, req.SendID, "", req.GroupID)
	err := json.Unmarshal([]byte(req.Content), &elem)
	if err != nil {
		apiresp.GinError(c, errs.ErrArgs.WrapMsg("json unmarshal error: "+err.Error()))
	}
	for _, userID := range elem.AtUserList {
		if botstruct.IsAgentUserID(userID) {
			reqs = append(reqs, &bot.SendBotMessageReq{
				AgentID:        userID,
				ConversationID: convID,
				ContentType:    req.ContentType,
				Content:        elem.Text,
				Ex:             req.Ex,
				Key:            key,
			})
		}
	}
	if len(reqs) == 0 {
		apiresp.GinSuccess(c, nil)
	}

	g := errgroup.Group{}
	g.SetLimit(min(len(reqs), 5))
	for i := 0; i < len(reqs); i++ {
		i := i
		g.Go(func() error {
			_, err := o.botClient.SendBotMessage(c, reqs[i])
			if err != nil {
				return err
			}
			return nil
		})
	}

	err = g.Wait()
	if err != nil {
		apiresp.GinError(c, err)
		return
	}

	apiresp.GinSuccess(c, nil)
}

func getConversationIDByMsg(sessionType int32, sendID, recvID, groupID string) string {
	switch sessionType {
	case constant.SingleChatType:
		l := []string{sendID, recvID}
		sort.Strings(l)
		return "si_" + strings.Join(l, "_") // single chat
	case constant.WriteGroupChatType:
		return "g_" + groupID // group chat
	case constant.ReadGroupChatType:
		return "sg_" + groupID // super group chat
	case constant.NotificationChatType:
		l := []string{sendID, recvID}
		sort.Strings(l)
		return "sn_" + strings.Join(l, "_")
	}
	return ""
}
