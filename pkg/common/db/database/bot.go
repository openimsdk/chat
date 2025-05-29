package database

import (
	"context"

	"github.com/openimsdk/chat/pkg/common/db/model/bot"
	tablebot "github.com/openimsdk/chat/pkg/common/db/table/bot"
	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/db/pagination"
	"github.com/openimsdk/tools/db/tx"
)

type BotDatabase interface {
	CreateAgent(ctx context.Context, jobs ...*tablebot.Agent) error
	TakeAgent(ctx context.Context, userID string) (*tablebot.Agent, error)
	FindAgents(ctx context.Context, userIDs []string) ([]*tablebot.Agent, error)
	UpdateAgent(ctx context.Context, userID string, data map[string]any) error
	DeleteAgents(ctx context.Context, userIDs []string) error
	PageAgents(ctx context.Context, userIDs []string, pagination pagination.Pagination) (int64, []*tablebot.Agent, error)

	TakeConversationRespID(ctx context.Context, convID, agentID string) (*tablebot.ConversationRespID, error)
	UpdateConversationRespID(ctx context.Context, convID, agentID string, data map[string]any) error
}

type botDatabase struct {
	tx         tx.Tx
	agent      tablebot.AgentInterface
	convRespID tablebot.ConversationRespIDInterface
}

func NewBotDatabase(cli *mongoutil.Client) (BotDatabase, error) {
	agent, err := bot.NewAgent(cli.GetDB())
	if err != nil {
		return nil, err
	}
	convRespID, err := bot.NewConversationRespID(cli.GetDB())
	if err != nil {
		return nil, err
	}
	return &botDatabase{
		tx:         cli.GetTx(),
		agent:      agent,
		convRespID: convRespID,
	}, nil
}

func (a *botDatabase) CreateAgent(ctx context.Context, agents ...*tablebot.Agent) error {
	return a.agent.Create(ctx, agents...)
}

func (a *botDatabase) TakeAgent(ctx context.Context, userID string) (*tablebot.Agent, error) {
	return a.agent.Take(ctx, userID)
}

func (a *botDatabase) FindAgents(ctx context.Context, userIDs []string) ([]*tablebot.Agent, error) {
	return a.agent.Find(ctx, userIDs)
}

func (a *botDatabase) UpdateAgent(ctx context.Context, userID string, data map[string]any) error {
	return a.agent.Update(ctx, userID, data)
}

func (a *botDatabase) DeleteAgents(ctx context.Context, userIDs []string) error {
	return a.agent.Delete(ctx, userIDs)
}

func (a *botDatabase) PageAgents(ctx context.Context, userIDs []string, pagination pagination.Pagination) (int64, []*tablebot.Agent, error) {
	return a.agent.Page(ctx, userIDs, pagination)
}

func (a *botDatabase) UpdateConversationRespID(ctx context.Context, convID, agentID string, data map[string]any) error {
	return a.convRespID.Update(ctx, convID, agentID, data)
}

func (a *botDatabase) TakeConversationRespID(ctx context.Context, convID, agentID string) (*tablebot.ConversationRespID, error) {
	return a.convRespID.Take(ctx, convID, agentID)
}
