package bot

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/openimsdk/chat/pkg/botstruct"
	"github.com/openimsdk/chat/pkg/common/imapi"
	"github.com/openimsdk/chat/pkg/common/mctx"
	"github.com/openimsdk/chat/pkg/protocol/bot"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/tools/errs"
	"go.mongodb.org/mongo-driver/mongo"
)

func (b *botSvr) SendBotMessage(ctx context.Context, req *bot.SendBotMessageReq) (*bot.SendBotMessageResp, error) {
	agent, err := b.database.TakeAgent(ctx, req.AgentID)
	if err != nil {
		return nil, errs.ErrArgs.WrapMsg("agent not found")
	}
	convRespID, err := b.database.TakeConversationRespID(ctx, req.ConversationID, req.AgentID)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	}
	var respID string
	if convRespID != nil {
		respID = convRespID.PreviousResponseID
	}
	botreq := &botstruct.Request{
		Model: agent.Model,
		Input: []botstruct.InputItem{
			{
				Role:    botstruct.RoleDeveloper,
				Content: agent.Prompts,
			},
			{
				Role:    botstruct.RoleUser,
				Content: req.Content,
			},
		},
		PreviousResponseID: respID,
	}
	header := map[string]string{}
	header["Content-Type"] = "application/json"
	header["Authorization"] = "Bearer " + agent.Key
	postResp, err := b.httpClient.Post(ctx, agent.Url, header, botreq, b.timeout)
	if err != nil {
		return nil, errs.WrapMsg(err, "post bot failed")
	}

	var botResp botstruct.Response
	err = json.Unmarshal(postResp, &botResp)
	if err != nil {
		return nil, errs.WrapMsg(err, fmt.Sprintf("unmarshal post body failed, body:%s", string(postResp)))
	}

	newRespID, respContent, err := botResp.GetContentAndID()
	if err != nil {
		return nil, err
	}

	imToken, err := b.imCaller.ImAdminTokenWithDefaultAdmin(ctx)
	if err != nil {
		return nil, err
	}
	ctx = mctx.WithApiToken(ctx, imToken)
	err = b.imCaller.SendSimpleMsg(ctx, &imapi.SendSingleMsgReq{
		SendID:  agent.UserID,
		Content: respContent,
	}, req.Key)
	if err != nil {
		return nil, err
	}

	err = b.database.UpdateConversationRespID(ctx, req.ConversationID, agent.UserID, ToDBConversationRespIDUpdate(newRespID))
	if err != nil {
		return nil, err
	}
	return &bot.SendBotMessageResp{}, nil
}

func getContent(contentType int32, content string) (string, error) {
	switch contentType {
	case constant.Text:
		var elem botstruct.TextElem
		err := json.Unmarshal([]byte(content), &elem)
		if err != nil {
			return "", errs.ErrArgs.WrapMsg(err.Error())
		}
		return elem.Content, nil
	default:
		return "", errs.New("un support contentType").Wrap()
	}
}
