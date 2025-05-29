package bot

import (
	"context"
)

type ConversationRespID struct {
	ConversationID     string `bson:"conversation_id"`
	AgentID            string `bson:"agent_id"`
	PreviousResponseID string `bson:"previous_response_id"`
}

func (ConversationRespID) TableName() string {
	return "conversation_resp_id"
}

type ConversationRespIDInterface interface {
	Create(ctx context.Context, elems ...*ConversationRespID) error
	Take(ctx context.Context, convID, agentID string) (*ConversationRespID, error)
	Update(ctx context.Context, convID, agentID string, data map[string]any) error
	Delete(ctx context.Context, convID, agentID string) error
}
