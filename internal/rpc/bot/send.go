package bot

import (
	"context"
	"encoding/json"
	"time"

	"github.com/openimsdk/chat/pkg/botstruct"
	"github.com/openimsdk/chat/pkg/common/imapi"
	"github.com/openimsdk/chat/pkg/common/mctx"
	"github.com/openimsdk/chat/pkg/protocol/bot"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/tools/errs"
	"github.com/sashabaranov/go-openai"
)

func (b *botSvr) SendBotMessage(ctx context.Context, req *bot.SendBotMessageReq) (*bot.SendBotMessageResp, error) {
	agent, err := b.database.TakeAgent(ctx, req.AgentID)
	if err != nil {
		return nil, errs.ErrArgs.WrapMsg("agent not found")
	}
	//convRespID, err := b.database.TakeConversationRespID(ctx, req.ConversationID, req.AgentID)
	//if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
	//	return nil, err
	//}
	//var respID string
	//if convRespID != nil {
	//	respID = convRespID.PreviousResponseID
	//}

	aiCfg := openai.DefaultConfig(agent.Key)
	aiCfg.BaseURL = agent.Url
	aiCfg.HTTPClient = b.httpClient
	client := openai.NewClientWithConfig(aiCfg)
	aiReq := openai.ChatCompletionRequest{
		Model: agent.Model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: agent.Prompts,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: req.Content,
			},
		},
	}
	aiCtx, cancel := context.WithTimeout(ctx, time.Duration(b.timeout)*time.Second)
	defer cancel()
	completion, err := client.CreateChatCompletion(aiCtx, aiReq)
	if err != nil {
		return nil, errs.Wrap(err)
	}

	imToken, err := b.imCaller.ImAdminTokenWithDefaultAdmin(ctx)
	if err != nil {
		return nil, err
	}
	ctx = mctx.WithApiToken(ctx, imToken)

	content := "no response"
	if len(completion.Choices) > 0 {
		content = completion.Choices[0].Message.Content
	}
	err = b.imCaller.SendSimpleMsg(ctx, &imapi.SendSingleMsgReq{
		SendID:  agent.UserID,
		Content: content,
	}, req.Key)
	if err != nil {
		return nil, err
	}

	//err = b.database.UpdateConversationRespID(ctx, req.ConversationID, agent.UserID, ToDBConversationRespIDUpdate(completion.ID))
	//if err != nil {
	//	return nil, err
	//}
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
